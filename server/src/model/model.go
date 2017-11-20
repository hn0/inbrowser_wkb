package model

import (
	"reflect"
	"strings"
)

type Field struct {
	Name string
	Typ  reflect.Kind
	// values can be map of Typ (holding all column values!)
	Vals []interface{}
}

type Fields []Field

func CreateFields(names []string) *Fields {

	var f Fields
	f = make([]Field, len(names))

	for i, v := range names {
		var vals []interface{}
		f[i] = Field{
			v,
			reflect.Invalid,
			vals,
		}
	}

	return &f
}

func (f *Fields) GetColumns(join string) string {

	var cols []string
	cols = make([]string, len(*f))
	for i, v := range *f {
		cols[i] = v.Name
	}

	return strings.Join(cols, join)
}

func (f *Fields) SetKind(index int, kind reflect.Kind) {
	(*f)[index].Typ = kind
}

func (f *Fields) AddRow(values []interface{}) {
	for i, _ := range values {
		v := (*f)[i].Vals
		(*f)[i].Vals = append(v, values[i])
	}
}
