package context

import (
	"github.com/jasonroelofs/late/object"
)

type Assigns map[string]interface{}

type Context struct {
	globalAssigns map[string]object.Object
}

func New() *Context {
	ctx := &Context{
		globalAssigns: make(map[string]object.Object),
	}

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
