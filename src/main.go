package main

import (
	"./evaluator"
	"./lexer"
	"./object"
	"./parser"
	"./repl"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func main() {
	if len(os.Args) > 1 {
		runFile()
	} else {
		startREPL()
	}
}

func runFile() {
	fileName := os.Args[1]

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}

	text := string(bytes)

	env := object.NewEnvironment()
	l := lexer.New(text)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	result := evaluator.Eval(program, env)
	if result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		return
	}
}

func printParserErrors(errors []string) {
	fmt.Println("parser errors:")
	for _, msg := range errors {
		fmt.Println("  " + msg)
	}
}

func startREPL() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
