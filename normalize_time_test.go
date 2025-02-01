package snippets_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/slzhffktm/go-snippets"
)

type subStruct struct {
	TimeField        time.Time
	TimeFieldPointer *time.Time
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
	}

	snippets.NormalizeTime(&obj)

	obj.assertTimeFields(t, timeNow.Truncate(time.Second).UTC())
}

func TestNormalizeTime_TimeParam(t *testing.T) {
	timeNow := time.Now().Local()
	obj := timeNow

	snippets.NormalizeTime(&obj)

	require.Equal(t, timeNow.Truncate(time.Second).UTC(), obj)
}
