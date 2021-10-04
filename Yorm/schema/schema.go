package schema

import (
	"Ghenyorm/Yorm/dialect"
	"go/ast"
	"reflect"
)

//a column in the database
type Field struct {
	Name string //name of column
	Type string //type of column data
	Tag  string //tag contrain (like`yorm:"some"``)
}

//a table in the database
type Schema struct {
	Model      interface{}       //model going to be reflect
	Name       string            //name of a table
	Fields     []*Field          //all columns
	FieldNames []string          //all column names
	fieldMap   map[string]*Field //relationship between Fields and FieldNames
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//the contrain is followed behind tag yorm
			if v, ok := p.Tag.Lookup("yorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
