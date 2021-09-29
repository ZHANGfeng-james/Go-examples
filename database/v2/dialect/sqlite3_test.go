package dialect

import (
	"reflect"
	"testing"
)

func TestDataTypeOf(t *testing.T) {
	p := []struct {
		Values interface{}
		Type   string
	}{
		{"Tom", "text"},
		{123, "integer"},
		{1.23, "real"},
		{[]int{1, 2, 3, 4}, "blob"},
	}
	// 在执行 TestDataTypeOf 时，已经调用了sqlite3.go中的init()
	sqlDB, ok := GetDialect("sqlite3")
	if ok {
		for _, parameter := range p {
			if typ := sqlDB.DataTypeOf(reflect.ValueOf(parameter.Values)); typ != parameter.Type {
				t.Fatalf("Type of %v is %s, got:%s", parameter.Values, parameter.Type, typ)
			}
		}
	}
}
