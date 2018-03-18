package evaluator

import (
	"testing"

	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/parser"
)

func TestRawStatements(t *testing.T) {
	input := "This is a raw, non-liquid template"

	results := evalInput(t, input)

	checkStatementCount(t, results, 1)
	checkObject(t, results[0], object.OBJ_STRING, input)
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"{{ 1 }}", 1},
		{"{{ 1 + 1 }}", 2},
		{"{{ 2 * 5 - 1 }}", 9},
		{"{{ 2 + 3 * 5 }}", 17},
		{"{{ (2 + 3) * 5 }}", 25},
		{"{{ -1 }}", -1},
		{"{{ -(12 + 3) }}", -15},
	}

	for _, test := range tests {
		results := evalInput(t, test.input)

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.OBJ_NUMBER, test.expected)
	}
}

func TestBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"{{ true }}", true},
		{"{{ false }}", false},
		{"{{ true == true }}", true},
		{"{{ true == false }}", false},
		{"{{ 1 < 2 }}", true},
		{"{{ 2 > 5 }}", false},
		{"{{ 3 <= 3 }}", true},
		{"{{ 4 >= 5 }}", false},
	}

	for _, test := range tests {
		results := evalInput(t, test.input)

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.OBJ_BOOL, test.expected)
	}
}

func TestStrings(t *testing.T) {
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
		results := evalInput(t, test.input)

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.OBJ_STRING, test.expected)
	}
}

func TestFilters(t *testing.T) {
	tests := []struct {
		input        string
		expectedType object.ObjectType
		expected     interface{}
	}{
		{`{{ "A String" | size }}`, object.OBJ_NUMBER, float64(8)},
		{`{{ "A String" | upcase }}`, object.OBJ_STRING, "A STRING"},
		{`{{ "A String" | upcase | size }}`, object.OBJ_NUMBER, float64(8)},
		{`{{ "Hello Mom" | replace: "Mom", with: "World" }}`, object.OBJ_STRING, "Hello World"},
		// TODO: Unknown filter
		//   Strict: error out
		//   Lax: treat as a pass-through no-op, trigger a warning
	}

	for _, test := range tests {
		results := evalInput(t, test.input)

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], test.expectedType, test.expected)
	}
}

func evalInput(t *testing.T, input string) []object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	tpl := p.Parse()

	if len(p.Errors) > 0 {
		for _, msg := range p.Errors {
			t.Errorf("Parser error %q", msg)
		}

		t.FailNow()
	}

	e := New(tpl)
	return e.Run()
}

func checkStatementCount(t *testing.T, results []object.Object, expected int) {
	if len(results) != expected {
		t.Fatalf("Wrong number of results. Expected %d Got %d", expected, len(results))
	}
}

func checkObject(t *testing.T, obj object.Object, objType object.ObjectType, expected interface{}) {
	if obj.Type() != objType {
		t.Fatalf("Wrong object type. Expected %s Got %T", objType, obj)
	}

	if obj.Value() != expected {
		t.Fatalf("Wrong value of object. Expected %v Got %v", expected, obj.Value())
	}
}
