package parser

import (
	"testing"

	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
)

func TestRawTemplates(t *testing.T) {
	input := "This is a raw template\nIt has no liquid code whatsoever"

	l := lexer.New(input)
	p := New(l)

	template := p.Parse()
	checkParserErrors(t, p)

	if template == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(template.Statements) != 1 {
		t.Fatalf("Parsing built the wrong number of nodes. got=%d", len(template.Statements))
	}

	stmt := template.Statements[0]

	if stmt.String() != "This is a raw template\nIt has no liquid code whatsoever" {
		t.Fatalf("Raw statement not correct. got='%s'", stmt.String())
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input    string
		errorStr string
	}{
		{"{{ %}", "expected IDENT or CLOSE_VAR or NUMBER or STRING, found CLOSE_TAG"},
		{"{{", "expected IDENT or CLOSE_VAR or NUMBER or STRING, found EOF"},
		{"{{ foobar ", "expected CLOSE_VAR, found EOF"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		p.Parse()

		if len(p.Errors) != 1 {
			t.Fatalf("Parser did not find errors when it should have")
		}

		if p.Errors[0] != test.errorStr {
			t.Fatalf("Wrong error. Wanted: \"%s\" Got: \"%s\"", test.errorStr, p.Errors[0])
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "{{ te }}"

	l := lexer.New(input)
	p := New(l)

	template := p.Parse()
	checkParserErrors(t, p)

	if len(template.Statements) != 1 {
		t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
	}

	stmt, ok := template.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression not an IDENT, got %T", stmt.Expression)
	}

	if ident.Value != "te" {
		t.Fatalf("Identifier has the wrong name, got %s", ident.Value)
	}
}

func TestNumberExpression(t *testing.T) {
	input := "{{ 400 }} {{ 3.1415 }}"
	// TODO test invalid numbers

	l := lexer.New(input)
	p := New(l)

	template := p.Parse()
	checkParserErrors(t, p)

	if len(template.Statements) != 3 {
		t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
	}

	stmt, ok := template.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("Expression not a NUMBER, got %T", stmt.Expression)
	}

	if ident.Value != 400 {
		t.Fatalf("Identifier has the wrong name, got %f", ident.Value)
	}

	// The 2nd statement is a Raw node between the two

	stmt, ok = template.Statements[2].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	ident, ok = stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("Expression not a NUMBER, got %T", stmt.Expression)
	}

	if ident.Value != 3.1415 {
		t.Fatalf("Identifier has the wrong name, got %f", ident.Value)
	}
}

//func TestBasicLiquid(t *testing.T) {
//	input := `This is the header {{ "In Liquid" }} This is the footer`
//
//	l := lexer.New(input)
//	p := New(l)
//
//	template := p.Parse()
//	checkParserErrors(t, p)
//
//	if template == nil {
//		t.Fatalf("Parse() returned nil")
//	}
//
//	if len(template.Statements) != 3 {
//		t.Fatalf("Parsing built the wrong number of nodes. got=%d", len(template.Statements))
//	}
//
//	stmt := template.Statements[0]
//
//	if stmt.String() != "This is the header " {
//		t.Fatalf("First statement not correct. got='%s'", stmt.String())
//	}
//
//	stmt = template.Statements[1]
//
//	if stmt.String() != `"In Liquid"` {
//		t.Fatalf("Second statement not correct. got='%s'", stmt.String())
//	}
//
//	stmt = template.Statements[2]
//
//	if stmt.String() != " This is the footer" {
//		t.Fatalf("Third statement not correct. got='%s'", stmt.String())
//	}
//
//	program := template.String()
//	if program != input {
//		t.Fatalf("The rendered result didn't match the input. got '%s'", program)
//	}
//}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error %q", msg)
	}

	t.FailNow()
}
