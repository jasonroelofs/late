package tag

import (
	"github.com/jasonroelofs/late/context"
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

func (i *If) Eval(ctx *context.Context, results *ParseResult) object.Object {
	if object.Truthy(results.Nodes[0]) {
		return ctx.EvalAll(results.Statements)
	} else {
		for _, subTag := range results.SubTagResults {
			if (subTag.TagName == "elsif" && object.Truthy(subTag.Nodes[0])) ||
				subTag.TagName == "else" {
				return ctx.EvalAll(subTag.Statements)
			}
		}
	}

	return object.NULL
}
