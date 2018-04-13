package tag

import (
	"strings"

	"github.com/jasonroelofs/late/object"
)

/**
 * The veritable `if` tag!
 */
type If struct{}

func (i *If) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "if",
		Block:   true,
		Rules:   []ParseRule{Expression()},
		SubTags: []ParseConfig{
			{
				TagName: "elsif",
				Block:   true,
				Rules:   []ParseRule{Expression()},
			},
			{
				TagName: "else",
				Block:   true,
			},
		},
	}
}

func (i *If) Eval(env Environment, results *ParseResult) object.Object {
	if object.Truthy(results.Nodes[0]) {
		return i.evalStatements(env, results.Statements)
	} else {
		for _, subTag := range results.SubTagResults {
			if (subTag.TagName == "elsif" && object.Truthy(subTag.Nodes[0])) ||
				subTag.TagName == "else" {
				return i.evalStatements(env, subTag.Statements)
			}
		}
	}

	return object.NULL
}

func (i *If) evalStatements(env Environment, statements []Statement) object.Object {
	content := strings.Builder{}

	for _, stmt := range statements {
		content.WriteString(env.Eval(stmt).Inspect())
	}

	return object.New(content.String())
}
