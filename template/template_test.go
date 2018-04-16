package template

import (
	"strings"
	"testing"

	"github.com/jasonroelofs/late/context"
)

func TestNew(t *testing.T) {
	tpl := New("This is a template")

	if tpl.body != "This is a template" {
		t.Errorf("Did not store the template body")
	}
}

func TestRender(t *testing.T) {
	tpl := New("This is a template")
	results := tpl.Render(context.New())
	checkNoErrors(t, tpl)

	if results != "This is a template" {
		t.Errorf("Failed to render the template")
	}
}

func TestRenderLiquidWithLiterals(t *testing.T) {
	tests := []struct {
		liquidInput    string
		expectedOutput string
	}{
		{"{{ 3 }}", "3"},
		{"{{ 1 + 2 }}", "3"},
		{"{{ 1 / 2 }}", "0.5"},
		//		{"{{ \"Hi\" }}", "Hi"},
		//		{"{{ 'Hi' + ' ' + 'Bye' }}", "Hi Bye"},
	}

	for _, test := range tests {
		tpl := New(test.liquidInput)
		results := tpl.Render(context.New())
		checkNoErrors(t, tpl)

		if results != test.expectedOutput {
			t.Errorf("Failed to render the template. Expected '%s' got '%s'", test.expectedOutput, results)
		}
	}
}

func TestRenderWithInitialState(t *testing.T) {
	tests := []struct {
		input    string
		assigns  context.Assigns
		expected string
	}{
		{"{{ page }}", context.Assigns{"page": "home"}, "home"},
	}

	for _, test := range tests {
		tpl := New(test.input)
		ctx := context.New()
		ctx.Assign(test.assigns)

		results := tpl.Render(ctx)
		checkNoErrors(t, tpl)

		if results != test.expected {
			t.Errorf("Failed to render. Expected '%s' got '%s'", test.expected, results)
		}
	}
}

func TestRenderWithComplexObject(t *testing.T) {
	input := "{{ site.root.title }}"

	tpl := New(input)
	ctx := context.New()
	ctx.Set("site", map[string]interface{}{"root": map[string]interface{}{"title": "Site Title"}})

	results := tpl.Render(ctx)
	checkNoErrors(t, tpl)

	if results != "Site Title" {
		t.Errorf("Failed to render. Expected got '%s'", results)
	}
}

func TestRender_RawAndComments(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Before {# Middle #} End`, "Before  End"},
		{`{# {{ "hi" }} #}`, ""},

		{`This is {{{ {{ "Raw Late Code" }} }}}`, `This is  {{ "Raw Late Code" }} `},

		// Raw and Comment should not try to parse code inside of their blocks.
		// I would expect invalid liquid to be ignored, not erroring out.
		{`This is {{{ Invalid {% {{ >=}}}`, `This is  Invalid {% {{ >=`},
		{`{# Don't {{ break {% ==#}`, ""},
	}

	for i, test := range tests {
		tpl := New(test.input)
		results := tpl.Render(context.New())
		checkNoErrors(t, tpl)

		if results != test.expected {
			t.Errorf("(%d) Failed to render. Expected '%s' got '%s'", i, test.expected, results)
		}
	}
}

func TestRender_Tags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 0
		{`{% assign page = "home" %}{{ page }}`, "home"},
		{`{% assign page = "home" | upcase %}{{ page }}`, "HOME"},
		{`{% assign page = "Here" | replace: "Here", with: "There" %}{{ page }}`, "There"},

		// 3
		{`{% capture math %}1 + 2 == {{ 1 + 2 }}{% end %}{{ math }}`, "1 + 2 == 3"},
		{`{% capture outer %}
				{% capture inner %}
					{{ 1 + 2 }}
				{% end %}
				1 + 2 == {{ inner }}
			{% end %}
			{{ outer }}`,
			"1 + 2 == 3",
		},

		// 5
		{`{% if true %}True{% end %}`, "True"},
		{`{% if false %}True{% end %}`, ""},
		{`{% if false %}True{% else %}False{% end %}`, "False"},
		{`{% assign num = 7 %}
			{% if num > 10 %}
				Big
			{% elsif num > 7 %}
				Big-ish
			{% elsif num > 5 %}
				Medium
			{% else %}
				Small
			{% end %}`,
			"Medium",
		},

		// 9
		{`{% for num in [1,2,3] %}{{ num }}{% end %}`, "123"},
		{`{% for num in [1,2,3] %}
				{% if num == 1 %}{% continue %}{% end %}
				{{ num }}
			{% end %}`,
			"23",
		},
		{`{% for x in [1,2,3] %}
				{% for y in [1,2,3] %}
					{% if y == 3 %}{% break %}{% end %}
					{{ x }},{{ y }}-
				{% end %}
			{% end %}`,
			"1,1-1,2-2,1-2,2-3,1-3,2-",
		},
		{`{% for num in [1,2,3] %}
				{{ num }}
				{% if num == 2 %}{% break %}{% end %}
			{% end %}`,
			"12",
		},
		{`{% for num in [1,2,3] %}
				{{ num }}
				{% if num == 2 %}Break{% break %}It Up{% end %}
			{% end %}`,
			"12Break",
		},
	}

	// TODO: Build a set of rules around whitespace management.
	// For now, clear out all newlines and tab characters for test comparisons.
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "")

	for i, test := range tests {
		tpl := New(test.input)
		results := tpl.Render(context.New())
		checkNoErrors(t, tpl)

		trimmed := replacer.Replace(results)

		if trimmed != test.expected {
			t.Errorf("(%d) Failed to render. Expected '%s' got '%s'", i, test.expected, trimmed)
		}
	}
}

func checkNoErrors(t *testing.T, tpl *Template) {
	if len(tpl.Errors) != 0 {
		t.Fatalf("Errors rendering the template:\n%s", strings.Join(tpl.Errors, "\n"))
	}
}
