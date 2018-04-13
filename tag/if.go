package tag

import (
	"strings"

	"github.com/jasonroelofs/late/object"
)

/**
 * The veritable `if` tag!
 */
type If struct {
}

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
	initialValue := results.Nodes[0] == object.TRUE

	if initialValue {
		return i.evalStatements(env, results.Statements)
	} else {
		for _, subTag := range results.SubTagResults {
			if (subTag.TagName == "elsif" && subTag.Nodes[0] == object.TRUE) ||
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
