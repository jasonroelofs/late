package lexer

import (
	"github.com/jasonroelofs/late/template/token"
)

type Lexer struct {
	input string

	// Current position in the input string
	position int

	// Next position in the input string
	peekPosition int

	// The currenty character under advisulment.
	ch byte

	// Are we parsing actual Liquid?
	inLiquid bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	l.checkStartState()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.ch == 0 {
		tok.Type = token.EOF
		return tok
	}

	if l.inLiquid {
		tok = l.parseNextLiquidToken()
	} else {
		tok = l.parseUntilLiquid()
		l.inLiquid = true
	}

	return tok
}

func (l *Lexer) parseUntilLiquid() token.Token {
	var tok token.Token
	startPos := l.position

	for {
		if l.atLiquidStart() {
			break
		}

		if l.ch == 0 {
			break
		}

		l.readChar()
	}

	tok.Type = token.RAW
	tok.Literal = l.input[startPos:l.position]

	return tok
}

func (l *Lexer) parseNextLiquidToken() (tok token.Token) {
	l.skipWhitespace()

	switch l.ch {
	case '{':
		if l.peek() == '{' {
			tok.Type = token.OPEN_VAR
			tok.Literal = "{{"
			l.readChar()
		} else if l.peek() == '%' {
			tok.Type = token.OPEN_TAG
			tok.Literal = "{%"
			l.readChar()
		} else {
			tok = newToken(token.LBRACKET, l.ch)
		}
	case '}':
		if l.peek() == '}' {
			tok.Type = token.CLOSE_VAR
			tok.Literal = "}}"

			l.readChar()
			l.inLiquid = false
		} else {
			tok = newToken(token.RBRACKET, l.ch)
		}
	case '%':
		if l.peek() == '}' {
			tok.Type = token.CLOSE_TAG
			tok.Literal = "%}"

			l.readChar()
			l.inLiquid = false
		} else {
			tok = newToken(token.PERCENT, l.ch)
		}

	case '[':
		tok = newToken(token.LSQUARE, l.ch)
	case ']':
		tok = newToken(token.RSQUARE, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '|':
		tok = newToken(token.PIPE, l.ch)
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isNumber(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return
		} else if isIdentifier(l.ch) {
			tok.Type = token.IDENT
			tok.Literal = l.readIdentifier()
			return
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()

	return
}

func (l *Lexer) readNumber() string {
	startPosition := l.position

	// We don't currently support starting a float number
	// without a leading number
	notFirst := false

	for isNumber(l.ch) || (notFirst && l.ch == '.') {
		l.readChar()
		notFirst = true
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.position

	for isIdentifier(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) readString() string {
	// Keep track of what character opened our string (' or ")
	// so we can find our matching closing quote,
	// making sure that we look for escaped quotes and don't prematurely
	// end our string parsing.
	prev := l.ch
	openWith := l.ch
	l.readChar()

	startPosition := l.position
	for {
		if (l.ch == openWith && prev != '\\') || l.ch == 0 {
			break
		}

		prev = l.ch
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) checkStartState() {
	l.inLiquid = l.atLiquidStart()
}

func (l *Lexer) atLiquidStart() bool {
	return (l.ch == '{' && (l.peek() == '{' || l.peek() == '%'))
}

func (l *Lexer) readChar() {
	if l.peekPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.peekPosition]
	}

	l.position = l.peekPosition
	l.peekPosition += 1
}

func (l *Lexer) peek() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isIdentifier(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
