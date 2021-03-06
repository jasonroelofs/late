package filter

import (
	"testing"

	"github.com/jasonroelofs/late/object"
)

func TestSize(t *testing.T) {
	tests := []struct {
		input    object.Object
		expected object.Object
	}{
		{object.New("A String"), object.New(8)},
	}

	for _, test := range tests {
		got := Size(test.input, make(Parameters))

		if got.Value() != test.expected.Value() {
			t.Errorf("Returned the wrong value. Expected %#v got %#v", test.expected, got)
		}
	}
}

func TestUpcase(t *testing.T) {
	tests := []struct {
		input    object.Object
		expected object.Object
	}{
		{object.New("A String"), object.New("A STRING")},
		{object.New("ALREADY BIG"), object.New("ALREADY BIG")},
	}

	for _, test := range tests {
		got := Upcase(test.input, make(Parameters))

		if got.Value() != test.expected.Value() {
			t.Errorf("Returned the wrong value. Expected %#v got %#v", test.expected, got)
		}
	}
}
