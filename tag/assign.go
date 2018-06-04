package tag

import (
	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/template/token"
)

/**
 * The assign tag lets users set new variables in their templates.
 * These variables are initially set in Template Only scope, but can be promoted
 * to Render scope via the `promote` tag.
 *
 *   {% assign var_name = EXPRESSION %}
 *
 */
type Assign struct {
}

func (a *Assign) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "assign",
		Rules:   []ParseRule{Identifier(), Token(token.ASSIGN), Expression()},
	}
}

func (a *Assign) Eval(ctx *context.Context, results *ParseResult) object.Object {
	varName := results.Nodes[0].Value().(string)
	result := results.Nodes[2]

	ctx.Set(varName, result)
	return object.NULL
}
