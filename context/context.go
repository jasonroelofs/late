package context

import (
	"github.com/jasonroelofs/late/object"
	s "github.com/jasonroelofs/late/template/statement"
)

type Assigns map[string]interface{}

/**
 * Evaluator is the set of functions that the Context needs to be
 * able to delegate down to evaluation when processing tags.
 */
type Evaluator interface {
	Eval(s.Statement) object.Object
	EvalAll([]s.Statement) object.Object
	Interrupt() string
	ClearInterrupt()
}

type Context struct {
	RenderFunc func(string, *Context) string

	evaluator    Evaluator
	globalScope  *Scope
	currentScope *Scope
	reader       FileReader
}

func New(options ...func(*Context)) *Context {
	ctx := &Context{
		globalScope: NewScope(nil),
		reader:      new(NullReader),
	}

	ctx.currentScope = ctx.globalScope

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

func (c *Context) SetEvaluator(e Evaluator) {
	c.evaluator = e
}

func (c *Context) Assign(assigns Assigns) {
	for key, value := range assigns {
		c.Set(key, value)
	}
}

/**
 * Store a value under a given name for the current scope.
 */
func (c *Context) Set(name string, value interface{}) {
	c.currentScope.Set(name, object.New(value))
}

func (c *Context) ShadowSet(name string, value interface{}) {
	c.currentScope.ShadowSet(name, object.New(value))
}

func (c *Context) Get(name string) object.Object {
	return c.currentScope.Get(name)
}

/**
 * Promoting a variable takes it out of the current scope
 * and moves it up to live at the global scope. This then ensures
 * that an include at any level of nesting is able to promote
 * a variable to allow any template or partial in the current render
 * to access that variable.
 */
func (c *Context) Promote(name string) {
	c.globalScope.Set(name, c.currentScope.Get(name))
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

func (c *Context) Render(input string) object.Object {
	if c.RenderFunc == nil {
		return object.NULL
	}

	return object.New(c.RenderFunc(input, c))
}

func (c *Context) Eval(s s.Statement) object.Object {
	return c.evaluator.Eval(s)
}

func (c *Context) EvalAll(stmts []s.Statement) object.Object {
	return c.evaluator.EvalAll(stmts)
}

func (c *Context) Interrupt() string {
	return c.evaluator.Interrupt()
}

func (c *Context) ClearInterrupt() {
	c.evaluator.ClearInterrupt()
}
