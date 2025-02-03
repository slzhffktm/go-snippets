package snippets

import (
	"reflect"
	"strings"
)

const scrubbedField = "***scrubbed***"

// ScrubFields replaces the values of the fields in the given map with "***scrubbed***"
// if the field name is in the toScrub key.
// The toScrubs is in map[string]struct{} to improve the lookup time.
// Make sure that the toScrub keys are in lowercase.
// It'll recurse into nested maps & convert all struct into map to make all fields scrub-able.
func ScrubFields(fields map[string]any, toScrubs map[string]struct{}) (res map[string]any) {
	// Return value as-is if we panic.
	defer func() {
		if r := recover(); r != nil {
			res = fields
		}
	}()

	return scrubFields(reflect.ValueOf(fields), toScrubs).Interface().(map[string]any)
}

func scrubFields(value reflect.Value, toScrubs map[string]struct{}) reflect.Value {
	scrubbedFieldVal := reflect.ValueOf(scrubbedField)

	switch value.Kind() {
	case reflect.Pointer:
		// If pointer, we recurse through if not nil.
		if value.IsNil() {
			return value
		}
		// Recurse through.
		res := scrubFields(value.Elem(), toScrubs)
		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(res.Addr())

		return newVal
	case reflect.Struct:
		// If struct, create a map to make all fields scrub-able.

		// Create the map.
		newVal := reflect.MakeMap(reflect.TypeOf(map[string]any{}))

		// Iterate through the fields.
		for i := 0; i < value.NumField(); i++ {
			// If anonymous or unexported, skip.
			if value.Type().Field(i).Anonymous || !value.Type().Field(i).IsExported() {
				continue
			}

			f := value.Field(i)

			// If the field name is in the toScrubs, replace the value with scrubbedField.
			// Use JSON tag if available.
			fieldName := value.Type().Field(i).Name
			// If JSON tag is "-", skip.
			jsonTag := value.Type().Field(i).Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			if jsonTag != "" {
				fieldName = jsonTag
			}
			if _, ok := toScrubs[strings.ToLower(fieldName)]; ok {
				newVal.SetMapIndex(reflect.ValueOf(fieldName), scrubbedFieldVal)
				continue
			}

			// Else, recurse through.

			// If value type is interface/pointer/struct, no need to create copy.
			if f.Kind() == reflect.Interface || f.Kind() == reflect.Ptr || f.Kind() == reflect.Struct {
				res := scrubFields(f, toScrubs)
				newVal.SetMapIndex(reflect.ValueOf(fieldName), res)
				continue
			}

			// Create copy.
			newField := reflect.New(f.Type()).Elem()
			// Recurse through.
			res := scrubFields(f, toScrubs)
			newField.Set(res)

			newVal.SetMapIndex(reflect.ValueOf(fieldName), newField)
		}

		return newVal
	case reflect.Array, reflect.Slice:
		// If array/slice, iterate through the elements.

		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		// Set empty value with the same length.
		newVal.Set(reflect.MakeSlice(value.Type(), value.Len(), value.Len()))

		// If array/slice, iterate through the elements.
		for i := 0; i < value.Len(); i++ {
			// Create copy.
			newField := reflect.New(value.Index(i).Type()).Elem()
			// Recurse through.
			res := scrubFields(value.Index(i), toScrubs)
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

			// If the key is convertible to string, check if we need to scrub.
			if value.Type().Key().ConvertibleTo(reflect.TypeOf("")) {
				// If the field name is in the toScrubs, replace the value with scrubbedField.
				fieldName := key.String()
				if _, ok := toScrubs[strings.ToLower(fieldName)]; ok {
					newVal.SetMapIndex(reflect.ValueOf(fieldName), scrubbedFieldVal)
					continue
				}
			}

			// If not, recurse through.

			// If value type is interface/pointer/struct, no need to create copy.
			realKind := reflect.ValueOf(v.Interface()).Kind()
			if realKind == reflect.Interface || realKind == reflect.Ptr || realKind == reflect.Struct {
				res := scrubFields(v, toScrubs)
				newVal.SetMapIndex(key, res)
				continue
			}

			// Else, create copy.
			// Make v typed first, just in case it is interface.
			vWithType := reflect.ValueOf(v.Interface())
			newField := reflect.New(vWithType.Type()).Elem()
			// Recurse through.
			res := scrubFields(vWithType, toScrubs)
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
		res := scrubFields(value.Elem(), toScrubs)
		// Create copy.
		newVal := reflect.New(value.Type()).Elem()
		newVal.Set(res)

		return newVal
	default:
		return value
	}
}
