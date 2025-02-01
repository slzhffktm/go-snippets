package snippets

import (
	"reflect"
	"time"
)

// NormalizeTime normalizes time.Time fields in a struct or a slice of structs to be second-only precision & in UTC.
// It'll recurse into nested structs, arrays, and pointers.
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

			// Normalize time.Time fields.
			if fieldVal.Type().String() == "time.Time" {
				utc := fieldVal.Interface().(time.Time).Truncate(time.Second).UTC()
				fieldVal.Set(reflect.ValueOf(utc))
				continue
			}

			// Recurse to field values.
			if (objVal.Kind() == reflect.Struct) ||
				(objVal.Kind() == reflect.Pointer && objVal.Type().Elem().Kind() == reflect.Struct) {
				normalizeTime(fieldVal)
			}
		}
	} else if objVal.Kind() == reflect.Slice {
		// If slice, iterate through the elements.
		for i := 0; i < objVal.Len(); i++ {
			normalizeTime(objVal.Index(i))
		}
	}
}
