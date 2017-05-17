package repl

import (
	"../evaluator"
	"../lexer"
	"../object"
	"../parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// io.WriteString(out, " -> "+program.String()+"\n")

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, " => "+evaluated.Inspect()+"\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "  "+msg+"\n")
	}
}
