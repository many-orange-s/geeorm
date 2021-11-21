package orm

import (
	"database/sql"
	"geeorm/clause"
	"geeorm/log"
	"reflect"
)

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars); err != nil {
		log.Error(err)
	}
	return result, err
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Insert(values ...interface{}) (int64, error) {
	var recordvalues []interface{}
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordvalues = append(recordvalues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordvalues...)
	sqls, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sqls, vars).Exec()
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Find 我也没咋看懂
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
