package context

import (
	"github.com/jasonroelofs/late/filter"
	"github.com/jasonroelofs/late/object"
)

type Assigns map[string]interface{}

type Context struct {
	filters map[string]*filter.Filter

	globalAssigns map[string]object.Object
}

func New() *Context {
	ctx := &Context{
		filters:       make(map[string]*filter.Filter),
		globalAssigns: make(map[string]object.Object),
	}

	filter.ApplyStandardFilters(ctx)

	return ctx
}

func (c *Context) Assign(assigns Assigns) {
	for key, value := range assigns {
		c.Set(key, value)
	}
}

func (c *Context) Set(name string, value interface{}) {
	c.globalAssigns[name] = object.New(value)
}

func (c *Context) Get(name string) object.Object {
	obj, ok := c.globalAssigns[name]
	if !ok {
		// WARN: Referenced undefined variable {name}
		return object.NULL
	}

	return obj
}

func (c *Context) AddFilter(name string, filterFunc filter.FilterFunc) {
	c.filters[name] = filter.New(filterFunc)
}

func (c *Context) FindFilter(name string) *filter.Filter {
	return c.filters[name]
}
