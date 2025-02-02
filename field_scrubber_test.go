package snippets_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	loggerservice "github.com/xendit/xsh-go-logger/v2"

	"github.com/slzhffktm/go-snippets"
)

type StructTest struct {
	FieldWithoutJson     string
	privateField         string
	FieldWithJson        string `json:"field_with_json"`
	FieldIgnoreJson      string `json:"-"`
	ListField            []string
	FieldPointerToString *string

	Ring string
}

func TestLogger(t *testing.T) {
	logger, err := loggerservice.InitLogger(
		"test", "test",
		loggerservice.WithAdditionalBlacklist([]string{"ring", "ListField"}))
	require.NoError(t, err)

	s := StructTest{
		FieldWithoutJson:     "field_without_json",
		privateField:         "private_field",
		FieldWithJson:        "field_with_json",
		ListField:            []string{"asdf"},
		FieldPointerToString: snippets.ToPtr("asdfasdf"),
		FieldIgnoreJson:      "field_ignore_json",
		Ring:                 "asdf",
	}

	logger.Info(context.Background(), loggerservice.Fields{
		"struct": &s,
	}, "Woi")

	fmt.Println(s)
}
