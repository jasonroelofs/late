package filter

import (
	"github.com/jasonroelofs/late/object"
)

// All secondary parameters are passed in as a hash map
// to the filter function that the function is expected to pull out
// by name.
type Parameters map[string]object.Object

// Filters are functions that can act on any Late data type, manipulating and returning
// new values.
type FilterFunc func(object.Object, Parameters) object.Object

type Filter struct {
	FilterFunc FilterFunc
}

func New(filterFunc FilterFunc) *Filter {
	return &Filter{
		FilterFunc: filterFunc,
	}
}

func (f *Filter) Call(input object.Object, params Parameters) object.Object {
	return f.FilterFunc(input, params)
}
