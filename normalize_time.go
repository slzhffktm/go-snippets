package snippets

import (
	"fmt"
	"reflect"
	"time"
)

// NormalizeTime normalizes time.Time fields in a struct or a slice of structs to be second-only precision & in UTC.
// It'll recurse into nested structs, arrays, and pointers.
// To-add features:
// - WithTimezone
// - WithPrecision
func NormalizeTime(obj any) {
	if obj == nil {
		return
	}
	normalizeTime(reflect.ValueOf(obj))
}

func normalizeTime(objVal reflect.Value) {
	if objVal.Type().String() == "time.Time" {
		utc := objVal.Interface().(time.Time).Truncate(time.Second).UTC()
		objVal.Set(reflect.ValueOf(utc))
	} else if objVal.Kind() == reflect.Pointer {
		// If pointer, go inside.
		if objVal.IsNil() {
			return
		}
		normalizeTime(objVal.Elem())
	} else if objVal.Kind() == reflect.Struct {
		// If struct, iterate through the fields & update time.Time.
		fields := reflect.VisibleFields(objVal.Type())

		// Iterate through the source fields
		for _, field := range fields {
			if field.Anonymous {
				continue
			}

			fieldVal := objVal.FieldByName(field.Name)
			normalizeTime(fieldVal)
		}
	} else if objVal.Kind() == reflect.Slice {
		// If slice, iterate through the elements.
		for i := 0; i < objVal.Len(); i++ {
			normalizeTime(objVal.Index(i))
		}
	} else if objVal.Kind() == reflect.Map {
		// If map, iterate through the elements.
		for _, key := range objVal.MapKeys() {
			fmt.Println(key)
			fmt.Println("ITSTIME", objVal.MapIndex(key))
			normalizeTime(objVal.MapIndex(key))
		}
	}
}
