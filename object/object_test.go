package object

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected Object
	}{
		{42, &Number{value: 42}},
		{int32(42), &Number{value: 42}},
		{int64(42), &Number{value: 42}},
		{1.5, &Number{value: 1.5}},
		{float32(1.5), &Number{value: 1.5}},
		{float64(1.5), &Number{value: 1.5}},
		{"String", &String{value: "String"}},
		{true, &Boolean{value: true}},
	}

	// To cover:
	// Unknown types should be converted to strings
	// Handle array types

	for _, test := range tests {
		results := New(test.input)

		if results.Inspect() != test.expected.Inspect() {
			t.Errorf("Did not get the right object back. Expected %#v got %#v", test.expected, results)
		}
	}
}
