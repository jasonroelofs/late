package parser

import (
	"testing"

	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
)

func TestRawTemplates(t *testing.T) {
	input := "This is a raw template\nIt has no liquid code whatsoever"

	template := parseTest(t, input)
	checkStatementCount(t, template, 1)

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
	checkStatementCount(t, template, 1)

	stmt := getVariableStatement(t, template, 0)
	checkIdentifierExpression(t, stmt.Expression, "te")
}

func TestNumberLiteral(t *testing.T) {
	input := "{{ 400 }} {{ 3.1415 }}"
	// TODO test invalid numbers

	template := parseTest(t, input)
	checkStatementCount(t, template, 3)

	stmt := getVariableStatement(t, template, 0)
	checkNumberExpression(t, stmt.Expression, 400)

	// The 2nd statement is a Raw node between the two

	stmt = getVariableStatement(t, template, 2)
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
		checkStatementCount(t, template, 1)

		stmt := getVariableStatement(t, template, 0)
		checkBooleanLiteral(t, stmt.Expression, test.expected)
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
		checkStatementCount(t, template, 1)

		stmt := getVariableStatement(t, template, 0)
		checkStringLiteral(t, stmt.Expression, test.expected)
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
		checkStatementCount(t, template, 1)

		stmt := getVariableStatement(t, template, 0)

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
		checkStatementCount(t, template, 1)

		stmt := getVariableStatement(t, template, 0)

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
		// PIPE needs to be super low
		{"{{ a | size }}", "(a | size)"},
		{"{{ b | upcase | size }}", "((b | upcase) | size)"},
		{"{{ 5 * 6 + 1 | filter }}", "(((5 * 6) + 1) | filter)"},
		{`{{ a | remove: "this" | upcase }}`, `((a | (remove: "this")) | upcase)`},
		{`{{ a | replace: "this", with: "that" | upcase }}`, `((a | (replace: "this", with: "that")) | upcase)`},
		// We can explicitly nest filters inside of expressions with grouping!
		{`{{ a | replace: "this", with: ("that" | upcase) }}`, `(a | (replace: "this", with: ("that" | upcase)))`},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)

		got := template.String()
		if got != test.expected {
			t.Errorf("Precedence result wrong. Expected '%s', got '%s'", test.expected, got)
		}
	}
}

func TestFilters(t *testing.T) {
	tests := []struct {
		input          string
		expectedVar    string
		operator       string
		expectedFilter string
	}{
		{`{{ "A String" | upcase }}`, "A String", "|", "upcase"},
		{`{{ "A String" | size }}`, "A String", "|", "size"},
	}

	for i, test := range tests {
		template := parseTest(t, test.input)
		checkStatementCount(t, template, 1)

		stmt := getVariableStatement(t, template, 0)

		exp, ok := stmt.Expression.(*ast.FilterExpression)
		if !ok {
			t.Fatalf("(%d) stmt is not a FilterExpression, got %T", i, stmt.Expression)
		}

		checkStringLiteral(t, exp.Input, test.expectedVar)

		checkFilterLiteral(t, exp.Filter, test.expectedFilter)
	}
}

func TestFiltersWithParameters(t *testing.T) {
	input := `{{ "Hello Mom" | replace: "Mom", with: "World" }}`
	template := parseTest(t, input)
	checkStatementCount(t, template, 1)

	stmt := getVariableStatement(t, template, 0)

	exp, ok := stmt.Expression.(*ast.FilterExpression)
	if !ok {
		t.Fatalf("stmt is not a FilterExpression, got %T", stmt.Expression)
	}

	checkFilterLiteral(t, exp.Filter, "replace")
	filter := exp.Filter.(*ast.FilterLiteral)

	params := filter.Parameters
	if len(params) != 2 {
		t.Fatalf("Wrong number of parameters, got %d", len(params))
	}

	expr, ok := params["replace"]
	if !ok {
		t.Fatalf("Didn't set the initial `replace` parameter")
	}

	checkStringLiteral(t, expr, "Mom")

	expr, ok = params["with"]
	if !ok {
		t.Fatalf("Didn't set the explicit `with` parameter")
	}

	checkStringLiteral(t, expr, "World")
}

func TestTags(t *testing.T) {
	tests := []struct {
		input    string
		tagName  string
		numNodes int
	}{
		{`{% assign this = "that" %}`, "assign", 3},
	}

	for _, test := range tests {
		template := parseTest(t, test.input)
		checkStatementCount(t, template, 1)

		stmt := getTagStatement(t, template, 0)

		if stmt.TagName != test.tagName {
			t.Fatalf("Did not parse out the right tag name, Expected %s Got %s", test.tagName, stmt.TagName)
		}

		if stmt.Tag == nil {
			t.Fatalf("Did not store the instantiated tag in the tree")
		}

		if len(stmt.Nodes) != test.numNodes {
			t.Fatalf("Did not store the right number of ast nodes. Expected %d got %d", test.numNodes, len(stmt.Nodes))
		}
	}
}

/**
 * Helper methods
 */

func parseTest(t *testing.T, input string) *ast.Template {
	l := lexer.New(input)
	p := New(l)

	template := p.Parse()
	checkParserErrors(t, p)

	return template
}

func checkStatementCount(t *testing.T, template *ast.Template, expected int) {
	if len(template.Statements) != expected {
		t.Fatalf("Wrong number of statements. Expected %d, Got %d", expected, len(template.Statements))
	}
}

func checkBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) {
	boolean, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Expression not BOOLEAN, got %T", exp)
	}

	if boolean.Value != expected {
		t.Fatalf("Bool has the wrong value. Expected '%t' got '%t'", expected, boolean.Value)
	}
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

func checkStringLiteral(t *testing.T, exp ast.Expression, expected string) {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expression not STRING, got %T", exp)
	}

	if str.Value != expected {
		t.Fatalf("String has the wrong value. Expected '%s', got '%s'", expected, str.Value)
	}
}

func checkFilterLiteral(t *testing.T, exp ast.Expression, expected string) {
	filter, ok := exp.(*ast.FilterLiteral)
	if !ok {
		t.Fatalf("Expression not FILTER, got %T", exp)
	}

	if filter.Name != expected {
		t.Fatalf("Filter has the wrong name, expected '%s' got '%s'", expected, filter.Name)
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

func getVariableStatement(t *testing.T, template *ast.Template, index int) *ast.VariableStatement {
	stmt, ok := template.Statements[index].(*ast.VariableStatement)

	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	return stmt
}

func getTagStatement(t *testing.T, template *ast.Template, index int) *ast.TagStatement {
	stmt, ok := template.Statements[index].(*ast.TagStatement)

	if !ok {
		t.Fatalf("Template statement was the wrong type, got %T", template.Statements[0])
	}

	return stmt
}
