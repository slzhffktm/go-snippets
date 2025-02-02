package snippets

import (
	"reflect"
	"time"
)

// NormalizeTime normalizes time.Time fields in a struct or a slice of structs to be second-only precision & in UTC.
// It'll recurse into nested structs, arrays, and pointers.
// To-add features:
// - WithTimezone
// - WithPrecision
//func NormalizeTime(obj any) {
//	if obj == nil {
//		return
//	}
//	normalizeTime(reflect.ValueOf(obj))
//}
//
//func normalizeTime(value reflect.Value) {
//	if value.Type().String() == "time.Time" {
//		utc := value.Interface().(time.Time).Truncate(time.Second).UTC()
//		value.Set(reflect.ValueOf(utc))
//		return
//	}
//
//	switch value.Kind() {
//	case reflect.Pointer:
//		if value.IsNil() {
//			return
//		}
//		normalizeTime(value.Elem())
//	case reflect.Struct:
//		// If struct, iterate through the fields.
//		for i := 0; i < value.NumField(); i++ {
//			f := value.Field(i)
//			normalizeTime(f)
//		}
//	case reflect.Array, reflect.Slice:
//		// If slice, iterate through the elements.
//		for i := 0; i < value.Len(); i++ {
//			normalizeTime(value.Index(i))
//		}
//	case reflect.Map:
//		// If map, iterate through the elements.
//		for _, key := range value.MapKeys() {
//			value := value.MapIndex(key)
//			// Normalize time.Time fields.
//			// Map element is not settable, so we set it in place.
//			if reflect.TypeOf(value.Interface()) == reflect.TypeOf(time.Time{}) {
//				utc := value.Interface().(time.Time).Truncate(time.Second).UTC()
//				value.SetMapIndex(key, reflect.ValueOf(utc))
//				continue
//			}
//			// Pointer is settable, so we recurse right away.
//			if value.Kind() == reflect.Pointer {
//				// If pointer, go inside.
//				if value.IsNil() {
//					continue
//				}
//				normalizeTime(value.Elem())
//				continue
//			}
//
//			// Else, we need to create new value of the map element and recurse into it
//			// since map element is not settable.
//			newVal := reflect.New(reflect.TypeOf(value.Interface())).Elem()
//			newVal.Set(reflect.ValueOf(value.Interface()))
//			normalizeTime(newVal)
//			value.SetMapIndex(key, newVal)
//		}
//	case reflect.Interface:
//		if value.IsNil() {
//			return
//		}
//		normalizeTime(value.Elem())
//	default:
//		// Do nothing.
//	}
//}

func NormalizeTime[T any](obj any) T {
	res := normalizeTime(reflect.ValueOf(obj))
	return res.Interface().(T)
}

func normalizeTime(value reflect.Value) reflect.Value {
	if value.Type().String() == "time.Time" {
		newVal := reflect.New(value.Type()).Elem()

		utc := value.Interface().(time.Time).Truncate(time.Second).UTC()
		newVal.Set(reflect.ValueOf(utc))
		return newVal
	}

	switch value.Kind() {
	case reflect.Pointer:
		// If pointer, we recurse through if not nil.
		if value.IsNil() {
			return value
		}
		// Recurse through.
		res := normalizeTime(value.Elem())
		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(res.Addr())

		return newVal
	case reflect.Struct:
		// Copy the struct.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(value)

		// If struct, iterate through the fields.
		for i := 0; i < value.NumField(); i++ {
			// If anonymous, skip.
			if value.Type().Field(i).Anonymous {
				continue
			}

			f := value.Field(i)
			// Create copy.
			newField := reflect.New(f.Type()).Elem()
			// Recurse through.
			res := normalizeTime(f)
			newField.Set(res)

			newVal.Field(i).Set(newField)
		}

		return newVal
	case reflect.Array, reflect.Slice:
		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(value)

		// If array/slice, iterate through the elements.
		for i := 0; i < value.Len(); i++ {
			// Create copy.
			newField := reflect.New(value.Index(i).Type()).Elem()
			// Recurse through.
			res := normalizeTime(value.Index(i))
			newField.Set(res)

			newVal.Index(i).Set(newField)
		}

		return newVal
	case reflect.Map:
		// If map, iterate through the elements.

		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(reflect.MakeMap(value.Type()))

		for _, key := range value.MapKeys() {
			v := value.MapIndex(key)

			// Create copy.
			// Make v typed first, just in case it is interface.
			vWithType := reflect.ValueOf(v.Interface())
			newField := reflect.New(vWithType.Type()).Elem()
			// Recurse through.
			res := normalizeTime(vWithType)
			newField.Set(res)

			newVal.SetMapIndex(key, newField)
		}

		return newVal
	case reflect.Interface:
		// If interface, we recurse through if not nil.
		if value.IsNil() {
			return value
		}
		// Recurse through.
		res := normalizeTime(value.Elem())
		// Create copy.
		newVal := reflect.New(value.Type())
		newVal.Elem().Set(res)

		return newVal
	default:
		return value
	}
}
