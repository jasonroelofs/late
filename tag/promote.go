package tag

import (
	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
)

/**
 * The promote tag pushes the given variable up the scope tree, making it
 * available in other templates.
 * Mainly for use when using `include` to make variables set in partials available
 * to the template that called it.
 *
 *   {% promote variable %}
 *
 */
type Promote struct {
}

func (p *Promote) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "promote",
		Rules:   []ParseRule{Identifier()},
	}
}

func (p *Promote) Eval(ctx *context.Context, results *ParseResult) object.Object {
	varName := results.Nodes[0].Value().(string)

	ctx.Promote(varName)

	return object.NULL
}
