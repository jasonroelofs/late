package context

import (
	"github.com/jasonroelofs/late/object"
)

type Assigns map[string]interface{}

type Context struct {
	RenderFunc func(string, *Context) string

	currentScope *Scope
	reader       FileReader
}

func New(options ...func(*Context)) *Context {
	ctx := &Context{
		currentScope: NewScope(nil),
		reader:       new(NullReader),
	}

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

func (c *Context) Assign(assigns Assigns) {
	for key, value := range assigns {
		c.Set(key, value)
	}
}

func (c *Context) Set(name string, value interface{}) {
	c.currentScope.Set(name, object.New(value))
}

func (c *Context) ShadowSet(name string, value interface{}) {
	c.currentScope.ShadowSet(name, object.New(value))
}

func (c *Context) Get(name string) object.Object {
	return c.currentScope.Get(name)
}

func (c *Context) Promote(name string) {
	c.currentScope.Promote(name)
}

func (c *Context) PushScope() {
	c.currentScope = NewScope(c.currentScope)
}

func (c *Context) PushShadowScope() {
	c.currentScope = NewShadowScope(c.currentScope)
}

func (c *Context) PopScope() {
	if c.currentScope.Parent == nil {
		return
	}

	c.currentScope = c.currentScope.Parent
}

func (c *Context) ReadFile(path string) string {
	return c.reader.Read(path)
}
