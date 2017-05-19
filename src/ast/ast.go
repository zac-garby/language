package ast

import (
	"../token"
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program node

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier expression

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Number literal

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (il *NumberLiteral) expressionNode()      {}
func (il *NumberLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *NumberLiteral) String() string       { return il.Token.Literal }

// Boolean literal

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// Null literal

type Null struct {
	Token token.Token
}

func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return n.Token.Literal }
func (n *Null) String() string       { return n.Token.Literal }

// Expression statement

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Declare statement

type DeclareExpression struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (ds *DeclareExpression) expressionNode()      {}
func (ds *DeclareExpression) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeclareExpression) String() string {
	return fmt.Sprintf("(%v := %v)", ds.Name.String(), ds.Value.String())
}

// Assign statement

type AssignExpression struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (as *AssignExpression) expressionNode()      {}
func (as *AssignExpression) TokenLiteral() string { return as.Token.Literal }
func (as *AssignExpression) String() string {
	return fmt.Sprintf("(%v = %v)", as.Name.String(), as.Value.String())
}

// Return statement

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %v;", rs.ReturnValue.String())
}

// Prefix expression

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%v%v)", pe.Operator, pe.Right.String())
}

// Infix expression

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%v %v %v)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// If expression

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String() + ")"
}

// Block statement

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Function literal

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return "(" + out.String() + ")"
}

// Model literal

type ModelLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	ParentName *Expression
	ParentArgs []Expression
}

func (ml *ModelLiteral) expressionNode()      {}
func (ml *ModelLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *ModelLiteral) String() string {
	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	if ml.ParentName != nil {
		parentArgs := []string{}
		for _, p := range ml.ParentArgs {
			parentArgs = append(parentArgs, p.String())
		}

		return fmt.Sprintf("(model (%v) : model (%v))",
			strings.Join(params, ", "), strings.Join(parentArgs, ", "))
	}

	return fmt.Sprintf("(model (%v))",
		strings.Join(params, ", "))
}

// Lambda expression

type LambdaExpression struct {
	Token      token.Token
	Parameters []*Identifier
	Body       Expression
}

func (le *LambdaExpression) expressionNode()      {}
func (le *LambdaExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LambdaExpression) String() string {
	params := []string{}
	for _, p := range le.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("(\\(%v) = %v)",
		strings.Join(params, ", "), le.Body.String())
}

// Call expression

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return "(" + out.String() + ")"
}

// String Literal

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return `"` + sl.Token.Literal + `"` }

// Array Literal

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Index expression

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// Hash literal

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// While expression

type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" { ")
	out.WriteString(we.Body.String())
	out.WriteString(" })")

	return out.String()
}

// For expression

type ForExpression struct {
	Token token.Token
	Var   Expression
	Set   Expression
	Body  *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	return fmt.Sprintf("(for (%v | %v) { %v })",
		fe.Var.String(), fe.Set.String(), fe.Body.String())
}

// Break statement

type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.Token.Literal + ";" }

// Next statement

type NextStatement struct {
	Token token.Token
}

func (ns *NextStatement) statementNode()       {}
func (ns *NextStatement) TokenLiteral() string { return ns.Token.Literal }
func (ns *NextStatement) String() string       { return ns.Token.Literal + ";" }
