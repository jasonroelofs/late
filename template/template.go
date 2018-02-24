package template

import (
	"github.com/jasonroelofs/late/parser"
)

type Template struct {
	body string

	parser *parser.Parser
}

type Params map[string]string

func New(templateBody string) *Template {
	return &Template{
		body:   templateBody,
		parser: parser.New(templateBody),
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render() string {
	parsedTemplate := t.parser.Parse()

	return parsedTemplate
}
