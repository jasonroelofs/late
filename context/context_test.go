package context

import (
	"testing"

	"github.com/jasonroelofs/late/object"
)

func TestAddAndFindingFilters(t *testing.T) {
	c := New()

	c.AddFilter("size", func(input object.Object) object.Object {
		switch inputT := input.Value().(type) {
		case string:
			return object.New(len(inputT))
		default:
			return input
		}
	})

	filter := c.FindFilter("size")

	if filter == nil {
		t.Fatalf("Unable to find filter with name size")
	}

	resultObj := filter.Call(object.New("String"))

	if resultObj.Type() != object.OBJ_NUMBER {
		t.Fatalf("The resulting object is not a number, got %T", resultObj)
	}

	result := resultObj.Value().(float64)
	if result != 6 {
		t.Fatalf("Calling the filter did not return the right size, got %f", result)
	}
}
