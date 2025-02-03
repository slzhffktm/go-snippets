package snippets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slzhffktm/go-snippets"
)

const scrubbedField = "***scrubbed***"

type ScrubFieldsSubStruct struct {
	FieldToScrub int
}

type ScrubFieldsStruct struct {
	privateField         string
	FieldWithJson        string `json:"field_with_json"`
	JsonToScrub          string `json:"json_to_scrub"`
	FieldIgnoreJson      string `json:"-"`
	ListField            []string
	ListFieldToScrub     []string
	FieldPointerToString *string
	SubStruct            ScrubFieldsSubStruct
}

func TestScrubFields(t *testing.T) {
	s := ScrubFieldsStruct{
		privateField:         "private_field",
		FieldWithJson:        "field_with_json",
		ListField:            []string{"asdf"},
		ListFieldToScrub:     []string{"asdf"},
		FieldPointerToString: snippets.ToPtr("asdfasdf"),
		FieldIgnoreJson:      "field_ignore_json",
		SubStruct:            ScrubFieldsSubStruct{FieldToScrub: 123},
	}

	obj := map[string]any{
		"struct":  s,
		"string":  "string",
		"toscrub": true,
		"maptoscrub": map[int]int{
			1: 1,
		},
		"nested_map": map[int]any{
			1: map[string]any{
				"nestedtoscrub": "nested",
			},
		},
		"list_of_map": []map[string]any{
			{
				"nestedtoscrub": "nestedtoscrub",
				"dontscrub":     "dontscrub",
				"nestedagain": []any{
					map[string]any{
						"nestedtoscrub": "nestedtoscrub",
					},
				},
			},
		},
	}

	res := snippets.ScrubFields(obj, map[string]struct{}{
		"listfieldtoscrub": {},
		"fieldtoscrub":     {},
		"toscrub":          {},
		"json_to_scrub":    {},
		"maptoscrub":       {},
		"nestedtoscrub":    {},
	})

	// Assert that the scrubbed fields are replaced with scrubbedField.
	assert.Equal(t, map[string]any{
		"struct": map[string]any{
			"field_with_json":      "field_with_json",
			"json_to_scrub":        scrubbedField,
			"ListField":            []string{"asdf"},
			"ListFieldToScrub":     scrubbedField,
			"FieldPointerToString": snippets.ToPtr("asdfasdf"),
			"SubStruct": map[string]any{
				"FieldToScrub": scrubbedField,
			},
		},
		"string":     "string",
		"toscrub":    scrubbedField,
		"maptoscrub": scrubbedField,
		"nested_map": map[int]any{
			1: map[string]any{
				"nestedtoscrub": scrubbedField,
			},
		},
		"list_of_map": []map[string]any{
			{
				"nestedtoscrub": scrubbedField,
				"dontscrub":     "dontscrub",
				"nestedagain": []any{
					map[string]any{
						"nestedtoscrub": scrubbedField,
					},
				},
			},
		},
	}, res)
	// Assert that the original object is not modified.
	assert.Equal(t, map[string]any{
		"struct":  s,
		"string":  "string",
		"toscrub": true,
		"maptoscrub": map[int]int{
			1: 1,
		},
		"nested_map": map[int]any{
			1: map[string]any{
				"nestedtoscrub": "nested",
			},
		},
		"list_of_map": []map[string]any{
			{
				"nestedtoscrub": "nestedtoscrub",
				"dontscrub":     "dontscrub",
				"nestedagain": []any{
					map[string]any{
						"nestedtoscrub": "nestedtoscrub",
					},
				},
			},
		},
	}, obj)
}
