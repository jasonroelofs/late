package lexer

import (
	"github.com/jasonroelofs/late/template/token"
)

type Lexer struct {
	input       string
	eofPosition int

	tokenStart   int
	lookPosition int

	// Are we parsing actual Late?
	// For anything outside of Late tags, we want to combine all text
	// into a single Raw token that can be trivially included back in the
	// expected output.
	inCode bool
}

func New(input string) *Lexer {
	return &Lexer{
		input:        input,
		eofPosition:  len(input) - 1,
		tokenStart:   0,
		lookPosition: 0,
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// With the assumption that the end of the previous token will include
	// all of the characters from l.tokenStart to l.lookPosition (exclusive),
	// the next token obviously starts at the current look position.
	l.tokenStart = l.lookPosition

	if l.atEOF() {
		tok.Type = token.EOF
		return tok
	}

	// Handle the case where two code tags are right next to
	// each other, e.g. {{ 1 }}{{ 2 }}. We don't want an empty
	// intermediate RAW token inbetween them.
	if !l.inCode && l.atCodeStart() {
		l.inCode = true
	}

	if l.inCode {
		tok = l.parseNextCodeToken()
	} else {
		tok = l.parseUntilCode()
		l.inCode = true
	}

	return tok
}

func (l *Lexer) parseNextCodeToken() (tok token.Token) {
	l.skipWhitespace()

	switch {
	case l.test("{%end%}"):
		l.inCode = false
		tok = l.stringToken(token.END)
	case l.test("{%"):
		tok = l.stringToken(token.OPEN_TAG)
	case l.test("%}"):
		l.inCode = false
		tok = l.stringToken(token.CLOSE_TAG)
	case l.test("{{"):
		tok = l.stringToken(token.OPEN_VAR)
	case l.test("}}"):
		l.inCode = false
		tok = l.stringToken(token.CLOSE_VAR)
	case l.test(">="):
		tok = l.stringToken(token.GT_EQ)
	case l.test("<="):
		tok = l.stringToken(token.LT_EQ)
	case l.test("=="):
		tok = l.stringToken(token.EQ)
	case l.test("!="):
		tok = l.stringToken(token.NOT_EQ)
	}

	if tok.Type != "" {
		return
	}

	switch l.peek() {
	case '{':
		tok = l.charToken(token.LBRACKET)
	case '}':
		tok = l.charToken(token.RBRACKET)
	case '[':
		tok = l.charToken(token.LSQUARE)
	case ']':
		tok = l.charToken(token.RSQUARE)
	case '.':
		tok = l.charToken(token.DOT)
	case ',':
		tok = l.charToken(token.COMMA)
	case ':':
		tok = l.charToken(token.COLON)
	case '|':
		tok = l.charToken(token.PIPE)
	case '-':
		tok = l.charToken(token.MINUS)
	case '+':
		tok = l.charToken(token.PLUS)
	case '*':
		tok = l.charToken(token.TIMES)
	case '/':
		tok = l.charToken(token.SLASH)
	case '(':
		tok = l.charToken(token.LPAREN)
	case ')':
		tok = l.charToken(token.RPAREN)
	case '>':
		tok = l.charToken(token.GT)
	case '<':
		tok = l.charToken(token.LT)
	case '=':
		tok = l.charToken(token.ASSIGN)
	case '"', '\'':
		tok = l.manualToken(token.STRING, l.readString())
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isNumber(l.peek()) {
			tok = l.manualToken(token.NUMBER, l.readNumber())
			return
		} else if isIdentifier(l.peek()) {
			tok = l.manualToken(token.IDENT, l.readIdentifier())

			switch tok.Literal {
			case "true":
				tok.Type = token.TRUE
			case "false":
				tok.Type = token.FALSE
			}
			return
		} else {
			tok = l.charToken(token.ILLEGAL)
		}
	}

	return
}

func (l *Lexer) parseUntilCode() token.Token {
	for {
		if l.atEOF() {
			break
		}

		if l.atCodeStart() {
			break
		}

		l.step()
	}

	return l.stringToken(token.RAW)
}

func (l *Lexer) test(expect string) bool {
	peekWas := l.lookPosition

	for i := 0; i < len(expect); i++ {
		l.skipWhitespace()

		if l.peek() != expect[i] {
			l.lookPosition = peekWas
			return false
		}

		l.step()
	}

	return true
}

func (l *Lexer) skipWhitespace() {
	for l.isWhitespace(l.peek()) {
		l.step()
	}
}

func (l *Lexer) isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (l *Lexer) atEOF() bool {
	return l.lookPosition > l.eofPosition
}

func (l *Lexer) atCodeStart() bool {
	// Manual checking our characters from the input string here
	// as we don't want to step the lookPosition
	if l.lookPosition == l.eofPosition {
		return false
	}

	at := l.input[l.lookPosition]
	peek := l.input[l.lookPosition+1]

	return at == '{' && (peek == '{' || peek == '%')
}

func (l *Lexer) peek() byte {
	return l.input[l.lookPosition]
}

func (l *Lexer) step() {
	l.lookPosition += 1
}

func (l *Lexer) readString() string {
	// Keep track of what character opened our string (' or ")
	// so we can find our matching closing quote,
	// making sure that we look for escaped quotes and don't prematurely
	// end our string parsing.
	prev := l.peek()
	openWith := l.peek()
	var buffer []byte
	l.step()

	for {
		if (l.peek() == openWith && prev != '\\') || l.atEOF() {
			break
		}

		prev = l.peek()
		buffer = append(buffer, l.peek())
		l.step()
	}

	// And finally move onto the closing quote
	l.step()

	return string(buffer)
}

func (l *Lexer) readNumber() string {
	// We don't currently support starting a float number
	// without a leading number, e.g. don't support ".1234".
	notFirst := false
	var buffer []byte

	for isNumber(l.peek()) || (notFirst && l.peek() == '.') {
		buffer = append(buffer, l.peek())
		l.step()
		notFirst = true
	}

	return string(buffer)
}

func (l *Lexer) readIdentifier() string {
	var buffer []byte

	for isIdentifier(l.peek()) {
		buffer = append(buffer, l.peek())
		l.step()
	}

	return string(buffer)
}

// Build a token out of the current single character
func (l *Lexer) charToken(t token.TokenType) token.Token {
	char := string(l.peek())

	tok := token.Token{
		Type:    t,
		Literal: char,
		Raw:     char,
	}

	// If we're here, it means that lookPosition is the character we need to store
	// in the token so to fit our rule the lookPosition is always at the start
	// of the next token, kick us forward.
	l.step()

	return tok
}

// Build a token out of the full string between tokenStart and lookPosition
func (l *Lexer) stringToken(t token.TokenType) token.Token {
	var buffer []byte

	for i := l.tokenStart; i < l.lookPosition; i++ {
		if l.isWhitespace(l.input[i]) {
			continue
		}

		buffer = append(buffer, l.input[i])
	}

	tok := l.manualToken(t, string(buffer))

	if tok.Type == token.RAW {
		tok.Literal = tok.Raw
	}

	return tok
}

func (l *Lexer) manualToken(t token.TokenType, literal string) token.Token {
	tok := token.Token{
		Type:    t,
		Literal: literal,
		Raw:     l.input[l.tokenStart:l.lookPosition],
	}

	return tok
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
