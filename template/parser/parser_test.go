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

	if len(template.Nodes) != 1 {
		t.Fatalf("Parsing built the wrong number of nodes. got=%d", len(template.Nodes))
	}

	stmt := template.Nodes[0]

	if stmt.String() != "This is a raw template\nIt has no liquid code whatsoever" {
		t.Fatalf("Raw statement not correct. got='%s'", stmt.String())
	}
}

func TestBasicLiquid(t *testing.T) {
	input := `This is the header {{ "In Liquid" }} This is the footer`

	l := lexer.New(input)
	p := New(l)

	template := p.Parse()

	if template == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(template.Nodes) != 3 {
		t.Fatalf("Parsing built the wrong number of nodes. got=%d", len(template.Nodes))
	}

	stmt := template.Nodes[0]

	if stmt.String() != "This is the header " {
		t.Fatalf("First statement not correct. got='%s'", stmt.String())
	}

	stmt = template.Nodes[1]

	if stmt.String() != `"In Liquid"` {
		t.Fatalf("Second statement not correct. got='%s'", stmt.String())
	}

	stmt = template.Nodes[2]

	if stmt.String() != " This is the footer" {
		t.Fatalf("Third statement not correct. got='%s'", stmt.String())
	}

	program := template.String()
	if program != input {
		t.Fatalf("The rendered result didn't match the input. got '%s'", program)
	}
}
