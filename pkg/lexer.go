package veritas

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	VAR     = "VAR"

	OR  = "|"
	AND = "&"
	XOR = "*"
	NOT = "!"
	IMP = "@"

	LPAREN = "("
	RPAREN = ")"
)

// order
const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[string]int{
	OR:     SUM,
	AND:    SUM,
	XOR:    SUM,
	IMP:    PRODUCT,
	NOT:    PREFIX,
	LPAREN: CALL,
}

type Token struct {
	Type    string
	Literal string
}

func newToken(tokenType string, c byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(c),
	}
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLex(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '&':
		tok = newToken(AND, l.ch)
	case '|':
		tok = newToken(OR, l.ch)
	case '^':
		tok = newToken(XOR, l.ch)
	case '!':
		tok = newToken(NOT, l.ch)
	case '@':
		tok = newToken(IMP, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isVar(l.ch) {
			tok.Type = VAR
			tok.Literal = l.readByte()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func isVar(ch byte) bool {
	return 'A' <= ch && ch <= 'z'
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readByte() string {
	position := l.position
	if isVar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
