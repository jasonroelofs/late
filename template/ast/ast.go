package ast

import (
	"github.com/jasonroelofs/late/template/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Template struct {
	Statements []Statement
}

func (t *Template) TokenLiteral() string {
	if len(t.Statements) > 0 {
		return t.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type RawStatement struct {
	Token token.Token
}

func (r *RawStatement) statementNode()       {}
func (r *RawStatement) TokenLiteral() string { return r.Token.Literal }
