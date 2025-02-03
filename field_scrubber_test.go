package snippets_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/slzhffktm/go-snippets"
)

type ScrubFieldsSubStruct struct {
	FieldToScrub int
}

type ScrubFieldsStruct struct {
	FieldWithoutJson     string
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
		FieldWithoutJson:     "field_without_json",
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
	}

	res := snippets.ScrubFields(obj, map[string]struct{}{
		"listfieldtoscrub": {},
		"fieldtoscrub":     {},
		"toscrub":          {},
		"json_to_scrub":    {},
		"maptoscrub":       {},
		"nestedtoscrub":    {},
	})

	b, _ := json.Marshal(res)
	fmt.Println("SCRUBBED", string(b))

	b, _ = json.Marshal(obj)
	fmt.Println("ORI", string(b))
}
