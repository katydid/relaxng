package relaxng

import (
	"reflect"
)

func removeTODOs(v reflect.Value) {
	num := v.NumField()
	for i := 0; i < num; i++ {
		f := v.Field(i)
		if f.Type().Kind() == reflect.Ptr {
			if !f.IsNil() {
				removeTODOs(f.Elem())
			}
		}
		if f.Type().Kind() == reflect.Slice {
			for i := 0; i < f.Len(); i++ {
				removeTODOs(f.Index(i))
			}
		}
		if f.Type().Kind() == reflect.Struct {
			removeTODOs(f)
		}
		name := v.Type().Field(i).Name
		if name == "Ns" {
			if f.String() == "TODO" {
				f.SetString("")
			}
		}
	}
}
