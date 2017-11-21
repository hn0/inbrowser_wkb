package model

import (
	"fmt"
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

func (f *Fields) GetRecord(i int) map[string]interface{} {
	var ret map[string]interface{}
	ret = make(map[string]interface{}, len(*f))
	for _, fl := range *f {
		if fl.Typ == reflect.Int64 {
			ret[fl.Name] = fl.Vals[i]
		} else if fl.Typ == reflect.String && fl.Vals[i] != nil {
			v := fl.Vals[i].(*interface{})
			ret[fl.Name] = fmt.Sprintf("%s", *v)
		} else {
			ret[fl.Name] = nil
		}
	}
	return ret
}
