package parser

import (
	"../lexer"
	"testing"
)

type test struct {
	input          string
	expectedOutput string
}

func TestID(t *testing.T) {
	runTests(t, []test{
		{"hello;", "hello"},
		{"hello_world;", "hello_world"},
		{"abc9;", "abc9"},
		{"_x;", "_x"},
		{"1a", "ERROR: expected next token to be ;, but got ID"},
	})
}

func TestLiterals(t *testing.T) {
	runTests(t, []test{
		{"5;", "5"},
		{"5.3;", "5.3"},
		{"1.;", "1"},
		{"true;", "true"},
		{"false;", "false"},
		{`"Hello, world";`, "\"Hello, world\""},
		{`"Hello\nWorld";`, "\"Hello\nWorld\""},
		{"null;", "null"},
	})
}

func TestUnaryOps(t *testing.T) {
	runTests(t, []test{
		{"-5;", "(-5)"},
		{"!3;", "(!3)"},
		{"-++--1;", "(-(+(+(-(-1)))))"},
		{"=5", "ERROR: no prefix parse function for = found"},
	})
}

func TestInfixOps(t *testing.T) {
	runTests(t, []test{
		{"a + b;", "(a + b)"},
		{"a - b;", "(a - b)"},
		{"a / b;", "(a / b)"},
		{"a * b;", "(a * b)"},
		{"a % b;", "(a % b)"},
		{"a == b;", "(a == b)"},
		{"a != b;", "(a != b)"},
		{"a < b;", "(a < b)"},
		{"a > b;", "(a > b)"},
		{"a >= b;", "(a >= b)"},
		{"a <= b;", "(a <= b)"},
		{"a.b;", "(a . b)"},
		{"a := b;", "(a := b)"},
		{"a = b;", "(a = b)"},
		{"a..b;", "(a .. b)"},
		{"a..<b;", "(a ..< b)"},
		{"a || b;", "(a || b)"},
		{"a && b;", "(a && b)"},
		{"a ** b;", "(a ** b)"},
		{"a << b;", "(a << b)"},
		{"a >> b;", "(a >> b)"},
		{"a & b;", "(a & b)"},
		{"a | b;", "(a | b)"},
		{"a ^ b;", "(a ^ b)"},
		{"a in b;", "(a in b)"},
	})
}

func TestIfExpr(t *testing.T) {
	runTests(t, []test{
		{"if cond { a + b; };", "(if cond (a + b))"},
		{"if cond { a + b; } else { a - b; };", "(if cond (a + b) else (a - b))"},
		{
			"if cond { a + b; } elif otherCond { a - b; } elif another { a * b; } else { a / b; };",
			"(if cond (a + b) else (if otherCond (a - b) else (if another (a * b) else (a / b))))",
		},
	})
}

func TestLoops(t *testing.T) {
	runTests(t, []test{
		{"for i | array { i + 1; };", "(for (i | array) { (i + 1) })"},
		{"for (i | array) { i + 1; };", "(for (i | array) { (i + 1) })"},
		{"while cond { a + b; };", "(while cond { (a + b) })"},
		{"while (cond) { a + b; };", "(while cond { (a + b) })"},
	})
}

func TestFunctions(t *testing.T) {
	runTests(t, []test{
		{"fn (x, y) {};", "(fn(x, y) )"},
		{"\\(x, y) = x;", "(\\(x, y) = x)"},
		{"fn (x, y) { x + y; };", "(fn(x, y) (x + y))"},
		{
			"fn (x, y) { if x > 0 { return; }; x * y; };",
			"(fn(x, y) (if (x > 0) return ;)(x * y))",
		},
		{`\(x, y) = x * y;`, `(\(x, y) = (x * y))`},
	})
}

func TestModels(t *testing.T) {
	runTests(t, []test{
		{"model (x, y);", "(model (x, y))"},
		{`model (x, y) : p (a, "y");`, `(model (x, y) : model (a, "y"))`},
	})
}

func runTests(t *testing.T, tests []test) {
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		node := p.ParseProgram()

		if len(test.expectedOutput) > 5 && test.expectedOutput[0:5] == "ERROR" {
			if len(p.Errors()) == 0 {
				t.Errorf("expected at least one error when parsing %v",
					test.input)
			}
		} else {
			if node.String() != test.expectedOutput {
				t.Errorf("expected %v but got %v",
					test.expectedOutput, node.String())
			}
		}
	}
}
