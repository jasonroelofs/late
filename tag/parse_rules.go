package tag

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

type ExpressionRule struct {
}

func Identifier() ParseRule          { return &IdentifierRule{} }
func Literal(value string) ParseRule { return &LiteralRule{Value: value} }
func Expression() ParseRule          { return &ExpressionRule{} }
