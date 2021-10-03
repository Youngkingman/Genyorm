package session

import (
	"Ghenyorm/Yorm/log"
	"Ghenyorm/Yorm/schema"
	"fmt"
	"reflect"
	"strings"
)

//parse the input struct into a table(schema), then cache it
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

//get current schema of the session
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model has not set yet")
	}
	return s.refTable
}

func (s *Session) CreateTable() (err error) {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err = s.Raw("CREATE TABLE %s (%s);", table.Name, desc).Exec()
	return
}

func (s *Session) DropTable() (err error) {
	_, err = s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
