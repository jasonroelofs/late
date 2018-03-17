package context

import (
	"github.com/jasonroelofs/late/filter"
)

type Context struct {
	filters map[string]*filter.Filter
}

func New() *Context {
	ctx := &Context{
		filters: make(map[string]*filter.Filter),
	}

	filter.ApplyStandardFilters(ctx)

	return ctx
}

func (c *Context) AddFilter(name string, filterFunc filter.FilterFunc) {
	c.filters[name] = filter.New(filterFunc)
}

func (c *Context) FindFilter(name string) *filter.Filter {
	return c.filters[name]
}
