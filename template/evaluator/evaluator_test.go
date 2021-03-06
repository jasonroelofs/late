package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/parser"
)

func TestRawStatements(t *testing.T) {
	input := "This is a raw, non-liquid template"

	results := evalInput(t, input, context.New())

	checkStatementCount(t, results, 1)
	checkObject(t, results[0], object.TYPE_STRING, input)
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
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.TYPE_NUMBER, test.expected)
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
		{"{{ 1 == 1 }}", true},
		{"{{ 1 != 1 }}", false},
		{`{{ "this" == "this" }}`, true},
		{`{{ "this" != "that" }}`, true},
		{`{{ "this" == "that" }}`, false},
	}

	for _, test := range tests {
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.TYPE_BOOL, test.expected)
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
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.TYPE_STRING, test.expected)
	}
}

type ExpectedElement struct {
	objType object.ObjectType
	value   interface{}
}

func TestArrays(t *testing.T) {
	tests := []struct {
		input          string
		expectedLength int
		expectedValues []ExpectedElement
	}{
		{`{{ [] }}`, 0, []ExpectedElement{}},
		{`{{ [1,2,3] }}`, 3,
			[]ExpectedElement{
				{object.TYPE_NUMBER, float64(1)},
				{object.TYPE_NUMBER, float64(2)},
				{object.TYPE_NUMBER, float64(3)},
			},
		},
		{`{{ ["two", 1 + 3] }}`, 2,
			[]ExpectedElement{
				{object.TYPE_STRING, "two"},
				{object.TYPE_NUMBER, float64(4)},
			},
		},
	}

	for _, test := range tests {
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)

		obj := results[0]
		if obj.Type() != object.TYPE_ARRAY {
			t.Fatalf("Did not get an array. Got %T", obj)
		}

		array := obj.(*object.Array)
		if len(array.Elements) != test.expectedLength {
			t.Fatalf("Wrong number of array elements. Expected %d got %d",
				test.expectedLength, len(array.Elements),
			)
		}

		for i, expValue := range test.expectedValues {
			checkObject(t, array.Elements[i], expValue.objType, expValue.value)
		}
	}
}

func TestArrayAccess(t *testing.T) {
	tests := []struct {
		input        string
		expectedType object.ObjectType
		expected     interface{}
	}{
		{`{{ [][0] }}`, object.TYPE_NULL, nil},
		{`{{ [][-1] }}`, object.TYPE_NULL, nil},
		{`{{ [1, 2, 3][0] }}`, object.TYPE_NUMBER, float64(1)},
		{`{{ [1, 2, 3][1] }}`, object.TYPE_NUMBER, float64(2)},
		{`{{ [1, 2, 3][2] }}`, object.TYPE_NUMBER, float64(3)},
		{`{{ [1, 2, 3][3] }}`, object.TYPE_NULL, nil},
		{`{{ ["one", 2][0] }}`, object.TYPE_STRING, "one"},
		{`{{ [1,2,5][ [1,2][1] ] }}`, object.TYPE_NUMBER, float64(5)},
	}

	for _, test := range tests {
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], test.expectedType, test.expected)
	}
}

func TestHashAccess(t *testing.T) {
	tests := []string{
		`{{ site.root.title }}`,
		`{{ site["root"]["title"] }}`,
		`{{ site[l1][l2] }}`,
	}

	ctx := context.New()
	ctx.Set("site", map[string]interface{}{"root": map[string]interface{}{"title": "Site Title"}})
	ctx.Set("l1", "root")
	ctx.Set("l2", "title")

	for _, test := range tests {
		results := evalInput(t, test, ctx)

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], object.TYPE_STRING, "Site Title")
	}
}

