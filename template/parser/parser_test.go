package parser

import (
	"testing"

	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
)

func TestRawTemplates(t *testing.T) {
	input := "This is a raw template\nIt has no liquid code whatsoever"

	template := parseTest(t, input)

	if template == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(template.Statements) != 1 {
		t.Fatalf("Parsing built the wrong number of statements. got=%d", len(template.Statements))
	}

	stmt := template.Statements[0]

	if stmt.String() != "This is a raw template\nIt has no liquid code whatsoever" {
		t.Fatalf("Raw statement not correct. got='%s'", stmt.String())
	}
}

// Figure this out properly
//func TestParserErrors(t *testing.T) {
//	tests := []struct {
//		input    string
//		errorStr string
//	}{
//		{"{{ %}", "expected CLOSE_VAR, found EOF"},
//		{"{{", "expected CLOSE_VAR, found EOF"},
//		{"{{ foobar ", "expected CLOSE_VAR, found EOF"},
//	}
//
//	for i, test := range tests {
//		l := lexer.New(test.input)
//		p := New(l)
//		p.Parse()
//
//		if len(p.Errors) != 1 {
//			fmt.Printf("Parser errors: %#v", p.Errors)
//			t.Fatalf("(%d) Parser didn't find right # of errors. Found %d", i, len(p.Errors))
//		}
//
//		if p.Errors[0] != test.errorStr {
//			t.Fatalf("(%d) Wrong error. Wanted: \"%s\" Got: \"%s\"", i, test.errorStr, p.Errors[0])
//		}
//	}
//}

func TestIdentifierExpression(t *testing.T) {
	input := "{{ te }}"

	template := parseTest(t, input)

	if len(template.Statements) != 1 {
		t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
	}

	stmt, ok := template.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	checkIdentifierExpression(t, stmt.Expression, "te")
}

func TestNumberLiteral(t *testing.T) {
	input := "{{ 400 }} {{ 3.1415 }}"
	// TODO test invalid numbers

	template := parseTest(t, input)

	if len(template.Statements) != 3 {
		t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
	}

	stmt, ok := template.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	checkNumberExpression(t, stmt.Expression, 400)

	// The 2nd statement is a Raw node between the two

	stmt, ok = template.Statements[2].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	checkNumberExpression(t, stmt.Expression, 3.1415)
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"{{ true }}", true},
		{"{{ false }}", false},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)

		if len(template.Statements) != 1 {
			t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
		}

		stmt, ok := template.Statements[0].(*ast.VariableStatement)
		if !ok {
			t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("stmt is not a BooleanLiteral, got %T", stmt.Expression)
		}

		if exp.Value != test.expected {
			t.Fatalf("BooleanLiteral has wrong value, expected %t got %t", test.expected, exp.Value)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ "A string" }}`, "A string"},
		{`{{ 'Single Quotes' }}`, "Single Quotes"},
		{`{{ "Mixe'd Quotes" }}`, "Mixe'd Quotes"},
		{`{{ 'Escape\'d Quotes' }}`, "Escape\\'d Quotes"},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)

		if len(template.Statements) != 1 {
			t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
		}

		stmt, ok := template.Statements[0].(*ast.VariableStatement)
		if !ok {
			t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("stmt is not a StringLiteral, got %T", stmt.Expression)
		}

		if exp.Value != test.expected {
			t.Fatalf("StringLiteral has wrong value, expected %s got %s", test.expected, exp.Value)
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		expected string
	}{
		{"{{ -15 }}", "-", "(-15)"},
		{"{{ -te }}", "-", "(-te)"},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)

		if len(template.Statements) != 1 {
			t.Fatalf("Template did not have the right number of statements, got %d", len(template.Statements))
		}

		stmt, ok := template.Statements[0].(*ast.VariableStatement)
		if !ok {
			t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not a PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != test.operator {
			t.Fatalf("Wrong operator, expected '%s' got '%s'", test.operator, exp.Operator)
		}

		if exp.String() != test.expected {
			t.Fatalf("Expression incorrect. Expected '%s' got '%s'", test.expected, exp.String())
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  float64
		operator   string
		rightValue float64
	}{
		{"{{ 1 + 1 }}", 1, "+", 1},
		{"{{ 2 - 2 }}", 2, "-", 2},
		{"{{ 3 * 3 }}", 3, "*", 3},
		{"{{ 4 / 4 }}", 4, "/", 4},
		{"{{ 5 < 5 }}", 5, "<", 5},
		{"{{ 5 <= 5 }}", 5, "<=", 5},
		{"{{ 6 > 6 }}", 6, ">", 6},
		{"{{ 6 >= 6 }}", 6, ">=", 6},
		{"{{ 7 == 7 }}", 7, "==", 7},
		{"{{ 8 != 8 }}", 8, "!=", 8},
	}

	for i, test := range tests {
		template := parseTest(t, test.input)

		if len(template.Statements) != 1 {
			t.Fatalf("(%d) Template did not have the right number of statements, got %d", i, len(template.Statements))
		}

		stmt, ok := template.Statements[0].(*ast.VariableStatement)
		if !ok {
			t.Fatalf("(%d) Template statement was the wrong type, got %T", i, template.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("(%d) stmt is not an InfixExpression, got %T", i, stmt.Expression)
		}

		checkNumberExpression(t, exp.Left, test.leftValue)

		if exp.Operator != test.operator {
			t.Fatalf("Operator was wrong, expected '%s' got '%s'", test.operator, exp.Operator)
		}

		checkNumberExpression(t, exp.Right, test.rightValue)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{{ -a * b }}", "((-a) * b)"},
		{"{{ a + b - c }}", "((a + b) - c)"},
		{"{{ a * b / c }}", "((a * b) / c)"},
		{"{{ a + b * c }}", "(a + (b * c))"},
		{"{{ 5 > 4 == 3 < 4 }}", "((5 > 4) == (3 < 4))"},
		{"{{ 3 + 2 * 6 != 12 <= 100 / 2 }}", "((3 + (2 * 6)) != (12 <= (100 / 2)))"},
		// Explicit grouping
		{"{{ (3 + 2) * 6 }}", "((3 + 2) * 6)"},
		{"{{ a + (b - c) * c}}", "(a + ((b - c) * c))"},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)

		got := template.String()
		if got != test.expected {
			t.Errorf("Precedence result wrong. Expected '%s', got '%s'", test.expected, got)
		}
	}
}

func parseTest(t *testing.T, input string) *ast.Template {
	l := lexer.New(input)
	p := New(l)

	template := p.Parse()
	checkParserErrors(t, p)

	return template
}

func checkNumberExpression(t *testing.T, exp ast.Expression, expected float64) {
	number, ok := exp.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("Expression not NUMBER, got %T", exp)
	}

	if number.Value != expected {
		t.Fatalf("Number has the wrong value. Expected '%f' got '%f'", expected, number.Value)
	}
}

func checkIdentifierExpression(t *testing.T, exp ast.Expression, expected string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression not IDENT, got %T", exp)
	}

	if ident.Value != expected {
		t.Fatalf("Identifier has the wrong name, expected '%s' got '%s'", expected, ident.Value)
	}
}

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
