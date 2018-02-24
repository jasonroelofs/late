package parser

type Parser struct {
	input string
}

func New(input string) *Parser {
	return &Parser{input: input}
}

func (p *Parser) Parse() string {
	return p.input
}
