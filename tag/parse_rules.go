package tag

import (
	"github.com/jasonroelofs/late/template/token"
)

/**
 * Parsing Rules
 * The following constructs are how tags define to the template how to parse and evaluate
 * code for the given tag.
 */

type ParseRule interface {
}

type IdentifierRule struct {
}

type LiteralRule struct {
	Value string
}

type TokenRule struct {
	Type token.TokenType
}

type ExpressionRule struct {
}

func Identifier() ParseRule             { return &IdentifierRule{} }
func Token(t token.TokenType) ParseRule { return &TokenRule{Type: t} }
func Literal(value string) ParseRule    { return &LiteralRule{Value: value} }
func Expression() ParseRule             { return &ExpressionRule{} }
