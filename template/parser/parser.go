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
	INDEX   // []
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:  ASSIGN,
	token.PIPE:    PIPE,
	token.EQ:      EQUALS,
	token.NOT_EQ:  EQUALS,
	token.LT:      COMPARE,
	token.GT:      COMPARE,
	token.LT_EQ:   COMPARE,
	token.GT_EQ:   COMPARE,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.SLASH:   PRODUCT,
	token.TIMES:   PRODUCT,
	token.LSQUARE: INDEX,
	token.DOT:     INDEX,
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

	// Keep track of a stack of Tags we're parsing, as they
	// can be nested quite deeply and will often have SubTags that we need
	// to associate with the correct parent Tag.
	currentTagStack []*ast.TagStatement
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
	p.registerPrefix(token.LSQUARE, p.parseArrayLiteral)
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
	p.registerInfix(token.LSQUARE, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseDotExpression)

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

func (p *Parser) currTokenIs(tokenType token.TokenType) bool {
	return p.currToken.Type == tokenType
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

	for !p.currTokenIs(token.EOF) {
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
	case token.OPEN_COMMENT:
		return p.parseCommentStatement()
	case token.OPEN_RAW:
		return p.parseVerbatimStatement()
	default:
		return p.parseRawStatement()
	}
}

func (p *Parser) parseRawStatement() *ast.RawStatement {
	return &ast.RawStatement{Token: p.currToken}
}

func (p *Parser) parseCommentStatement() *ast.RawStatement {
	p.nextToken()

	if !p.expectPeek(token.CLOSE_COMMENT) {
		return nil
	}

	p.nextToken()
	return &ast.RawStatement{}
}

func (p *Parser) parseVerbatimStatement() *ast.RawStatement {
	// Skip the starting {{{
	p.nextToken()

	stmt := p.parseRawStatement()

	if !p.expectPeek(token.CLOSE_RAW) {
		return nil
	}

	// Move to the ending }}}
	p.nextToken()

	return stmt
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

	// Move past the opening {%
	p.nextToken()

	var currentParseConfig *tag.ParseConfig
	var inSubTag bool
	currTagName := p.currToken.Literal

	stmt := p.currentTag()

	// If there's no statement at all, start one
	// If there's a statement, check that we're on a sub tag
	// - If no matching subtag, try to start a new tag and push on the stack
	// - If matching subtag, continue and store

	if stmt == nil {
		// Store the first token as our name
		// and find the tag object with that name.
		tagName := p.currToken.Literal
		stmt = &ast.TagStatement{
			Token:   p.currToken,
			TagName: tagName,
			Tag:     late.FindTag(tagName),
		}

		if stmt.Tag == nil {
			p.parserErrorf("Unknown tag '%s'", stmt.TagName)
			return nil
		}

		p.pushCurrentTag(stmt)
		currentParseConfig = stmt.Tag.Parse()

	} else if stmt.HasSubTag(currTagName) {
		// We have a subtag!
		inSubTag = true
		subStmt := &ast.TagStatement{
			Token:   p.currToken,
			TagName: currTagName,
			Owner:   stmt,
		}

		stmt.SubTags = append(stmt.SubTags, subStmt)
		currentParseConfig = stmt.SubTagConfig(currTagName)

		// From here on out, the subtag now behaves as a tag in its own right,
		// but is not pushed onto the stack so further sub-tags can be applied.
		stmt = subStmt
	} else if nestedTag := late.FindTag(currTagName); nestedTag != nil {
		// We actually are starting a new tag in the nested context of the current tag.
		// Build our new tag statement and make it the current
		stmt = &ast.TagStatement{
			Token:   p.currToken,
			TagName: currTagName,
			Tag:     nestedTag,
		}

		p.pushCurrentTag(stmt)
		currentParseConfig = stmt.Tag.Parse()
	} else {
		p.parserErrorf("Unknown tag '%s'", stmt.TagName)
		return nil
	}

	for _, parseRule := range currentParseConfig.Rules {
		expectedTokenType := p.parseRuleToTokenType(parseRule)

		if p.peekTokenIs(token.CLOSE_TAG) || p.peekTokenIs(token.EOF) {
			p.parserErrorf("Error parsing tag '%s': expected %s", stmt.TagName, expectedTokenType)
			break
		}

		p.nextToken()

		if expectedTokenType != token.EXPRESSION && !p.currTokenIs(expectedTokenType) {
			p.parserErrorf("Error parsing tag '%s': expected %s found %s", stmt.TagName, expectedTokenType, p.currToken.Type)
			break
		}

		switch parseRule := parseRule.(type) {
		case *tag.IdentifierRule:
			stmt.Nodes = append(stmt.Nodes, &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})
		case *tag.LiteralRule:
			stmt.Nodes = append(stmt.Nodes, &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal})
		case *tag.ExpressionRule:
			stmt.Nodes = append(stmt.Nodes, p.parseExpression(LOWEST))
		case *tag.TokenRule:
			stmt.Nodes = append(stmt.Nodes, &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal})
		default:
			p.parserErrorf("Error parsing tag '%s': Don't know how to handle ParseRule of type %T", stmt.TagName, parseRule)
			return nil
		}
	}

	if !p.expectPeek(token.CLOSE_TAG) {
		return nil
	}

	// Move to our %} token so we can continue
	p.nextToken()

	if currentParseConfig.Block {
		stmt.BlockStatement = p.parseBlockStatement()
	}

	// Done with this tag, pop us on out, and make sure we move
	// ourselves to the final END tag for block tags so the parser
	// can continue it's work.
	if !inSubTag {
		p.popCurrentTag()

		if currentParseConfig.Block {
			if !p.peekTokenIs(token.END) {
				p.parserErrorf("Error parsing tag '%s': expected %s found %s", stmt.TagName, token.END, p.peekToken.Type)
				return nil
			}

			p.nextToken()
		}
	}

	return stmt
}

