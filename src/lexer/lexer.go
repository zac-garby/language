package lexer

import (
	"../token"
	"strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// skip comments and whitespace
	for {
		if l.ch == '#' {
			l.skipComment()
		} else if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			l.skipWhitespace()
		} else {
			break
		}
	}

	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.DECLARE, ":=")
		} else {
			tok = token.New(token.COLON, ":")
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.EQ, "==")
		} else {
			tok = token.New(token.ASSIGN, "=")
		}
	case '+':
		tok = token.New(token.PLUS, "+")
	case '-':
		tok = token.New(token.MINUS, "-")
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			tok = token.New(token.EXP, "**")
		} else {
			tok = token.New(token.STAR, "*")
		}
	case '/':
		tok = token.New(token.SLASH, "/")
	case '%':
		tok = token.New(token.MOD, "%")
	case '\\':
		tok = token.New(token.BACKSLASH, "\\")
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.NOT_EQ, "!=")
		} else {
			tok = token.New(token.BANG, "!")
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LTE, "<=")
		} else if l.peekChar() == '<' {
			l.readChar()
			tok = token.New(token.BIT_LEFT, "<<")
		} else {
			tok = token.New(token.LT, "<")
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GTE, ">=")
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token.New(token.BIT_RIGHT, ">>")
		} else {
			tok = token.New(token.GT, ">")
		}
	case '.':
		if l.startsWith("..<") {
			l.readChar()
			l.readChar()
			tok = token.New(token.XRANGE, "..<")
		} else if l.peekChar() == '.' {
			l.readChar()
			tok = token.New(token.RANGE, "..")
		} else {
			tok = token.New(token.DOT, ".")
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.New(token.AND, "&&")
		} else {
			tok = token.New(token.BIT_AND, "&")
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.New(token.OR, "||")
		} else {
			tok = token.New(token.VLINE, "|")
		}
	case '^':
		tok = token.New(token.BIT_XOR, "^")
	case ',':
		tok = token.New(token.COMMA, ",")
	case ';':
		tok = token.New(token.SEMI, ";")
	case '(':
		tok = token.New(token.LPAREN, "(")
	case ')':
		tok = token.New(token.RPAREN, ")")
	case '{':
		tok = token.New(token.LBRACE, "{")
	case '}':
		tok = token.New(token.RBRACE, "}")
	case '[':
		tok = token.New(token.LBRACKET, "[")
	case ']':
		tok = token.New(token.RBRACKET, "]")
	case '~':
		tok = token.New(token.BIT_NOT, "~")
	case '"':
		tok = token.New(token.STRING, l.readString())
	case 0:
		tok = token.New(token.EOF, "")
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.NUM
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) startsWith(str string) bool {
	section := l.input[l.position:len(l.input)]
	return strings.HasPrefix(section, str)
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	dotted := false
	literal := ""

	for {
		if isDigit(l.ch) {
			literal += string(l.ch)
		} else if l.ch == '.' && !dotted && len(literal) > 0 {
			dotted = true

			l.readChar()
			if isDigit(l.ch) {
				literal += "." + string(l.ch)
			} else {
				l.readPosition--
				l.position--
				return literal
			}
		} else {
			break
		}

		l.readChar()
	}

	return literal
}

func (l *Lexer) readString() string {
	str := ""
	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}

		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '\n', '\\':
				str += "\\"
			case '\'':
				str += "'"
			case '"':
				str += "\""
			case 'a':
				str += "\a"
			case 'b':
				str += "\b"
			case 'f':
				str += "\f"
			case 'n':
				str += "\n"
			case 'r':
				str += "\r"
			case 't':
				str += "\t"
			case 'v':
				str += "\v"
			}
		} else {
			str += string(l.ch)
		}
	}

	return str
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?'
}
