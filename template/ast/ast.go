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

// Template is always the root node of the AST.
type Template struct {
	Statements []Statement
}

func (t *Template) AddStatement(stmt Statement) {
	t.Statements = append(t.Statements, stmt)
}

func (t *Template) String() string {
	builder := strings.Builder{}

	for _, stmt := range t.Statements {
		builder.WriteString(stmt.String())
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

type VariableStatement struct {
	Token      token.Token
	Expression Expression
}

func (v *VariableStatement) statementNode() {}
func (v *VariableStatement) String() string {
	// TODO This case should never be possible. Will probably solve itself
	// as the rest of the parser fleshes itself out.
	// This is causing panics in `make docs` as it's trying to parse and print out
	// parts of liquid we don't parse into the tree yet.
	if v != nil {
		return v.Expression.String()
	} else {
		return ""
	}
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

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

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string  { return b.Token.Literal }
