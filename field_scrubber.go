package snippets

import "reflect"

const scrubbedField = "***scrubbed***"

// ScrubFields replaces the values of the fields in the given map with "***scrubbed***"
// if the field name is in the blacklist key.
// The blacklist is in map[string]struct{} to improve the lookup time.
// Make sure that the blacklist key is in lowercase.
// It'll recurse into nested maps & convert all struct into map to make all fields scrub-able.
func ScrubFields(fields map[string]any, blacklist map[string]struct{}) map[string]any {
	return scrubFields(reflect.ValueOf(fields), blacklist).Interface().(map[string]any)
}

func scrubFields(value reflect.Value, blacklist map[string]struct{}) reflect.Value {
	scrubbedFieldVal := reflect.ValueOf(scrubbedField)

	switch value.Kind() {
	case reflect.Pointer:
		// If pointer, we recurse through if not nil.
		if value.IsNil() {
			return value
		}
		// Recurse through.
		res := scrubFields(value.Elem(), blacklist)
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
			// If anonymous, skip.
			if value.Type().Field(i).Anonymous {
				continue
			}

			f := value.Field(i)

			// If the field name is in the blacklist, replace the value with scrubbedField.
			// Use JSON tag if available.
			fieldName := value.Type().Field(i).Name
			if jsonTag := value.Type().Field(i).Tag.Get("json"); jsonTag != "" {
				fieldName = jsonTag
			}
			if _, ok := blacklist[fieldName]; ok {
				newVal.SetMapIndex(reflect.ValueOf(fieldName), scrubbedFieldVal)
				continue
			}

			// Else, recurse through.

			// Create copy.
			newField := reflect.New(f.Type()).Elem()
			// Recurse through.
			res := scrubFields(f, blacklist)
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
			res := scrubFields(value.Index(i), blacklist)
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

			// If the field name is in the blacklist, replace the value with scrubbedField.
			fieldName := v.Type().Name()
			if _, ok := blacklist[fieldName]; ok {
				newVal.SetMapIndex(reflect.ValueOf(fieldName), scrubbedFieldVal)
				continue
			}

			// Create copy.
			// Make v typed first, just in case it is interface.
			vWithType := reflect.ValueOf(v.Interface())
			newField := reflect.New(vWithType.Type()).Elem()
			// Recurse through.
			res := scrubFields(vWithType, blacklist)
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
		res := scrubFields(value.Elem(), blacklist)
		// Create copy.
		newVal := reflect.New(value.Type())
		newVal.Elem().Set(res)

		return newVal
	default:
		return value
	}
}
