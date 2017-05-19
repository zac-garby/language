package evaluator

import (
	"../object"
	"bufio"
	"fmt"
	"os"
	"time"
)

var builtins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect() + " ")
			}

			fmt.Println()

			return NULL
		},
	},
	"err": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
			msg := ""

			for _, arg := range args {
				msg += arg.Inspect() + " "
			}

			return newError(msg)
		},
	},
	"str": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'str'")
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},
	"input": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
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

			return &object.String{Value: text[0 : len(text)-1]}
		},
	},
	"type": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
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
		Fn: func(this object.Object, args ...object.Object) object.Object {
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
	"sleep": &object.Builtin{
		Fn: func(this object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected exactly one argument to 'sleep'")
			}

			arg := args[0]

			seconds, ok := arg.(*object.Number)
			if !ok {
				return newError("expected a number to be passed to 'sleep'")
			}

			time.Sleep(time.Duration(seconds.Value) * time.Second)

			return NULL
		},
	},
}
