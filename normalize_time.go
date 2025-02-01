package snippets

import (
	"reflect"
	"time"
)

//// NormalizeTime normalizes time.Time fields in a struct or a slice of structs to be second-only precision & in UTC.
//// It'll recurse into nested structs, arrays, and pointers.
//// To-add features:
//// - WithTimezone
//// - WithPrecision
//func NormalizeTime(obj any) {
//	if obj == nil {
//		return
//	}
//	normalizeTime(reflect.ValueOf(obj))
//}
//
//func normalizeTime(objVal reflect.Value) {
//	if objVal.Type().String() == "time.Time" {
//		utc := objVal.Interface().(time.Time).Truncate(time.Second).UTC()
//		objVal.Set(reflect.ValueOf(utc))
//	} else if objVal.Kind() == reflect.Pointer {
//		// If pointer, go inside.
//		if objVal.IsNil() {
//			return
//		}
//		normalizeTime(objVal.Elem())
//	} else if objVal.Kind() == reflect.Struct {
//		// If struct, iterate through the fields.
//		for i := 0; i < objVal.NumField(); i++ {
//			f := objVal.Field(i)
//			normalizeTime(f)
//		}
//	} else if objVal.Kind() == reflect.Array || objVal.Kind() == reflect.Slice {
//		// If slice, iterate through the elements.
//		for i := 0; i < objVal.Len(); i++ {
//			normalizeTime(objVal.Index(i))
//		}
//	} else if objVal.Kind() == reflect.Map {
//		// If map, iterate through the elements.
//		for _, key := range objVal.MapKeys() {
//			value := objVal.MapIndex(key)
//			// Normalize time.Time fields.
//			// Map element is not settable, so we set it in place.
//			if reflect.TypeOf(value.Interface()) == reflect.TypeOf(time.Time{}) {
//				utc := value.Interface().(time.Time).Truncate(time.Second).UTC()
//				objVal.SetMapIndex(key, reflect.ValueOf(utc))
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
//			newVal := reflect.New(value.Type()).Elem()
//			newVal.Set(value)
//			normalizeTime(newVal)
//			objVal.SetMapIndex(key, newVal)
//		}
//	} else if objVal.Kind() == reflect.Interface {
//		normalizeTime(objVal.Elem())
//	}
//}

// NormalizeTime normalizes time.Time fields in a struct or a slice of structs to be second-only precision & in UTC.
// It'll recurse into nested structs, arrays, and pointers.
// To-add features:
// - WithTimezone
// - WithPrecision
func NormalizeTime[T any](obj T) T {
	val := reflect.ValueOf(obj)
	retVal := reflect.New(val.Type()).Elem()
	retVal.Set(val)

	normalizeTime(retVal)

	return retVal.Interface().(T)
}

func normalizeTime(objVal reflect.Value) {
	if objVal.Type().String() == "time.Time" {
		utc := objVal.Interface().(time.Time).Truncate(time.Second).UTC()
		objVal.Set(reflect.ValueOf(utc))
		return
	}

	switch objVal.Kind() {
	case reflect.Pointer:
		if objVal.IsNil() {
			return
		}
		normalizeTime(objVal.Elem())
	case reflect.Struct:
		// If struct, iterate through the fields.
		for i := 0; i < objVal.NumField(); i++ {
			f := objVal.Field(i)
			normalizeTime(f)
		}
	case reflect.Array, reflect.Slice:
		// If slice, iterate through the elements.
		for i := 0; i < objVal.Len(); i++ {
			normalizeTime(objVal.Index(i))
		}
	case reflect.Map:
		// If map, iterate through the elements.
		for _, key := range objVal.MapKeys() {
			value := objVal.MapIndex(key)
			// Normalize time.Time fields.
			// Map element is not settable, so we set it in place.
			if reflect.TypeOf(value.Interface()) == reflect.TypeOf(time.Time{}) {
				utc := value.Interface().(time.Time).Truncate(time.Second).UTC()
				objVal.SetMapIndex(key, reflect.ValueOf(utc))
				continue
			}
			// Pointer is settable, so we recurse right away.
			if value.Kind() == reflect.Pointer {
				// If pointer, go inside.
				if value.IsNil() {
					continue
				}
				normalizeTime(value.Elem())
				continue
			}

			// Else, we need to create new value of the map element and recurse into it
			// since map element is not settable.
			newVal := reflect.New(reflect.TypeOf(value.Interface())).Elem()
			newVal.Set(reflect.ValueOf(value.Interface()))
			normalizeTime(newVal)
			objVal.SetMapIndex(key, newVal)
		}
	case reflect.Interface:
		if objVal.IsNil() {
			return
		}
		normalizeTime(objVal.Elem())
	default:
		// Do nothing.
	}
}
