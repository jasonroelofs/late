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
		So much { Not % quite { { code } % } here.
		"This is stringy"
		{{ "This is a string" | 'that is a string' | 100 | -437.6 }}
		{{ 1 < 2 > 3 <= 4 >= 5 * 6 + 7 - 8 / 9 == 0 != 10 }}
		{{ true | false }}
		{{ (1 + 2) }}{{ 3 }}
		{{ [1, 2] }}
		One more raw token`

	// There's a lot going on here
	// Parsing late is itself a stateful system. We want to keep around, untouched,
	// the raw text from the template and only replace the parts that are actually late.
	// Thus as we parse we need to keep track of the raw text itself, and only go into
	// serious parsing mode when we hit {{ or {%.

	tests := []ExpectedToken{
		{token.RAW, "\n\t\tRaw Text "}, // 0
		{token.OPEN_VAR, "{{"},
		{token.IDENT, "variable"},
		{token.DOT, "."},
		{token.IDENT, "method"},
		{token.CLOSE_VAR, "}}"}, // 5
		{token.RAW, "\n\t\t"},
		{token.OPEN_TAG, "{%"},
		{token.IDENT, "tag"},
		{token.CLOSE_TAG, "%}"},
		{token.RAW, "\n\t\t\tStuff here\n\t\t"}, // 10
		{token.END, "{%end%}"},
		{token.RAW, "\n\t\tSo much { Not % quite { { code } % } here.\n\t\t\"This is stringy\"\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.STRING, "This is a string"},
		{token.PIPE, "|"}, // 15
		{token.STRING, "that is a string"},
		{token.PIPE, "|"},
		{token.NUMBER, "100"},
		{token.PIPE, "|"},
		{token.MINUS, "-"}, // 20
		{token.NUMBER, "437.6"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.NUMBER, "1"}, // 25
		{token.LT, "<"},
		{token.NUMBER, "2"},
		{token.GT, ">"},
		{token.NUMBER, "3"},
		{token.LT_EQ, "<="}, // 30
		{token.NUMBER, "4"},
		{token.GT_EQ, ">="},
		{token.NUMBER, "5"},
		{token.TIMES, "*"},
		{token.NUMBER, "6"}, // 35
		{token.PLUS, "+"},
		{token.NUMBER, "7"},
		{token.MINUS, "-"},
		{token.NUMBER, "8"},
		{token.SLASH, "/"}, // 40
		{token.NUMBER, "9"},
		{token.EQ, "=="},
		{token.NUMBER, "0"},
		{token.NOT_EQ, "!="},
		{token.NUMBER, "10"}, // 45
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.TRUE, "true"},
		{token.PIPE, "|"}, // 50
		{token.FALSE, "false"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"},
		{token.LPAREN, "("}, // 55
		{token.NUMBER, "1"},
		{token.PLUS, "+"},
		{token.NUMBER, "2"},
		{token.RPAREN, ")"},
		{token.CLOSE_VAR, "}}"}, // 60
		{token.OPEN_VAR, "{{"},
		{token.NUMBER, "3"},
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_VAR, "{{"}, // 65
		{token.LSQUARE, "["},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RSQUARE, "]"}, // 70
		{token.CLOSE_VAR, "}}"},
		{token.RAW, "\n\t\tOne more raw token"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestOnlyCodeTemplates(t *testing.T) {
	input := "{{ variable }}"

	tests := []ExpectedToken{
		{token.OPEN_VAR, "{{"},
		{token.IDENT, "variable"},
		{token.CLOSE_VAR, "}}"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestCodeStrings(t *testing.T) {
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

func TestNoCode(t *testing.T) {
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

func TestTokenRawWithWhitespace(t *testing.T) {
	input := `
		{{ "Some text here" }} {{ 1 }}
		{% assign %} {% end %}`

	tests := []struct {
		expectedLiteral string
		expectedRaw     string
	}{
		{"\n\t\t", "\n\t\t"},
		{"{{", "{{"},
		{"Some text here", " \"Some text here\""},
		{"}}", " }}"},
		{" ", " "},
		{"{{", "{{"},
		{"1", " 1"},
		{"}}", " }}"},
		{"\n\t\t", "\n\t\t"},
		{"{%", "{%"},
		{"assign", " assign"},
		{"%}", " %}"},
		{" ", " "},
		{"{%end%}", "{% end %}"},
	}

	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Literal != test.expectedLiteral {
			t.Fatalf(
				"(%d) Wrong literal content, expected=%q, got=%q (%s)",
				i, test.expectedLiteral, tok.Literal, tok.Type,
			)
		}

		if tok.Raw != test.expectedRaw {
			t.Fatalf(
				"(%d) Wrong raw content, expected=%q, got=%q (%s)",
				i, test.expectedRaw, tok.Raw, tok.Type,
			)
		}
	}
}

func TestCommentsAndRaw(t *testing.T) {
	input := `
		{{{ This is {{ "Raw Liquid" }} }}}
		{# Ignore me {% {{ #}
		{{{ Invalid {{ !Liquid {% }}}`

	tests := []ExpectedToken{
		{token.RAW, "\n\t\t"},
		{token.OPEN_RAW, "{{{"},
		{token.RAW, " This is {{ \"Raw Liquid\" }} "},
		{token.CLOSE_RAW, "}}}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_COMMENT, "{#"},
		{token.RAW, " Ignore me {% {{ "},
		{token.CLOSE_COMMENT, "#}"},
		{token.RAW, "\n\t\t"},
		{token.OPEN_RAW, "{{{"},
		{token.RAW, " Invalid {{ !Liquid {% "},
		{token.CLOSE_RAW, "}}}"},
		{token.EOF, ""},
	}

	testTemplateGeneratesTokens(t, input, tests)
}

func TestLineAndCharacterNumbers(t *testing.T) {
	input := `First
		{{ "Second" }}
		{% assign
				third = "third" %}{{ third }}`

	tests := []struct {
		literal string
		line    int
		char    int
	}{
		{"First\n\t\t", 1, 1},
		{"{{", 2, 3},
		{`Second`, 2, 6},
		{"}}", 2, 15},
		{"\n\t\t", 2, 17},
		{"{%", 3, 3},
		{"assign", 3, 6},
		{"third", 4, 5},
		{"=", 4, 11},
		{`third`, 4, 13},
		{"%}", 4, 21},
		{"{{", 4, 23},
		{"third", 4, 26},
		{"}}", 4, 32},
	}

	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Literal != test.literal {
			t.Fatalf("(%d) Wrong token returned from lexer. expected=%s got=%s", i, test.literal, tok.Literal)
		}

		if tok.Line != test.line {
			t.Fatalf("(%d) Wrong line number on %#v, expected=%d got=%d", i, test.literal, test.line, tok.Line)
		}

		if tok.Char != test.char {
			t.Fatalf("(%d) Wrong character number on %#v, expected=%d got=%d", i, test.literal, test.char, tok.Char)
		}
	}
}

func testTemplateGeneratesTokens(t *testing.T, template string, expectedTokens []ExpectedToken) {
	l := New(template)

	for i, tt := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("(%d) Wrong token type, expected=%q, got=%q (%s)", i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("(%d) Wrong literal, expected=%q, got=%q (%s)", i, tt.expectedLiteral, tok.Literal, tok.Type)
		}
	}
}
