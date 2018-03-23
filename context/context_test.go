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
