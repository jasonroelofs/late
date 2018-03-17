package filter

import (
	"github.com/jasonroelofs/late/object"
)

// Filters are functions that can act on any Late data type, manipulating and returning
// new values.
type FilterFunc func(object.Object) object.Object

type Filter struct {
	FilterFunc FilterFunc
}

type FilterContainer interface {
	AddFilter(string, FilterFunc)
}

func New(filterFunc FilterFunc) *Filter {
	return &Filter{
		FilterFunc: filterFunc,
	}
}

func (f *Filter) Call(input object.Object) object.Object {
	return f.FilterFunc(input)
}
