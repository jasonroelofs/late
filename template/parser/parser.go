package parser

import (
	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/lexer"
	"github.com/jasonroelofs/late/template/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
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
		stmt := p.parseStatement()
		if stmt != nil {
			template.Statements = append(template.Statements, stmt)
		}

		p.nextToken()
	}

	return template
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.RAW:
		return p.parseRawStatement()
	default:
		return nil
	}
}

func (p *Parser) parseRawStatement() *ast.RawStatement {
	return &ast.RawStatement{Token: p.currToken}
}
