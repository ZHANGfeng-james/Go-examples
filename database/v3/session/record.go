package session

import (
	"errors"
	"reflect"

	"github.com/go-examples-with-tests/database/v3/clause"
	"github.com/go-examples-with-tests/database/v3/log"
)

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	// 用于链式调用
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) Insert(values ...interface{}) (int64, error) {
	// INSERT INTO table_name(col1, col2, col3,...) VALUES (a1, a2, a3, ...), (b1, b2, b3, ...),...
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallHoookMethod(BeforeInsert, value)
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
	s.CallHoookMethod(AfterInsert, nil)
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	s.CallHoookMethod(BeforeQuery, nil)
	// var users []User --> Find(&users)
	destSlice := reflect.Indirect(reflect.ValueOf(values)) // reflect.Value --> []User
	destType := destSlice.Type().Elem()                    // Array, Chan, Map, Ptr, or Slice reflect.Type --> User

	// reflect.New(destType) --> reflect.Value
	log.Info(reflect.New(destType).Kind())
	table := s.Model(reflect.New(destType).Interface()).RefTable()

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
			// 向 values 中添加 dest 按序铺平的各个字段指针
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 依据数据库查询值，为 values 赋值
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallHoookMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallHoookMethod(BeforeQuery, nil)
	// support map[string]interface{}
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		// also support: "Name", "Tom", "Age", 18
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallHoookMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallHoookMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallHoookMethod(AfterDelete, nil)
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) First(value interface{}) error {
	// var user &User --> session.First(user)
	dest := reflect.Indirect(reflect.ValueOf(value))

	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem() // 创建 []User 值
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
