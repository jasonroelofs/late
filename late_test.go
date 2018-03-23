package late

import (
	"testing"

	"github.com/jasonroelofs/late/filter"
	"github.com/jasonroelofs/late/object"
)

func TestAddAndFindingFilters(t *testing.T) {

	AddFilter("size", func(input object.Object, _ filter.Parameters) object.Object {
		switch inputT := input.Value().(type) {
		case string:
			return object.New(len(inputT))
		default:
			return input
		}
	})

	method := FindFilter("size")

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

func TestFindTag(t *testing.T) {
	a1 := FindTag("assign")
	a2 := FindTag("assign")

	if &a1 == &a2 {
		t.Fatalf("Should have created new versions of the assign tag each call")
	}
}
