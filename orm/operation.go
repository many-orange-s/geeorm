package orm

import (
	"database/sql"
	"geeorm/clause"
	"geeorm/log"
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
