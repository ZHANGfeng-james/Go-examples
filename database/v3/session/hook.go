package session

import (
	"reflect"

	"github.com/go-examples-with-tests/database/v3/log"
)

const (
	BeforeQuery = "BeforeQuery"
	AfterQuery  = "AfterQuery"

	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"

	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"

	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallHoookMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		// 表示 s.RefTable().Model 对应结构体的某个指定变量
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 每个钩子的入参都是 *Session
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
