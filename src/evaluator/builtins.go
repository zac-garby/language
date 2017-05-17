package evaluator

import (
	"../object"
	"bufio"
	"fmt"
	"os"
)

var builtins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect() + " ")
			}

			fmt.Println()

			return NULL
		},
	},
	"err": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			msg := ""

			for _, arg := range args {
				msg += arg.Inspect() + " "
			}

			return newError(msg)
		},
	},
	"str": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'str'")
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},
	"input": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'input'")
			}

			arg := args[0]

			reader := bufio.NewReader(os.Stdin)
			fmt.Print(arg.Inspect())

			text, err := reader.ReadString('\n')
			if err != nil {
				return newError("could not read a line")
			}

			return &object.String{Value: text}
		},
	},
	"type": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'type'")
			}

			arg := args[0]

			hash, ok := arg.(*object.Hash)
			if !ok {
				return newError("expected a hash to be passed to 'type'")
			}

			return hash.Model
		},
	},
	"parent": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'type'")
			}

			arg := args[0]

			hash, ok := arg.(*object.Hash)
			if !ok {
				return newError("expected a hash to be passed to 'type'")
			}

			parent := hash.Model.Parent
			if parent == nil {
				return NULL
			}

			return parent
		},
	},
}
