package geeorm

import (
	"github.com/kotlin2018/geeorm/dialect"
	"go/ast"
	"reflect"
)

type ITableName interface {
	TableName() string
}

// Field 数据库的字段
type Field struct {
	Name string //字段名
	Type string //字段类型
	Tag  string //约束条件，与结构体中的Tag对应。
}

// Schema 对象(object)和表(table)的转换。
type Schema struct {
	Model      interface{} // 被映射的对象 Model
	Name       string      // 表名 Name
	Fields     []*Field    // 所有字段 Fields
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField 获取一个字段
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// RecordValues 返回dest成员变量的值
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest)) //通过反射获取 dest的值
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// 将任意的对象解析为 Schema 实例
func Parse(model interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	var tableName string
	t, ok := model.(ITableName) // 断言dest实现了ITableName接口
	if !ok {
		tableName = modelType.Name() // 返回结构体的名称
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    model,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {  // 设置结构体Tag的key
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
