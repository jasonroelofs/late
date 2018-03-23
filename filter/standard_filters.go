package filter

import (
	"strings"

	"github.com/jasonroelofs/late/object"
)

func Size(input object.Object, _ Parameters) object.Object {
	switch input.Type() {
	case object.OBJ_STRING:
		return object.New(len(input.Value().(string)))
	default:
		return input
	}
}

func Upcase(input object.Object, _ Parameters) object.Object {
	switch inputT := input.Value().(type) {
	case string:
		return object.New(strings.ToUpper(inputT))
	default:
		return input
	}
}

func Replace(input object.Object, params Parameters) object.Object {
	// TODO Type checking and verification that the parameters are in-fact
	// Strings and will work here.
	in := input.Value().(string)
	replace := params["replace"].Value().(string)
	with := params["with"].Value().(string)

	out := strings.Replace(in, replace, with, -1)

	return object.New(out)
}
