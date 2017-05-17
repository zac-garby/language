package token

const (
	// Misc
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NL      = "NEWLINE"
	COMMENT = "COMMENT"

	// Values
	ID     = "ID"
	NUM    = "NUM"
	STRING = "STRING"

	// Operators
	DECLARE   = ":="
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	STAR      = "*"
	SLASH     = "/"
	MOD       = "%"
	BACKSLASH = "\\"
	BANG      = "!"
	LT        = "<"
	GT        = ">"
	LTE       = "<="
	GTE       = ">="
	EQ        = "=="
	NOT_EQ    = "!="
	RANGE     = ".."
	XRANGE    = "..<"
	AND       = "&&"
	OR        = "||"
	EXP       = "**"
	BIT_LEFT  = "<<"
	BIT_RIGHT = ">>"
	BIT_AND   = "&"
	BIT_XOR   = "^"
	BIT_NOT   = "~"
	IN        = "IN"

	// Separators
	COMMA   = ","
	SEMI    = ";"
	COLON   = ":"
	DOT     = "."
	VLINE   = "|"
	NEWLINE = "NEWLINE"

	// Parentheses
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	MODEL    = "MODEL"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"

	// Conditional keywords
	IF   = "IF"
	ELSE = "ELSE"
	ELIF = "ELIF"

	// Looping keywords
	WHILE = "WHILE"
	FOR   = "FOR"
	BREAK = "BREAK"
	NEXT  = "NEXT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"elif":   ELIF,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
	"break":  BREAK,
	"next":   NEXT,
	"null":   NULL,
	"in":     IN,
	"model":  MODEL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ID
}
