package parser

import "github.com/jasonroelofs/late/template/lexer"

type Parser struct {
	l *lexer.Lexer
}

func New(lexer *lexer.Lexer) *Parser {
	return &Parser{l: lexer}
}

func (p *Parser) Parse() string {
	return ""
}
