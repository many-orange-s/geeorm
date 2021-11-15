package orm

import (
	"fmt"
	"geeorm/schema"
	"log"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Println("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var tableDeta []string
	for _, file := range table.Fields {
		tableDeta = append(tableDeta, fmt.Sprintf("%s %s %s", file.Name, file.Type, file.Tag))
	}

	decs := strings.Join(tableDeta, ",")

	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, decs)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.refTable.Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
