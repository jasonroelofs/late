package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
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

	DOT     = "DOT"
	COMMA   = "COMMA"
	COLON   = "COLON"
	ASSIGN  = "ASSIGN"
	PERCENT = "PERCENT"
	PIPE    = "PIPE"
)
