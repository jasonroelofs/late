package evaluator

import (
	"testing"

	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/object"
	"github.com/jasonroelofs/late/template/parser"
)

func TestRawStatements(t *testing.T) {
	input := "This is a raw, non-liquid template"

	results := evalInput(input)

	if len(results) != 1 {
		t.Fatalf("Got the wrong number of results, got %d", len(results))
	}

	str, ok := results[0].(*object.String)
	if !ok {
		t.Fatalf("Expected a String, got %T", results)
	}

	if str.Value != input {
		t.Fatalf("The eval came through wrong, got '%s'", str.Value)
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
		results := evalInput(test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		number, ok := results[0].(*object.Number)
		if !ok {
			t.Fatalf("(%d) Expected a Number, got %T", i, results)
		}

		if number.Value != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%f'", i, number.Value)
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
		results := evalInput(test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		boolean, ok := results[0].(*object.Boolean)
		if !ok {
			t.Fatalf("(%d) Expected a Boolean, got %T", i, results)
		}

		if boolean.Value != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%v'", i, boolean.Value)
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
		results := evalInput(test.input)

		if len(results) != 1 {
			t.Fatalf("(%d) Got the wrong number of results, got %d", i, len(results))
		}

		str, ok := results[0].(*object.String)
		if !ok {
			t.Fatalf("(%d) Expected a String, got %T", i, results)
		}

		if str.Value != test.expected {
			t.Fatalf("(%d) The eval came through wrong, got '%v'", i, str.Value)
		}
	}
}

func evalInput(input string) []object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	t := p.Parse()
	e := New(t)
	return e.Run()
}
