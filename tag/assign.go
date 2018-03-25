package tag

import (
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

func (a *Assign) Block() bool {
	return false
}

func (a *Assign) Parse() []ParseRule {
	return []ParseRule{Identifier(), Token(token.ASSIGN), Expression()}
}

func (a *Assign) Eval(env Environment, results []object.Object, _ []Statement) object.Object {
	varName := results[0].Value().(string)
	result := results[2]

	env.Set(varName, result)
	return object.NULL
}
