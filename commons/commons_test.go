package commons

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var empty = make([]string, 0)

type testCaseSetOperation struct {
	name    string
	A, B, R []string
}

func TestDifference(t *testing.T) {
	tt := []testCaseSetOperation{
		{
			name: "Empty - Empty = Empty",
			A:    empty,
			B:    empty,
			R:    empty,
		},
		{
			name: "A - Empty = A",
			A:    []string{"a", "b"},
			B:    empty,
			R:    []string{"a", "b"},
		},
		{
			name: "Empty - B = Empty",
			A:    empty,
			B:    []string{"a", "b"},
			R:    empty,
		},
		{
			name: "(a, b, c) - (a, d) = (b, c)",
			A:    []string{"a", "b", "c"},
			B:    []string{"a", "d"},
			R:    []string{"b", "c"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			difference := Difference(tc.A, tc.B)
			assert.ElementsMatch(t, tc.R, difference)
		})
	}

}

func TestUnion(t *testing.T) {
	tt := []testCaseSetOperation{
		{
			name: "Empty U Empty = Empty",
			A:    empty,
			B:    empty,
			R:    empty,
		},
		{
			name: "A U Empty = A",
			A:    []string{"b", "e", "c", "d"},
			B:    empty,
			R:    []string{"b", "e", "c", "d"},
		},
		{
			name: "Empty U B = B",
			A:    empty,
			B:    []string{"b", "e", "c", "d"},
			R:    []string{"b", "e", "c", "d"},
		},
		{
			name: "(a, b, c) U (b, e, c, d) = (a, b, c, d, e)",
			A:    []string{"a", "b", "c"},
			B:    []string{"b", "e", "c", "d"},
			R:    []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			union := Union(tc.A, tc.B)
			assert.ElementsMatch(t, tc.R, union)
		})
	}

}

func TestExtractToken(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "When header is correct", value: "Bearer hello", expected: "hello"},
		{name: "When header value is malformed", value: "aaa Bearer", expected: ""},
		{name: "When token is missing", value: "Bearer ", expected: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			header := http.Header{}
			header.Set("Authorization", tc.value)
			token := ExtractToken(header)

			assert.Equal(t, tc.expected, token)
		})
	}

}
