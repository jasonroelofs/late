package template

import (
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

		if results != test.expected {
			t.Errorf("Failed to render. Expected '%s' got '%s'", test.expected, results)
		}
	}
}

func TestRender_Tags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{% assign page = "home" %}{{ page }}`, "home"},
		{`{% assign page = "home" | upcase %}{{ page }}`, "HOME"},
		{`{% assign page = "Here" | replace: "Here", with: "There" %}{{ page }}`, "There"},

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

		if results != test.expected {
			t.Errorf("(%d) Failed to render. Expected '%s' got '%s'", i, test.expected, results)
		}
	}
}
