package snippets_test

import (
	"testing"

	"github.com/slzhffktm/go-snippets"
)

func TestToPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		val := 42
		ptr := snippets.ToPtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("string", func(t *testing.T) {
		val := "hello"
		ptr := snippets.ToPtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type MyStruct struct {
			Field1 int
			Field2 string
		}
		val := MyStruct{Field1: 1, Field2: "test"}
		ptr := snippets.ToPtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("array", func(t *testing.T) {
		val := [3]int{1, 2, 3}
		ptr := snippets.ToPtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("map", func(t *testing.T) {
		val := map[string]int{"one": 1, "two": 2}
		ptr := snippets.ToPtr(val)
		if ptr == nil || (*ptr)["one"] != 1 || (*ptr)["two"] != 2 {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var val map[string]int
		ptr := snippets.ToPtr(val)
		if ptr == nil || *ptr != nil {
			t.Errorf("ToPtr() = %v, want %v", ptr, val)
		}
	})
}

func TestToVal(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		val := 42
		ptr := &val
		result := snippets.ToVal(ptr)
		if result != val {
			t.Errorf("ToVal() = %v, want %v", result, val)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		result := snippets.ToVal(ptr)
		if result != 0 {
			t.Errorf("ToVal() = %v, want %v", result, 0)
		}
	})

	t.Run("non-nil string pointer", func(t *testing.T) {
		val := "hello"
		ptr := &val
		result := snippets.ToVal(ptr)
		if result != val {
			t.Errorf("ToVal() = %v, want %v", result, val)
		}
	})

	t.Run("nil string pointer", func(t *testing.T) {
		var ptr *string
		result := snippets.ToVal(ptr)
		if result != "" {
			t.Errorf("ToVal() = %v, want %v", result, "")
		}
	})

	t.Run("non-nil array pointer", func(t *testing.T) {
		val := [3]int{1, 2, 3}
		ptr := &val
		result := snippets.ToVal(ptr)
		if result != val {
			t.Errorf("ToVal() = %v, want %v", result, val)
		}
	})

	t.Run("nil array pointer", func(t *testing.T) {
		var ptr *[3]int
		result := snippets.ToVal(ptr)
		if result != [3]int{} {
			t.Errorf("ToVal() = %v, want %v", result, [3]int{})
		}
	})

	t.Run("non-nil map pointer", func(t *testing.T) {
		val := map[string]int{"one": 1, "two": 2}
		ptr := &val
		result := snippets.ToVal(ptr)
		if result["one"] != 1 || result["two"] != 2 {
			t.Errorf("ToVal() = %v, want %v", result, val)
		}
	})

	t.Run("nil map pointer", func(t *testing.T) {
		var ptr *map[string]int
		result := snippets.ToVal(ptr)
		if result != nil {
			t.Errorf("ToVal() = %v, want %v", result, nil)
		}
	})
}
