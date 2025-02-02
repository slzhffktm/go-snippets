package snippets_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/slzhffktm/go-snippets"
)

type subStruct struct {
	TimeField             time.Time
	TimeFieldPointer      *time.Time
	TimeFieldPointerEmpty *time.Time
}

type someStruct struct {
	TimeField              time.Time
	TimeFieldPointer       *time.Time
	TimeFieldPointerEmpty  *time.Time
	NonTimeField           string
	SubStruct              subStruct
	SubStructPointer       *subStruct
	SubStructList          []subStruct
	SubStructPointerToList *[]subStruct
	TimeList               []time.Time
	MapField               map[string]time.Time
	MapToPointerField      map[string]*time.Time
	MapToStructField       map[string]subStruct
	// MapStringAny should be in this format:
	// MapStringAny: map[string]any{
	//			"key": timeNow,
	//			"mapInMap": map[string]time.Time{
	//				"nestedKey": timeNow,
	//			},
	//			"substruct": subStruct{
	//				TimeField:        timeNow,
	//				TimeFieldPointer: copyTimeToPtr(timeNow),
	//			},
	//		},
	MapStringAny map[string]any
}

func (s someStruct) assertTimeFields(t *testing.T, timeTest time.Time) {
	require.Equal(t, timeTest, s.TimeField)
	require.Equal(t, timeTest, *s.TimeFieldPointer)
	require.Nil(t, s.TimeFieldPointerEmpty)
	require.Equal(t, timeTest, s.SubStruct.TimeField)
	require.Equal(t, timeTest, *s.SubStruct.TimeFieldPointer)
	require.Equal(t, timeTest, s.SubStructPointer.TimeField)
	require.Equal(t, timeTest, *s.SubStructPointer.TimeFieldPointer)
	for _, sub := range s.SubStructList {
		require.Equal(t, timeTest, sub.TimeField)
		require.Equal(t, timeTest, *sub.TimeFieldPointer)
	}
	for _, sub := range *s.SubStructPointerToList {
		require.Equal(t, timeTest, sub.TimeField)
		require.Equal(t, timeTest, *sub.TimeFieldPointer)
	}
	for _, timeV := range s.TimeList {
		require.Equal(t, timeTest, timeV)
	}
	for _, timeV := range s.MapField {
		require.Equal(t, timeTest, timeV)
	}
	for _, timeV := range s.MapToPointerField {
		require.Equal(t, timeTest, *timeV)
	}
	for _, sub := range s.MapToStructField {
		require.Equal(t, timeTest, sub.TimeField)
		require.Equal(t, timeTest, *sub.TimeFieldPointer)
	}
	require.Equal(t, timeTest, s.MapStringAny["key"])
	require.Equal(t, timeTest, s.MapStringAny["mapInMap"].(map[string]time.Time)["nestedKey"])
	require.Equal(t, timeTest, s.MapStringAny["substruct"].(subStruct).TimeField)
	require.Equal(t, timeTest, *s.MapStringAny["substruct"].(subStruct).TimeFieldPointer)
}

func copyTimeToPtr(t time.Time) *time.Time {
	newT := t
	return snippets.ToPtr(newT)
}

func TestNormalizeTime(t *testing.T) {
	timeNow := time.Now().Local()
	obj := someStruct{
		TimeField:             timeNow,
		TimeFieldPointer:      copyTimeToPtr(timeNow),
		TimeFieldPointerEmpty: nil,
		NonTimeField:          "non-time",
		SubStruct: subStruct{
			TimeField:        timeNow,
			TimeFieldPointer: copyTimeToPtr(timeNow),
		},
		SubStructPointer: &subStruct{
			TimeField:        timeNow,
			TimeFieldPointer: copyTimeToPtr(timeNow),
		},
		SubStructList: []subStruct{
			{
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
			{
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
		},
		SubStructPointerToList: &[]subStruct{
			{
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
			{
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
		},
		TimeList: []time.Time{timeNow, timeNow},
		MapField: map[string]time.Time{
			"key1": timeNow,
		},
		MapToPointerField: map[string]*time.Time{
			"key1": copyTimeToPtr(timeNow),
		},
		MapToStructField: map[string]subStruct{
			"key1": {
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
		},
		MapStringAny: map[string]any{
			"key": timeNow,
			"mapInMap": map[string]time.Time{
				"nestedKey": timeNow,
			},
			"substruct": subStruct{
				TimeField:        timeNow,
				TimeFieldPointer: copyTimeToPtr(timeNow),
			},
		},
	}

	res := snippets.NormalizeTime[someStruct](obj)

	// Assert res time fields are normalized.
	res.assertTimeFields(t, timeNow.Truncate(time.Second).UTC())
	// Assert original obj time fields are not modified.
	obj.assertTimeFields(t, timeNow)
}

func TestNormalizeTime_TimeParam(t *testing.T) {
	timeNow := time.Now().Local()
	obj := timeNow

	snippets.NormalizeTime[someStruct](obj)

	require.Equal(t, timeNow.Truncate(time.Second).UTC(), obj)
}
