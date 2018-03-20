package template

import (
	"fmt"
	"strings"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/template/evaluator"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/parser"
)

type Template struct {
	body string
}

func New(templateBody string) *Template {
	return &Template{
		body: templateBody,
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render(context *context.Context) string {
	lexer := lexer.New(t.body)
	parser := parser.New(lexer)
	ast := parser.Parse()

	if len(parser.Errors) > 0 {
		fmt.Println("There were errors parsing the document")

		for _, msg := range parser.Errors {
			fmt.Println(msg)
		}

		// For now, just return the original document.
		return t.body
	}

	eval := evaluator.New(ast, context)

	final := strings.Builder{}

	for _, obj := range eval.Run() {
		final.WriteString(obj.Inspect())
	}

	return final.String()
}
