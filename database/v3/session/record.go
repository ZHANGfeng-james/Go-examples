package session

import (
	"reflect"

	"github.com/go-examples-with-tests/database/v3/clause"
	"github.com/go-examples-with-tests/database/v3/log"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	// INSERT INTO table_name(col1, col2, col3,...) VALUES (a1, a2, a3, ...), (b1, b2, b3, ...),...
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable() // 执行 Parse
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value)) // 解析出对象中各个字段的值
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values)) // reflect.Value
	log.Info(destSlice.Kind(), destSlice.Type().Name())

	destType := destSlice.Type().Elem() // Array, Chan, Map, Ptr, or Slice

}
