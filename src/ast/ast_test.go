package ast

import (
	"../token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &DeclareExpression{
					Token: token.Token{Type: token.DECLARE, Literal: ":="},
					Name: &Identifier{
						Token: token.Token{Type: token.ID, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &NumberLiteral{
							Token: token.Token{Type: token.NUM, Literal: "5"},
							Value: 5,
						},
						Operator: "+",
						Right: &NumberLiteral{
							Token: token.Token{Type: token.NUM, Literal: "6"},
							Value: 6,
						},
					},
				},
			},
		},
	}

	if program.String() != "(myVar := (5 + 6))" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
