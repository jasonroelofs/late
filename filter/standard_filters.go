package filter

import (
	"strings"

	"github.com/jasonroelofs/late/object"
)

func ApplyStandardFilters(container FilterContainer) {
	container.AddFilter("size", size)
	container.AddFilter("upcase", upcase)
}

func size(input object.Object) object.Object {
	switch input.Type() {
	case object.OBJ_STRING:
		return object.New(len(input.Value().(string)))
	default:
		return input
	}
}

func upcase(input object.Object) object.Object {
	switch inputT := input.Value().(type) {
	case string:
		return object.New(strings.ToUpper(inputT))
	default:
		return input
	}
}
