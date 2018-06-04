package tag

import (
	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
)

/**
 * The capture block evaluates all code in its block and assigns the result to
 * a variable usable elsewhere.
 *
 *   {% capture header %}
 *     <title>{{ site_title }}</title>
 *   {% end %}
 *
 *   {{ header }}
 *
 */
type Capture struct {
}

func (c *Capture) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "capture",
		Block:   true,
		Rules:   []ParseRule{Identifier()},
	}
}

func (c *Capture) Eval(ctx *context.Context, results *ParseResult) object.Object {
	varName := results.Nodes[0].Value().(string)

	ctx.Set(varName, ctx.EvalAll(results.Statements))
	return object.NULL
}
