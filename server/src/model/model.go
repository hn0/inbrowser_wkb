package model

import (
	"strings"
)

type Field struct {
	Name string
}

type Fields []Field

func CreateFields(names []string) *Fields {

	var f Fields
	f = make([]Field, len(names))

	for i, v := range names {
		f[i] = Field{
			Name: v,
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
