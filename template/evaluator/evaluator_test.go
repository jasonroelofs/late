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

	if len(results) != 1 {
		t.Fatalf("Got the wrong number of results, got %d", len(results))
	}

	str, ok := results[0].(*object.String)
	if !ok {
		t.Fatalf("Expected a String, got %T", results)
	}

	if str.Value().(string) != input {
		t.Fatalf("The eval came through wrong, got '%s'", str.Value())
	}
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

	for i, test := range tests {
		results := evalInput(t, test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		number, ok := results[0].(*object.Number)
		if !ok {
			t.Fatalf("(%d) Expected a Number, got %T", i, results)
		}

		if number.Value().(float64) != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%f'", i, number.Value())
		}
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

	for i, test := range tests {
		results := evalInput(t, test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		boolean, ok := results[0].(*object.Boolean)
		if !ok {
			t.Fatalf("(%d) Expected a Boolean, got %T", i, results)
		}

		if boolean.Value().(bool) != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%v'", i, boolean.Value())
		}
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

	for i, test := range tests {
		results := evalInput(t, test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		str, ok := results[0].(*object.String)
		if !ok {
			t.Fatalf("(%d) Expected a String, got %T", i, results)
		}

		if str.Value().(string) != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%v'", i, str.Value())
		}
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

	for i, test := range tests {
		results := evalInput(t, test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		if results[0].Type() != test.expectedType {
			t.Fatalf("(%d) Expected a %v, got %T", i, test.expectedType, results[0])
		}

		switch test.expectedType {
		case object.OBJ_NUMBER:
			val := results[0].Value().(float64)
			if val != test.expected {
				t.Fatalf("(%d) The eval came through wrong, expected %f got '%f'", i, test.expected, val)
			}
		case object.OBJ_STRING:
			val := results[0].Value().(string)
			if val != test.expected {
				t.Fatalf("(%d) The eval came through wrong, expected %v got '%v'", i, test.expected, val)
			}
		default:
			t.Fatalf("Don't know how to handle the type %v", test.expectedType)
		}

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
