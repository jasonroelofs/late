package lexer

import (
	"testing"

	"github.com/jasonroelofs/late/template/token"
)

type ExpectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestTokenizesInput(t *testing.T) {
	input := `
		Raw Text {{ variable.method }}
		{% tag %}
			Stuff here
		{% end %}
		So much { Not % quite { { liquid } % } here.
		"This is stringy"
		{{ "This is a string" | 'that is a string' | 100 | -437.6 }}
		{{ 1 < 2 > 3 <= 4 >= 5 * 6 + 7 - 8 / 9 == 0 != 10 }}
		{{ true | false }}
		One more raw token`

	// There's a lot going on here
	// Parsing liquid is itself a stateful system. We want to keep around, untouched,
	// the raw text from the template and only replace the parts that are actually liquid.
	// Thus as we parse we need to keep track of the raw text itself, and only go into
	// serious parsing mode when we hit {{ or {%.

	tests := []ExpectedToken{
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
		{token.PIPE, "|"},
		{token.NUMBER, "100"},
		{token.PIPE, "|"},
		{token.MINUS, "-"},
		{token.NUMBER, "437.6"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.NUMBER, "1"},
		{token.LT, "<"},
		{token.NUMBER, "2"},
		{token.GT, ">"},
		{token.NUMBER, "3"},
		{token.LT_EQ, "<="},
		{token.NUMBER, "4"},
		{token.GT_EQ, ">="},
		{token.NUMBER, "5"},
		{token.TIMES, "*"},
		{token.NUMBER, "6"},
		{token.PLUS, "+"},
		{token.NUMBER, "7"},
		{token.MINUS, "-"},
		{token.NUMBER, "8"},
		{token.SLASH, "/"},
		{token.NUMBER, "9"},
		{token.EQ, "=="},
		{token.NUMBER, "0"},
		{token.NOT_EQ, "!="},
		{token.NUMBER, "10"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.TRUE, "true"},
		{token.PIPE, "|"},
		{token.FALSE, "false"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\tOne more raw token"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestOnlyLiquidTemplates(t *testing.T) {
	input := "{{ variable }}"

	tests := []ExpectedToken{
		{token.OPEN_VAR, "{{"},
		{token.IDENT, "variable"},
		{token.CLOSE_VAR, "}}"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestLiquidStrings(t *testing.T) {
	input := `
		{{ "This is ' a string" }}
		{{ 'Single \' quotes' }}
		{{ "Double \" quotes" }}
		{{ "New \n Lines \n here" }}
		{{ "Other \t command chars" }}
	`

	tests := []ExpectedToken{
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "This is ' a string"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "Single \\' quotes"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "Double \\\" quotes"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "New \\n Lines \\n here"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "Other \\t command chars"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestRawAtEOF(t *testing.T) {
	input := `Before {{ variable }} After `

	tests := []ExpectedToken{
		{token.RAW, "Before "},
		{token.OPEN_VAR, "{{"},
		{token.IDENT, "variable"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, " After "},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestNoLiquid(t *testing.T) {
	input := `Before and After`

	tests := []ExpectedToken{
		{token.RAW, "Before and After"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestEmptyTemplate(t *testing.T) {
	input := ``

	tests := []ExpectedToken{
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func testTemplateGeneratesTokens(t *testing.T, template string, expectedTokens []ExpectedToken) {
	l := New(template)

	for i, tt := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("(%d) Wrong token type, expected=%q, got=%q (%s)", i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("(%d) Wrong literal, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
