package tag

import (
	"strings"

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

func (f *For) Eval(env Environment, results *ParseResult) object.Object {
	varName := results.Nodes[0].Value().(string)

	// TODO: Error if collection isn't an Array
	// Also TODO: Support iteration over a Hash
	collection := results.Nodes[2].(*object.Array)
	output := strings.Builder{}

	// Set up our shadow scope that keeps `forloop` and the loop variable
	// scoped to this for loop but allows users to assign values to the template's
	// scope.
	env.PushShadowScope()

	forLoopInfo := object.NewHash()
	forLoopInfo.Set(LENGTH, object.New(len(collection.Elements)))
	env.ShadowSet("forloop", forLoopInfo)

loop:
	for idx, entry := range collection.Elements {
		env.ShadowSet(varName, entry)

		forLoopInfo.Set(INDEX, object.New(idx))
		forLoopInfo.Set(FIRST, object.New(idx == 0))
		forLoopInfo.Set(LAST, object.New(idx == len(collection.Elements)-1))

		output.WriteString(env.EvalAll(results.Statements).Inspect())

		switch env.Interrupt() {
		case "continue":
			env.ClearInterrupt()
			continue loop
		case "break":
			env.ClearInterrupt()
			break loop
		}
	}

	env.PopScope()

	return object.New(output.String())
}

/**
 * Define the two Interrupts that we need to handle
 */

type Continue struct{}

func (c *Continue) Parse() *ParseConfig {
	return &ParseConfig{TagName: "continue", Interrupt: true}
}

func (c *Continue) Eval(_ Environment, _ *ParseResult) object.Object {
	return object.NULL
}

type Break struct{}

func (b *Break) Parse() *ParseConfig {
	return &ParseConfig{TagName: "break", Interrupt: true}
}

func (b *Break) Eval(_ Environment, _ *ParseResult) object.Object {
	return object.NULL
}
