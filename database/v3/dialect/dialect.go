package dialect

import (
	"fmt"
	"reflect"
)

var dialectsMap = map[string]Dialect{} // 进程全局保存注册的 name - Dialect

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                        // Go-type convert to RDMS-type
	TableExistSQLStmt(tableName string) (string, []interface{}) // 指定tablename是否存在的SQL语句
}

func RegisterDialect(name string, dialect Dialect) {
	_, ok := GetDialect(name)
	if ok {
		panic(fmt.Sprintf("dialect for %s just registe once", name))
	}
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
