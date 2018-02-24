package parser

import (
	"testing"

	"github.com/jasonroelofs/late/template/lexer"
)

func TestRawTemplates(t *testing.T) {
	input := "This is a raw template\nIt has no liquid code whatsoever"

	l := lexer.New(input)
	p := New(l)

	template := p.Parse()

	if template == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(template.Statements) != 1 {
		t.Fatalf("Parsing built the wrong number of nodes. got=%d", len(template.Statements))
	}

	stmt := template.Statements[0]

	if stmt.TokenLiteral() != "This is a raw template\nIt has no liquid code whatsoever" {
		t.Fatalf("Raw statement not correct. got='%s'", stmt.TokenLiteral())
	}
}
