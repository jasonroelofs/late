package tag

import (
	"strings"

	"github.com/jasonroelofs/late/object"
)

/**
 * The raw block will output everything in its block as it is without
 * any parsing or evaluation
 *
 *   {% raw %}
 *   	{{ "This is a variable" }}
 *   {% end %}
 *
 */
type Raw struct {
}

func (a *Raw) Block() bool {
	return true
}

func (a *Raw) Parse() []ParseRule {
	return []ParseRule{}
}

func (a *Raw) Eval(env Environment, _ []object.Object, statements []Statement) object.Object {
	out := &strings.Builder{}

	for _, stmt := range statements {
		out.WriteString(stmt.String())
	}

	return object.New(out.String())
}
