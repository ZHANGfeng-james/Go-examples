package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

func init() {
	//FIXME 此处可能会导致空指针异常，var _ Dialect = (*sqlite3)(nil) 强制初始化 dialog.go
	RegisterDialect("sqlite3", &sqlite3{})
}

// DataTypeOf convert Go-type to RDMS-type
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool" // type of RDBMS
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice: // 使用实例？看看别人是怎么使用的
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (s *sqlite3) TableExistSQLStmt(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name=?;", args
}
