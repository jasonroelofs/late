package tag

import (
	"github.com/jasonroelofs/late/object"
)

type Include struct{}

func (i *Include) Parse() *ParseConfig {
	return &ParseConfig{
		TagName: "include",
		Rules:   []ParseRule{Expression()},
	}
}

func (i *Include) Eval(env Environment, results *ParseResult) object.Object {
	// TODO: Make sure we have a string here
	// something like,
	//
	//   paritalName, err := object.GetString(result.Nodes[0])
	//   paritalName, err := object.Get(result.Nodes[0], object.TYPE_STRING)
	//

	partialName := results.Nodes[0].Value().(string)
	partialBody := env.ReadFile(partialName)

	env.PushScope()

	// Set up a new context stack?
	result := env.Render(partialBody)

	env.PopScope()

	return object.New(result)
}
