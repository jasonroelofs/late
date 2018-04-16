package tag

import (
	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/template/token"
)

/**
 * ParserConfig is used to apply new rules and options to the parser and evaluator
 * to properly handle the current tag.
 */
type ParseConfig struct {
	// TagName must contain the string name of the tag.
	TagName string

	// Block must be set to true if this tag expects to take a block of code
	// that ends with an {% end %}
	Block bool

	// Rules contains a list of rules informing the parser how to parse
	// the rest of the immediate tag. The rules will be mapped into object.Object records
	// and passed into Eval() during the evaluation phase.
	Rules []ParseRule

	// Block tags can also have sub-tags to provide for conditional execution (e.g. if/elsif/else).
	// Provide the definitions of each subtag type here. The results of these subtags, when found,
	// will be provided in the SubTagResults value of ParseResult.
	SubTags []ParseConfig

	// Flag this tag as one that will fire an Interrupt. When this tag is hit in the evaluator,
	// all further evaluation will halt until another tag handles and clears the interrupt.
	// For examples, see Continue and Break handling in the For tag.
	Interrupt bool
}

// ParseResult is passed into the tags Eval() method during evaulation phase.
type ParseResult struct {
	// The name of the tag we're currently processing
	TagName string

	// Nodes is the list of object.Object records that map directly to whatever
	// was provided in the ParseConfig.Rules field.
	Nodes []object.Object

	// For block-type tags, this list of statements correspond to the content of the
	// block and should be evaulated in order according to the rules of the tag.
	Statements []Statement

	// For block tags that can have sub-tags (e.g. if, case), each found and valid
	// subtag will have it's own set of ParseResults.
	SubTagResults []*ParseResult
}

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