func TestFilters(t *testing.T) {
	tests := []struct {
		input        string
		expectedType object.ObjectType
		expected     interface{}
	}{
		{`{{ "A String" | size }}`, object.TYPE_NUMBER, float64(8)},
		{`{{ "A String" | upcase }}`, object.TYPE_STRING, "A STRING"},
		{`{{ "Hello Mom" | replace: "Mom", with: "World" }}`, object.TYPE_STRING, "Hello World"},
		{`{{ "Hello Mom" | replace: " Mom", with: "" | upcase }}`, object.TYPE_STRING, "HELLO"},
		{`{{ "Hello Mom" | replace: "Mom", with: ("World" | upcase) }}`, object.TYPE_STRING, "Hello WORLD"},
		// TODO: Unknown filter
		//   Strict: error out
		//   Lax: treat as a pass-through no-op, trigger a warning
	}

	for _, test := range tests {
		results := evalInput(t, test.input, context.New())

		checkStatementCount(t, results, 1)
		checkObject(t, results[0], test.expectedType, test.expected)
	}
}

func TestVariables(t *testing.T) {
	tests := []struct {
		input        string
		assigns      context.Assigns
		expectedType object.ObjectType
		expected     interface{}
	}{
		{"{{ page }}", context.Assigns{"page": "home"}, object.TYPE_STRING, "home"},
		{"{{ count }}", context.Assigns{"count": 10}, object.TYPE_NUMBER, float64(10)},
		{"{{ unknown }}", context.Assigns{}, object.TYPE_NULL, nil},

		// Test variable usage as filter parameters
		{
			"{{ page | replace: page, with: changeTo | upcase }}",
			context.Assigns{"page": "home", "changeTo": "blog"},
			object.TYPE_STRING,
			"BLOG",
		},

		// TODO: if variables end up nil in the case above, what do we do?
	}

	for _, test := range tests {
		ctx := context.New()
		ctx.Assign(test.assigns)

		results := evalInput(t, test.input, ctx)
		checkStatementCount(t, results, 1)
		checkObject(t, results[0], test.expectedType, test.expected)
	}
}

func TestTags(t *testing.T) {
	tests := []struct {
		input        string
		expectedType object.ObjectType
		expected     interface{}
	}{
		{`{% assign page = "home" %}{{ page }}`, object.TYPE_STRING, "home"},
		{"{% assign count = 10 %}{{ count }}", object.TYPE_NUMBER, float64(10)},
		{`{% assign page_size = "home" | size %}{{ page_size }}`, object.TYPE_NUMBER, float64(4)},

		{`{% if true %}True{% end %}`, object.TYPE_STRING, "True"},
		{`{% if false %}True{% else %}False{% end %}`, object.TYPE_STRING, "False"},
		{`{% if 1 == 2 %}True{% elsif 2 == 3 %}False{% end %}`, object.TYPE_NULL, nil},
	}

	for _, test := range tests {
		results := evalInput(t, test.input, context.New())

		checkObject(t, results[len(results)-1], test.expectedType, test.expected)
	}
}

func TestNestingBlockTags(t *testing.T) {
	input := `
		{% if true %}
			{% capture nested %}
				{% if false %}
					I am false!
				{% else %}
					I am true!
				{% end %}
			{% end %}
		{% end %}
		{{ nested }}`

	results := evalInput(t, input, context.New())

	for i, r := range results {
		fmt.Printf("(%d) %s", i, r.Inspect())
	}

	output := results[len(results)-1].Value().(string)
	trimmed := strings.TrimSpace(output)

	if trimmed != "I am true!" {
		t.Fatalf("Wrong value of object. Expected %v Got %v", "I am true!", trimmed)
	}
}

func evalInput(t *testing.T, input string, ctx *context.Context) []object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	tpl := p.Parse()

	if len(p.Errors) > 0 {
		for _, msg := range p.Errors {
			t.Errorf("Parser error %q", msg)
		}

		t.FailNow()
	}

	e := New(tpl, ctx)
	return e.Run()
}

func checkStatementCount(t *testing.T, results []object.Object, expected int) {
	if len(results) != expected {
		t.Fatalf("Wrong number of results. Expected %d Got %d", expected, len(results))
	}
}

func checkObject(t *testing.T, obj object.Object, objType object.ObjectType, expected interface{}) {
	if obj.Type() != objType {
		t.Fatalf("Wrong object type. Expected %s Got %T. Test: \"%v\"", objType, obj, expected)
	}

	if obj.Value() != expected {
		t.Fatalf(
			"Wrong value of object. Expected %v (%T) Got %v (%T)",
			expected, expected, obj.Value(), obj.Value(),
		)
	}
}
