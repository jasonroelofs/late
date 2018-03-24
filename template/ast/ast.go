package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonroelofs/late/tag"
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
	return v.Expression.String()
}

type TagStatement struct {
	Token          token.Token
	TagName        string
	Tag            tag.Tag
	Nodes          []Expression
	BlockStatement *BlockStatement
}

func (t *TagStatement) statementNode() {}
func (t *TagStatement) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(t.TagName)

	for _, expr := range t.Nodes {
		out.WriteString(expr.String())
	}

	out.WriteString(")")

	return out.String()
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) String() string {
	out := strings.Builder{}

	for _, stmt := range b.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
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

// FilterExpression is the parent expression holder node that
// keeps track of filter requests. Specifically, it references
// the "input" side and the "filter" side of a call like:
//
//   {{ input | filter }}
//
// This is almost identical to an InfixExpression but it's handy
// to have an explicit type for this when it comes to evaluation.
type FilterExpression struct {
	Token  token.Token
	Input  Expression
	Filter Expression
}

func (f *FilterExpression) expressionNode() {}
func (f *FilterExpression) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(f.Input.String())
	out.WriteString(" | ")
	out.WriteString(f.Filter.String())
	out.WriteString(")")

	return out.String()
}

/**
 * Literals
 * These AST nodes evaluate to themselves
 */

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

type FilterLiteral struct {
	Token      token.Token
	Name       string
	Parameters map[string]Expression
}

func (f *FilterLiteral) expressionNode() {}
func (f *FilterLiteral) String() string {
	out := strings.Builder{}
	var params []string
	groupExpr := len(f.Parameters) > 0

	if groupExpr {
		out.WriteString("(")
	}

	out.WriteString(f.Name)

	if len(f.Parameters) > 0 {
		params = append(params, fmt.Sprintf(": %s", f.Parameters[f.Name].String()))

		for name, expr := range f.Parameters {
			if name == f.Name {
				continue
			}

			params = append(params, fmt.Sprintf("%s: %s", name, expr.String()))
		}

		out.WriteString(strings.Join(params, ", "))
	}

	if groupExpr {
		out.WriteString(")")
	}

	return out.String()
}
