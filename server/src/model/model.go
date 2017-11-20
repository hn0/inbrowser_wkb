package model

import (
	"reflect"
	"strings"
)

type Field struct {
	Name string
	Typ  reflect.Kind
	// values can be map of Typ (holding all column values!)
	Vals *Values
}

type Fields []Field
type Values []interface{}

func CreateFields(names []string) *Fields {

	var f Fields
	f = make([]Field, len(names))

	for i, v := range names {
		var vals Values
		vals = make([]interface{}, len(names))
		f[i] = Field{
			v,
			0,
			&vals,
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

// func (f *Fields) GetValues() []interface{} {
// acctually needed interface array for the row!?
// var vals []interface{}
// vals = make([]interface{}, len(*f))
// for i, _ := range *f {

//        todo create map for the intrface values!

// 	f[i].Values = vals[i]
// }
// return vals
// }
