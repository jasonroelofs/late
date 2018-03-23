package tag

import (
	"github.com/jasonroelofs/late/object"
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

func (a *Assign) TagName() string {
	return "assign"
}

func (a *Assign) Parse() []ParseRule {
	return []ParseRule{Identifier(), Literal("="), Expression()}
}

func (a *Assign) Eval(env Environment, results []object.Object) {
	varName := results[0].Value().(string)
	result := results[2]

	env.Set(varName, result)
}
