package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct {}

var _ Dialect = (*sqlite3)(nil) // *sqlite3断言nil，并赋值给Dialect接口，成功了；说明sqlite3实现了Dialect接口。

// 包在第一次加载时，会将 sqlite3 的 dialect 自动注册到全局。
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// 获取sqlite3 dialect的数据类型
func (s *sqlite3)DataTypeOf(goKind reflect.Value) string {
	switch goKind.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := goKind.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", goKind.Type().Name(), goKind.Kind()))
}

// TableExistSQL 返回判断表是否存在于数据库中的SQL
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}


