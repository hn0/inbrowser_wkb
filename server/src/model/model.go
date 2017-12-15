package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Filter struct {
	Field    string
	Relation string
	Value    interface{}
}

type Field struct {
	Name string
	Typ  reflect.Kind
	// values can be map of Typ (holding all column values!)
	Vals []interface{}
	// sizes in bytes! aggregate will do for now
	size        int
	constraints []Filter
}

type Fields []Field

func CreateFields(names []string) *Fields {

	var f Fields
	f = make([]Field, len(names))

	for i, v := range names {
		var vals []interface{}
		var con []Filter
		f[i] = Field{
			v,
			reflect.Invalid,
			vals,
			0,
			con,
		}
	}

	return &f
}

func (f *Field) AddConstraint(column string, relation string, value interface{}) {
	c := Filter{column, relation, value}
	f.constraints = append(f.constraints, c)
}

func (f *Fields) GetColumns(join string) string {

	var cols []string
	cols = make([]string, len(*f))
	for i, v := range *f {
		cols[i] = v.Name
	}

	return strings.Join(cols, join)
}

func (f *Fields) GetConstraints() string {
	var cond string

	for _, f := range *f {
		for _, c := range f.constraints {
			// current support is for AND only, after all this is just a proof of concept
			if len(cond) > 0 {
				cond = " AND " + cond
			}
			// TODO: support for other types of values is missing (current support is only for numeric types)
			cond = fmt.Sprintf("%s %s %s %d",
				cond,
				c.Field,
				c.Relation,
				c.Value)
		}
	}

	if len(cond) > 0 {
		return " WHERE " + cond
	}
	return ";"
}

func (f *Fields) SetKind(index int, kind reflect.Kind) {
	(*f)[index].Typ = kind
}

func (f *Fields) AddRow(values []interface{}) {
	for i, _ := range values {
		v := (*f)[i].Vals

		// ok, cannot access bytes so easily
		//  we need to detect blob type!
		var sz int
		v2 := values[i].(*interface{})
		if (*f)[i].Typ == reflect.Array {
			sz = len((*v2).([]byte))
		} else if (*f)[i].Typ == reflect.Int64 {
			// for demo, id will be truncated
			sz = 8
		} else {
			sz = int(reflect.TypeOf(v2).Size())
		}

		// fmt.Println(sz)

		(*f)[i].Vals = append(v, values[i])
		(*f)[i].size += sz
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
			ret[fl.Name] = fl.Vals[i].(*interface{})
		}
	}
	return ret
}

func (f *Fields) SizeOf() int {
	var sz int
	sz = 0
	for _, v := range *f {
		sz += v.size
	}
	return sz
}
