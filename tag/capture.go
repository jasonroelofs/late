package tag

import (
	"strings"

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

func (c *Capture) Eval(env Environment, results *ParseResult) object.Object {
	content := strings.Builder{}
	varName := results.Nodes[0].Value().(string)

	for _, stmt := range results.Statements {
		content.WriteString(env.Eval(stmt).Inspect())
	}

	env.Set(varName, content.String())
	return object.NULL
}
