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
	return l
}

func (l *Lexer) NextToken() (tok token.Token) {
	if l.inLiquid {
		tok = l.parseNextLiquidToken()
	} else {
		tok = l.parseUntilLiquid()
		l.inLiquid = true
	}

	return tok
}

func (l *Lexer) parseUntilLiquid() (tok token.Token) {
	startPos := l.position

	for {
		l.readChar()
		if l.ch == 0 || (l.ch == '{' && (l.peek() == '{' || l.peek() == '%')) {
			break
		}
	}

	tok.Type = token.RAW
	tok.Literal = l.input[startPos:l.position]

	return
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
		if isLetter(l.ch) {
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

func (l *Lexer) readIdentifier() string {
	startPosition := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) readString() string {
	// Move us past the quote mark that triggered this step
	l.readChar()

	startPosition := l.position
	for l.ch != '"' && l.ch != '\'' {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
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
