package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/token"
)

const (
	_ int = iota
	LOWEST
	PIPE    // '|' filter seperator
	ASSIGN  // =
	COMPARE // ==, < and >
	SUM     // +, -
	PRODUCT // *, /
	PREFIX  // -X, !X
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	Errors []string

	currToken token.Token
	peekToken token.Token

	// Pratt Parsing!
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{l: lexer}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)

	// Read the first two tokens to pre-fill
	// curr and peek token values
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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
			template.AddStatement(stmt)
		}

		p.nextToken()
	}

	return template
}

func (p *Parser) parseNext() ast.Statement {
	switch p.currToken.Type {
	case token.OPEN_VAR:
		return p.parseVariableStatement()
	default:
		return p.parseRawStatement()
	}
}

func (p *Parser) parseRawStatement() *ast.RawStatement {
	return &ast.RawStatement{Token: p.currToken}
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT, token.CLOSE_VAR, token.NUMBER, token.STRING) {
		return nil
	}

	// Skip past the opening {{
	// and parse the content into an expression tree
	p.nextToken()

	stmt.Expression = p.parseExpression(LOWEST)

	// Advance to make sure we are now at a closing }}
	p.nextToken()

	if p.currToken.Type != token.CLOSE_VAR {
		p.expectedTokenError(p.currToken.Type, token.CLOSE_VAR)
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	literal := &ast.NumberLiteral{Token: p.currToken}

	number, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		p.parserErrorf("could not parse %q as a number", p.currToken.Literal)
		return nil
	}

	literal.Value = number
	return literal
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
		p.expectedTokenError(currPeek, allowed...)

		return false
	}

	return true
}

func (p *Parser) expectedTokenError(got token.TokenType, expected ...token.TokenType) {
	var tokenNames []string
	for _, t := range expected {
		tokenNames = append(tokenNames, string(t))
	}

	msg := fmt.Sprintf("expected %s, found %s", strings.Join(tokenNames, " or "), got)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) parserErrorf(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args)
	p.Errors = append(p.Errors, msg)
}
