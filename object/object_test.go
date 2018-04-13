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

func TestNew_DeepObjects(t *testing.T) {
	obj := New(map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{"nested1": 1, "nested2": 2},
	})

	if obj.Type() != TYPE_HASH {
		t.Fatalf("Did not turn the map into a Hash, got %s", obj.Type())
	}

	hash := obj.(*Hash)
	key1 := hash.Get(New("key1"))
	if key1.Value() != "value1" {
		t.Fatalf("Did not find the right value for `key1`, got %#v", key1)
	}

	key2 := hash.Get(New("key2"))
	if key2.Type() != TYPE_HASH {
		t.Fatalf("Did not turn the key2 map into a Hash, got %#v", key2)
	}

	hash = key2.(*Hash)
	nested1 := hash.Get(New("nested1"))
	if nested1.Value() != float64(1) {
		t.Fatalf("Did not find the right value for `key2.nested1`, got %#v", nested1)
	}

	nested2 := hash.Get(New("nested2"))
	if nested2.Value() != float64(2) {
		t.Fatalf("Did not find the right value for `key2.nested1`, got %#v", nested2)
	}
}

func TestNewReturnsObjectsRaw(t *testing.T) {
	str := New("A test string")
	copy := New(str)

	if str != copy {
		t.Errorf("Did not properly copy the object. Got %#v", copy)
	}
}

func TestTruthy(t *testing.T) {
	tests := []struct {
		input    Object
		expected bool
	}{
		{TRUE, true},
		{New("String"), true},
		{New(12345), true},
		{&Array{}, true},
		{&Hash{}, true},
		{FALSE, false},
		{NULL, false},
	}

	for _, test := range tests {
		if Truthy(test.input) != test.expected {
			t.Errorf("Truthy test failed for %v, expected %t", test.input, test.expected)
		}
	}
}
