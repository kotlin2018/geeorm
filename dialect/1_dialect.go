// SQL 语句中的类型和 Go 语言中的类型是不同的， 例如Go 语言中的 int、int8、int16 等类型均对应 SQLite 中的 integer 类型。
//
// 因此实现 ORM 映射的第一步，需要思考如何将 Go 语言的类型映射为数据库中的类型。
//
// 同时，不同数据库支持的数据类型也是有差异的，即使功能相同，在 SQL 语句的表达上也可能有差异。

// ORM 框架往往需要兼容多种数据库，因此我们需要将差异的这一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。这部分代码称之为 dialect。
package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// Dialect 接口包含 2 个方法
type Dialect interface {
	DataTypeOf(goKind reflect.Value) string                    // 用于将 Go 语言的类型转换为数据库的数据类型。返回的是数据库的数据类型。
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表是否存在的 SQL 语句，参数是表名(table)。
}

// RegisterDialect 注册一个 dialect 到全局变量。如果新增加对某个数据库的支持，那么调用 RegisterDialect 即可注册到全局。
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 从全局变量获取一个 dialect
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}