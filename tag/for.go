package tag

import (
	"strings"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
)

var (
	FIRST  = object.New("first")
	LAST   = object.New("last")
	INDEX  = object.New("index")
	LENGTH = object.New("length")
)

/**
 * The for loop
 */
type For struct{}

func (f *For) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "for",
		Block:   true,
		Rules:   []ParseRule{Identifier(), Literal("in"), Expression()},
	}
}

func (f *For) Eval(ctx *context.Context, results *ParseResult) object.Object {
	varName := results.Nodes[0].Value().(string)

	// TODO: Error if collection isn't an Array
	// Also TODO: Support iteration over a Hash
	collection := results.Nodes[2].(*object.Array)
	output := strings.Builder{}

	// Set up our shadow scope that keeps `forloop` and the loop variable
	// scoped to this for loop but allows users to assign values to the template's
	// scope.
	ctx.PushShadowScope()

	forLoopInfo := object.NewHash()
	forLoopInfo.Set(LENGTH, object.New(len(collection.Elements)))
	ctx.ShadowSet("forloop", forLoopInfo)

loop:
	for idx, entry := range collection.Elements {
		ctx.ShadowSet(varName, entry)

		forLoopInfo.Set(INDEX, object.New(idx))
		forLoopInfo.Set(FIRST, object.New(idx == 0))
		forLoopInfo.Set(LAST, object.New(idx == len(collection.Elements)-1))

		output.WriteString(ctx.EvalAll(results.Statements).Inspect())

		switch ctx.Interrupt() {
		case "continue":
			ctx.ClearInterrupt()
			continue loop
		case "break":
			ctx.ClearInterrupt()
			break loop
		}
	}

	ctx.PopScope()

	return object.New(output.String())
}

/**
 * Define the two Interrupts that we need to handle
 */

type Continue struct{}

func (c *Continue) Parse() *ParseConfig {
	return &ParseConfig{TagName: "continue", Interrupt: true}
}

func (c *Continue) Eval(_ *context.Context, _ *ParseResult) object.Object {
	return object.NULL
}

type Break struct{}

func (b *Break) Parse() *ParseConfig {
	return &ParseConfig{TagName: "break", Interrupt: true}
}

func (b *Break) Eval(_ *context.Context, _ *ParseResult) object.Object {
	return object.NULL
}
