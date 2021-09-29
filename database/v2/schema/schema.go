package schema

import (
	"go/ast"
	"reflect"

	"github.com/go-examples-with-tests/database/v2/dialect"
)

// 一张 Table 中，Column 相关的信息
type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}       // 值，一般是指针类型的值
	Name       string            // 类型名，指针类型的值中解析出类型名，作为表名
	Fields     []*Field          // 表相关的所有列信息
	FieldNames []string          // 表相关的所有列名（字段名）
	fieldMap   map[string]*Field // 列名（字段名） - 列信息
}

type ITableName interface {
	TableName() string
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 依据具体的 dialect.Dialect 作类型转换
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	var tableName string
	t, ok := dest.(ITableName) // 是否实现ITableName接口
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i) // StructField 类型
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				// reflect.Indirect(reflect.New(p.Type)) --> 创建指针类型实例，并访问
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest)) // reflect.Value
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		// reflect.Value struct --> value
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
