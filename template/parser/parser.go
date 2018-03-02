package parser

import (
	"fmt"
	"strings"

	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	Errors []string
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{l: lexer}

	// Read the first two tokens to pre-fill
	// curr and peek token values
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Template {
	template := &ast.Template{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseNext()
		if stmt != nil {
			template.AddNode(stmt)
		}

		p.nextToken()
	}

	return template
}

func (p *Parser) parseNext() ast.Node {
	switch p.currToken.Type {
	case token.RAW:
		return p.parseRawStatement()
	case token.OPEN_VAR:
		return p.parseVariableExpression()
	default:
		return nil
	}
}

func (p *Parser) parseRawStatement() *ast.RawStatement {
	return &ast.RawStatement{Token: p.currToken}
}

func (p *Parser) parseVariableExpression() *ast.VariableExpression {
	node := &ast.VariableExpression{Token: p.currToken}

	if !p.expectPeek(token.IDENT, token.CLOSE_VAR) {
		return nil
	}

	// Skip past the opening {{
	p.nextToken()

	for p.currToken.Type != token.CLOSE_VAR && p.currToken.Type != token.EOF {
		p.nextToken()
	}

	if p.currToken.Type != token.CLOSE_VAR {
		p.parserError(p.currToken.Type, token.CLOSE_VAR)
	}

	return node
}

func (p *Parser) expectPeek(allowed ...token.TokenType) bool {
	matched := false
	currPeek := p.peekToken.Type

	for _, allowedType := range allowed {
		if currPeek == allowedType {
			matched = true
		}
	}

	if !matched {
		p.parserError(currPeek, allowed...)

		return false
	}

	return true
}

func (p *Parser) parserError(got token.TokenType, expected ...token.TokenType) {
	var tokenNames []string
	for _, t := range expected {
		tokenNames = append(tokenNames, string(t))
	}

	msg := fmt.Sprintf("expected %s, found %s", strings.Join(tokenNames, " or "), got)
	p.Errors = append(p.Errors, msg)
}
