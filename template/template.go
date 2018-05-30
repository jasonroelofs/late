package template

import (
	"strings"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/template/evaluator"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/parser"
)

type Template struct {
	body string

	Errors []string
}

func New(templateBody string) *Template {
	return &Template{
		body: templateBody,
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render(ctx *context.Context) string {

	// Setup ourselves as the re-entrant render function for this context
	// This is to ensure that the include tag (and anything else that wants to trigger
	// a full new render stack) doesn't need to depend on template, thus causing
	// an import cycle.
	ctx.RenderFunc = func(body string, ctx *context.Context) string {
		tpl := New(body)
		// TODO Propogating errors back up the stack
		return tpl.Render(ctx)
	}

	lexer := lexer.New(t.body)
	parser := parser.New(lexer)
	ast := parser.Parse()

	if len(parser.Errors) > 0 {
		t.Errors = parser.Errors
		// For now, just return the original document.
		return t.body
	}

	eval := evaluator.New(ast, ctx)
	final := strings.Builder{}
	results := eval.Run()

	for _, obj := range results {
		final.WriteString(obj.Inspect())
	}

	return final.String()
}
