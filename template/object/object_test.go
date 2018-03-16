package object

import (
	"testing"
)

func TestFromNativeType(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected Object
	}{
		{42, &Number{Value: 42}},
		{int32(42), &Number{Value: 42}},
		{int64(42), &Number{Value: 42}},
		{1.5, &Number{Value: 1.5}},
		{float32(1.5), &Number{Value: 1.5}},
		{float64(1.5), &Number{Value: 1.5}},
		{"String", &String{Value: "String"}},
	}

	for _, test := range tests {
		results := FromNativeType(test.input)

		if results.Inspect() != test.expected.Inspect() {
			t.Errorf("Did not get the right object back. Expected %#v got %#v", test.expected, results)
		}
	}
}
