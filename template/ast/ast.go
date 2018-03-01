package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonroelofs/late/template/token"
)

type Node interface {
	String() string
}

// Statements are code inside of {% %} brackets.
// They do not themselves ever produce output.
type Statement interface {
	Node

	// Nifty type-system hack to make sure that Statements
	// and Expressions don't get mixed up
	statementNode()
}

// Expressions are code inside of {{ }} brackets.
// They always produce a result to include in the output.
type Expression interface {
	Node

	// Nifty type-system hack to make sure that Statements
	// and Expressions don't get mixed up
	expressionNode()
}

// Template is always the root node of the AST
type Template struct {
	Nodes []Node
}

func (t *Template) AddNode(node Node) {
	t.Nodes = append(t.Nodes, node)
}

func (t *Template) String() string {
	builder := strings.Builder{}

	for _, n := range t.Nodes {
		builder.WriteString(n.String())
	}

	return builder.String()
}

/**
 * This token handles all of the non-Liquid raw text
 * that the template can contain.
 * This content is copied verbatim into the final results.
 */
type RawStatement struct {
	Token token.Token
}

func (r *RawStatement) statementNode() {}
func (r *RawStatement) String() string { return r.Token.Literal }

type ExpressionStatement struct {
	// First token of the statement
	Token      token.Token
	Expression Expression
}

func (v *ExpressionStatement) statementNode() {}
func (v *ExpressionStatement) String() string { return "" }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) String() string {
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Token.Literal)
}
