package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonroelofs/late"
	"github.com/jasonroelofs/late/tag"
	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/token"
)

// Precedence levels, from lowest to highest
const (
	_ int = iota
	LOWEST
	ASSIGN  // =
	PIPE    // '|' (filter seperator)
	EQUALS  // ==, !=
	COMPARE // <, >, <=, >=
	SUM     // +, -
	PRODUCT // *, /
	PREFIX  // -X
)

var precedences = map[token.TokenType]int{
	token.ASSIGN: ASSIGN,
	token.PIPE:   PIPE,
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     COMPARE,
	token.GT:     COMPARE,
	token.LT_EQ:  COMPARE,
	token.GT_EQ:  COMPARE,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.TIMES:  PRODUCT,
}

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
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.TIMES, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseFilterExpression)

	// Read the first two tokens to pre-fill
	// curr and peek token values
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) CurrentTokenName() string {
	return p.currToken.Literal
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

func (p *Parser) peekTokenIs(tokenTypes ...token.TokenType) bool {
	for _, tt := range tokenTypes {
		if p.peekToken.Type == tt {
			return true
		}
	}

	return false
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
	case token.OPEN_TAG:
		return p.parseTagStatement()
	default:
		return p.parseRawStatement()
	}
}

func (p *Parser) parseRawStatement() *ast.RawStatement {
	return &ast.RawStatement{Token: p.currToken}
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.currToken}

	// Skip past the opening {{
	// and parse the content into an expression tree
	p.nextToken()

	stmt.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(token.CLOSE_VAR) {
		return nil
	}

	// Advance to make sure we are now at a closing }}
	// and can continue
	p.nextToken()

	return stmt
}

func (p *Parser) parseTagStatement() *ast.TagStatement {
	stmt := &ast.TagStatement{Token: p.currToken}

	// Move past the opening {%
	p.nextToken()

	// Store the first token as our name
	stmt.TagName = p.currToken.Literal
	stmt.Tag = late.FindTag(stmt.TagName)

	if stmt.Tag == nil {
		p.parserErrorf("Unknown tag '%s'", stmt.TagName)
		return nil
	}

	expectedTokens := stmt.Tag.Parse()
	tokenIdx := 0
	for !p.peekTokenIs(token.CLOSE_TAG) && !p.peekTokenIs(token.EOF) && len(expectedTokens) > tokenIdx {
		p.nextToken()

		switch parseRule := expectedTokens[tokenIdx].(type) {
		case *tag.IdentifierRule:
			stmt.Nodes = append(stmt.Nodes, &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})
		case *tag.LiteralRule:
			stmt.Nodes = append(stmt.Nodes, &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal})
		case *tag.ExpressionRule:
			stmt.Nodes = append(stmt.Nodes, p.parseExpression(LOWEST))
		default:
			p.parserErrorf("Error parsing tag %s, don't know how to handle ParseRule of type %T", parseRule)
			return nil
		}

		tokenIdx += 1
	}

	if !p.expectPeek(token.CLOSE_TAG) {
		return nil
	}

	// Move to our %} token so we can continue
	p.nextToken()

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.CLOSE_VAR, token.CLOSE_TAG) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseFilterExpression(input ast.Expression) ast.Expression {
	expression := &ast.FilterExpression{
		Token: p.currToken,
		Input: input,
	}

	p.nextToken()
	expression.Filter = p.parseFilter()

	return expression
}

func (p *Parser) parseFilter() ast.Expression {
	expression := &ast.FilterLiteral{
		Token: p.currToken,
		Name:  p.currToken.Literal,
	}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		expression.Parameters = p.parseFilterParameters(expression.Name)
	}

	return expression
}

func (p *Parser) parseFilterParameters(initialParam string) map[string]ast.Expression {
	list := make(map[string]ast.Expression)

	// We need to make sure the parser doesn't accidentally chain parameter
	// expressions with further pipes, so we set PIPE as the lowest precendence here.
	// Then, when we hit something like `(replace: "this", with: "that") | upcase` the parser
	// stops at the `|` instead of seeing `replace: "this", with: ("that" | upcase)`.
	list[initialParam] = p.parseExpression(PIPE)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		// Error check, we should be on an identifier
		paramName := p.currToken.Literal

		if !p.peekTokenIs(token.COLON) {
			// Error something here, must have a colon after the name
		}

		p.nextToken()
		p.nextToken()
		list[paramName] = p.parseExpression(PIPE)
	}

	return list
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	p.nextToken()

	return exp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currToken, Value: p.currToken.Type == token.TRUE}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
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

func (p *Parser) noPrefixParseFnError(token token.TokenType) {
	p.parserErrorf("No known prefix parse function for token type %s", token)
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
