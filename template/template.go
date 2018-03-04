package template

import (
	"strings"

	"github.com/jasonroelofs/late/template/evaluator"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/parser"
)

type Template struct {
	body string
}

type Params map[string]string

func New(templateBody string) *Template {
	return &Template{
		body: templateBody,
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render() string {
	lexer := lexer.New(t.body)
	parser := parser.New(lexer)
	ast := parser.Parse()
	eval := evaluator.New(ast)

	final := strings.Builder{}

	for _, obj := range eval.Run() {
		final.WriteString(obj.Inspect())
	}

	return final.String()
}
