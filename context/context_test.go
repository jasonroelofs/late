package context

import (
	"testing"

	"github.com/jasonroelofs/late/object"
)

func TestGlobalAssigns(t *testing.T) {
	c := New()

	c.Assign(Assigns{"key1": "value1", "key2": "value2"})

	if c.Get("key1").Value() != "value1" {
		t.Fatalf("Set the wrong value for key1, got %s", c.Get("key1"))
	}

	if c.Get("key2").Value() != "value2" {
		t.Fatalf("Set the wrong value for key2, got %s", c.Get("key2"))
	}

	if c.Get("key3") != object.NULL {
		t.Fatalf("key3 should not have been set, got %s", c.Get("key3"))
	}

	// Merge into the hash new values
	c.Assign(Assigns{"key1": "value3", "key3": "value4"})

	if c.Get("key1").Value() != "value3" {
		t.Fatalf("key1 didn't get updated, got %s", c.Get("key1"))
	}

	if c.Get("key2").Value() != "value2" {
		t.Fatalf("key2 should not have been updated, got %s", c.Get("key2"))
	}

	if c.Get("key3").Value() != "value4" {
		t.Fatalf("key3 did not get added, got %s", c.Get("key3"))
	}
}

func TestScoping(t *testing.T) {
	c := New()
	c.Assign(Assigns{"global_key": "global_value"})
	checkValueExists(t, c.Get("global_key"), "global_value")
	checkNoValue(t, c.Get("missing"))

	// 1
	c.PushScope()
	c.Set("scope_key", "scope_value")
	checkValueExists(t, c.Get("global_key"), "global_value")
	checkValueExists(t, c.Get("scope_key"), "scope_value")
	checkNoValue(t, c.Get("missing"))

	// 2
	c.PushScope()
	c.Set("deeper_key", "deeper_value")
	checkValueExists(t, c.Get("global_key"), "global_value")
	checkValueExists(t, c.Get("scope_key"), "scope_value")
	checkValueExists(t, c.Get("deeper_key"), "deeper_value")

	// 3
	c.PushScope()

	// Deeper scopes can override values of higher scopes
	c.Set("deeper_key", "even_deeper_value")
	checkValueExists(t, c.Get("global_key"), "global_value")
	checkValueExists(t, c.Get("scope_key"), "scope_value")
	checkValueExists(t, c.Get("deeper_key"), "even_deeper_value")

	// Now make sure that things reset!
	// 2
	c.PopScope()

	checkValueExists(t, c.Get("global_key"), "global_value")
	checkValueExists(t, c.Get("scope_key"), "scope_value")
	checkValueExists(t, c.Get("deeper_key"), "deeper_value")

	// 1
	c.PopScope()

	checkValueExists(t, c.Get("global_key"), "global_value")
	checkValueExists(t, c.Get("scope_key"), "scope_value")
	checkNoValue(t, c.Get("deeper_key"))

	// global
	c.PopScope()

	checkValueExists(t, c.Get("global_key"), "global_value")
	checkNoValue(t, c.Get("scope_key"))
	checkNoValue(t, c.Get("deeper_key"))

	// Doesn't crash if we pop from the top
	c.PopScope()

	checkValueExists(t, c.Get("global_key"), "global_value")
	checkNoValue(t, c.Get("scope_key"))
	checkNoValue(t, c.Get("deeper_key"))
}

func TestPromote(t *testing.T) {
	c := New()
	c.PushScope()
	c.Set("var", "value")

	c.Promote("var")

	c.PopScope()
	checkValueExists(t, c.Get("var"), "value")
}

func TestReadFile_NullReader(t *testing.T) {
	c := New()

	if c.ReadFile("file/path") != "ERROR: Reader not implemented. Cannot read content at file/path" {
		t.Fatalf("Did not set up the Null Reader properly")
	}
}

type TestReader struct{}

func (t *TestReader) Read(path string) string {
	return "I read from " + path
}

func TestReadFile_CustomReader(t *testing.T) {
	c := New(Reader(new(TestReader)))

	file := c.ReadFile("file/path")
	if file != "I read from file/path" {
		t.Fatalf("Did not set up the Reader properly. Got `%s`", file)
	}
}

func checkValueExists(t *testing.T, got object.Object, expected interface{}) {
	if got.Value() != expected {
		t.Errorf("Expected to find %#v but got %#v", expected, got)
	}
}

func checkNoValue(t *testing.T, value object.Object) {
	if value != object.NULL {
		t.Errorf("Expected to find nil but found %#v", value)
	}
}
