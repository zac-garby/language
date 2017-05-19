package object

import (
	"../ast"
	"fmt"
	"strings"
)

type ObjectType string
type BuiltinFunction func(this Object, args ...Object) Object

const (
	ERROR_OBJ                  = "ERROR"
	NUMBER_OBJ                 = "NUMBER"
	BOOLEAN_OBJ                = "BOOLEAN"
	STRING_OBJ                 = "STRING"
	BUILTIN_OBJ                = "BUILTIN"
	ARRAY_OBJ                  = "ARRAY"
	HASH_OBJ                   = "HASH"
	NULL_OBJ                   = "NULL"
	FUNCTION_OBJ               = "FUNCTION"
	LAMBDA_OBJ                 = "LAMBDA"
	RETURN_VALUE_OBJ           = "RETURN_VALUE"
	LOOP_CONTROL_STATEMENT_OBJ = "LOOP_CONTROL_STATEMENT"
	MODEL_OBJ                  = "MODEL"
	METHOD_INSTANCE_OBJ        = "METHOD_INSTANCE"
)

// Object interface

type Object interface {
	Type() ObjectType
	Inspect() string
	Equals(Object) bool
}

// Error

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Equals(other Object) bool {
	switch other := other.(type) {
	case *Error:
		return e.Message == other.Message
	default:
		return false
	}
}

// Number

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%v", n.Value) }
func (n *Number) Equals(other Object) bool {
	switch other := other.(type) {
	case *Number:
		return n.Value == other.Value
	default:
		return false
	}
}

func (n *Number) IsInteger() bool {
	return float64(int64(n.Value)) == n.Value
}

func (n *Number) IsPositive() bool {
	return n.Value >= 0
}

// Boolean

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Equals(other Object) bool {
	switch other := other.(type) {
	case *Boolean:
		return b.Value == other.Value
	default:
		return false
	}
}

// String

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Equals(other Object) bool {
	switch other := other.(type) {
	case *String:
		return s.Value == other.Value
	default:
		return false
	}
}

// Null

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "<null>" }
func (n *Null) Equals(other Object) bool {
	switch other.(type) {
	case *Null:
		return true
	default:
		return false
	}
}

// ReturnValue

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Equals(other Object) bool {
	switch other := other.(type) {
	case *ReturnValue:
		return rv.Value.Equals(other)
	default:
		return false
	}
}

// Loop control statement

type LoopControlStatement struct {
	Literal string
}

func (lcs *LoopControlStatement) Type() ObjectType {
	return LOOP_CONTROL_STATEMENT_OBJ
}
func (lcs *LoopControlStatement) Inspect() string {
	return lcs.Literal
}
func (lcs *LoopControlStatement) Equals(other Object) bool {
	switch other := other.(type) {
	case *LoopControlStatement:
		return lcs.Literal == other.Literal
	default:
		return false
	}
}

// Function

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("fn (%v) { %v }",
		strings.Join(params, ", "), f.Body.String())
}
func (f *Function) Equals(other Object) bool {
	switch other := other.(type) {
	case *Function:
		return f.Inspect() == other.Inspect()
	default:
		return false
	}
}

// Lambda

type Lambda struct {
	Parameters []*ast.Identifier
	Body       *ast.Expression
	Env        *Environment
}

func (l *Lambda) Type() ObjectType { return LAMBDA_OBJ }
func (l *Lambda) Inspect() string {
	params := []string{}
	for _, p := range l.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("\\(%v) = %v",
		strings.Join(params, ", "), (*l.Body).String())
}
func (l *Lambda) Equals(other Object) bool {
	switch other := other.(type) {
	case *Lambda:
		return l.Inspect() == other.Inspect()
	default:
		return false
	}
}

// Builtin

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType         { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string          { return "<builtin>" }
func (b *Builtin) Equals(other Object) bool { return false }

// Array

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elems := []string{}
	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}

	str := fmt.Sprintf("[%v]", strings.Join(elems, ", "))

	return str
}
func (a *Array) Equals(other Object) bool {
	switch other := other.(type) {
	case *Array:
		for i, e := range a.Elements {
			if !e.Equals(other.Elements[i]) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
