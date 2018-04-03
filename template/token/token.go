package token

type TokenType string

type Token struct {
	Type TokenType

	// Literal is the post-lex'd content of the current token.
	// This will not include any stray punctuation or whitespace
	Literal string

	// Raw is the plain content that was processed to create this
	// token. It includes all punctuation and leading white space
	// and can be used to perfectly reproduce the input.
	Raw string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	RAW    = "RAW"
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	OPEN_VAR  = "OPEN_VAR"
	CLOSE_VAR = "CLOSE_VAR"
	OPEN_TAG  = "OPEN_TAG"
	CLOSE_TAG = "CLOSE_TAG"

	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"
	LSQUARE  = "LSQUARE"
	RSQUARE  = "RSQUARE"
	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"

	TRUE  = "TRUE"
	FALSE = "FALSE"

	DOT     = "DOT"
	COMMA   = "COMMA"
	COLON   = "COLON"
	ASSIGN  = "ASSIGN"
	PERCENT = "PERCENT"
	PIPE    = "PIPE"

	PLUS   = "PLUS"
	MINUS  = "MINUS"
	TIMES  = "TIMES"
	SLASH  = "SLASH"
	LT     = "LT"
	LT_EQ  = "LT_EQ"
	GT     = "GT"
	GT_EQ  = "GT_EQ"
	EQ     = "EQ"
	NOT_EQ = "NOT_EQ"

	// A special lexer token for the {% end %} token of a block.
	END = "END"

	// The following are meta-tokens used when
	// validating parse and evaluation rules for tags.
	EXPRESSION = "EXPRESSION"
)
