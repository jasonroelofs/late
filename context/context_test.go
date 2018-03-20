package context

import (
	"testing"

	"github.com/jasonroelofs/late/filter"
	"github.com/jasonroelofs/late/object"
)

func TestAddAndFindingFilters(t *testing.T) {
	c := New()

	c.AddFilter("size", func(input object.Object, _ filter.Parameters) object.Object {
		switch inputT := input.Value().(type) {
		case string:
			return object.New(len(inputT))
		default:
			return input
		}
	})

	method := c.FindFilter("size")

	if method == nil {
		t.Fatalf("Unable to find filter with name size")
	}

	resultObj := method.Call(object.New("String"), make(filter.Parameters))

	if resultObj.Type() != object.OBJ_NUMBER {
		t.Fatalf("The resulting object is not a number, got %T", resultObj)
	}

	result := resultObj.Value().(float64)
	if result != 6 {
		t.Fatalf("Calling the filter did not return the right size, got %f", result)
	}
}

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
