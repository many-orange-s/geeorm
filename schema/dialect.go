package schema

import "reflect"

// 接口的话都要加大括号
var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectMap[name]
	return dialect, ok
}
