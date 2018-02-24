package lexer

import (
	"testing"

	"github.com/jasonroelofs/late/parser/token"
)

func TestTokenizesInput(t *testing.T) {
	input := `
		Raw Text {{ variable.method }}
		{% tag %}
			Stuff here
		{% end %}
		So much { Not % quite { { liquid } % } here.
		"This is stringy"
		{{ "This is a string" | 'that is a string' }}
	`

	// There's a lot going on here
	// Parsing liquid is itself a stateful system. We want to keep around, untouched,
	// the raw text from the template and only replace the parts that are actually liquid.
	// Thus as we parse we need to keep track of the raw text itself, and only go into
	// serious parsing mode when we hit {{ or {%.

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.RAW, "\n\t\tRaw Text "},
		{token.OPEN_VAR, "{{"},
		{token.IDENT, "variable"},
		{token.DOT, "."},
		{token.IDENT, "method"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_TAG, "{%"},
		{token.IDENT, "tag"},
		{token.CLOSE_TAG, "%}"},
		{token.RAW, "\n\t\t\tStuff here\n\t\t"},
		{token.OPEN_TAG, "{%"},
		{token.IDENT, "end"},
		{token.CLOSE_TAG, "%}"},
		{token.RAW, "\n\t\tSo much { Not % quite { { liquid } % } here.\n\t\t\"This is stringy\"\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "This is a string"},
		{token.PIPE, "|"},
		{token.STRING, "that is a string"},
		{token.CLOSE_VAR, "}}"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("(%d) Wrong token type, expected=%q, got=%q (%s)", i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("(%d) Wrong literal, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// So many error cases to watch out for.
// End-of-file with anything.
// Un-terminated strings, both quote types
