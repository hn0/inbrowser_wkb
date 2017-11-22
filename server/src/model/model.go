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
	// sizes in bytes!
	sizes []int
}

type Fields []Field

func CreateFields(names []string) *Fields {

	var f Fields
	f = make([]Field, len(names))

	for i, v := range names {
		var vals []interface{}
		var szs []int
		szs = make([]int, len(names))
		f[i] = Field{
			v,
			reflect.Invalid,
			vals,
			szs,
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
		s := (*f)[i].sizes

		// ok, cannot access bytes so easily
		//  we need to detect blob type!
		var sz int
		v2 := values[i].(*interface{})
		if (*f)[i].Typ == reflect.Array {
			sz = len((*v2).([]byte))
		} else {
			sz = int(reflect.TypeOf(v2).Size())
		}

		// fmt.Println(sz)

		(*f)[i].Vals = append(v, values[i])
		(*f)[i].sizes = append(s, sz)
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
		} else if fl.Typ == reflect.Array {
			v := fl.Vals[i].(*interface{})
			ret[fl.Name] = (*v).([]byte)
		} else {
			ret[fl.Name] = nil
		}
	}
	return ret
}