func (p *Parser) pushCurrentTag(tagStmt *ast.TagStatement) {
	p.currentTagStack = append(p.currentTagStack, tagStmt)
}

func (p *Parser) currentTag() *ast.TagStatement {
	if len(p.currentTagStack) == 0 {
		return nil
	}

	return p.currentTagStack[len(p.currentTagStack)-1]
}

func (p *Parser) popCurrentTag() *ast.TagStatement {
	top := p.currentTag()
	p.currentTagStack = p.currentTagStack[0 : len(p.currentTagStack)-1]

	return top
}

func (p *Parser) parseRuleToTokenType(parseRule tag.ParseRule) token.TokenType {
	switch parseRule := parseRule.(type) {
	case *tag.IdentifierRule:
		return token.IDENT
	case *tag.LiteralRule:
		return token.STRING
	case *tag.TokenRule:
		return parseRule.Type
	case *tag.ExpressionRule:
		return token.EXPRESSION
	default:
		p.parserErrorf("Don't know how to convert parseRule type %T to a token.TokenType", parseRule)
		return token.ILLEGAL
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	currTag := p.currentTag()
	var nextStmt ast.Statement

	for !p.peekTokenIs(token.END) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		nextStmt = p.parseNext()
		// Due to the not-really-nested nature of block tags, we look for
		// any tag statements generated here and if the tag is actually a sub-tag
		// then we need to not include that tag in the block statements list of
		// the parent block tag.
		asTag, ok := nextStmt.(*ast.TagStatement)

		if !ok || asTag.Owner != currTag {
			block.Statements = append(block.Statements, nextStmt)
		}
	}

	return block
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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.currToken,
		Left:  left,
	}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RSQUARE) {
		return nil
	}

	p.nextToken()
	return expr
}

func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.currToken,
		Left:  left,
	}

	p.nextToken()
	// Dot access is syntax sugar for square bracket access,
	// so while this looks like an identifier, we hack it into a String
	// and evaluate it as an index access instead.
	//
	//   this.that.those == this["that"]["those"]
	//
	expr.Index = p.parseStringLiteral()

	return expr
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

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currToken}

	// Empty Array test
	if p.peekTokenIs(token.RSQUARE) {
		p.nextToken()
		return array
	}

	// First element
	p.nextToken()
	array.Expressions = append(array.Expressions, p.parseExpression(LOWEST))

	// Rest of the elements
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		array.Expressions = append(array.Expressions, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RSQUARE) {
		return nil
	}

	p.nextToken()

	return array
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
	msg := fmt.Sprintf(message, args...)
	p.Errors = append(p.Errors, msg)
}
