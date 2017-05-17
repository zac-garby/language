package lexer

import (
	"../token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `id 3.27 "string" := = + - * / \ !
  < > <= >= == != .. ..< && || , ; : . |
  (){}[] fn return true false if else elif while for break next`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ID, "id"},
		{token.NUM, "3.27"},
		{token.STRING, "string"},
		{token.DECLARE, ":="},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.BACKSLASH, "\\"},
		{token.BANG, "!"},
		{token.NL, "\n"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.RANGE, ".."},
		{token.XRANGE, "..<"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.COMMA, ","},
		{token.SEMI, ";"},
		{token.COLON, ":"},
		{token.DOT, "."},
		{token.VLINE, "|"},
		{token.NL, "\n"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.FUNCTION, "fn"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.ELIF, "elif"},
		{token.WHILE, "while"},
		{token.FOR, "for"},
		{token.BREAK, "break"},
		{token.NEXT, "next"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
