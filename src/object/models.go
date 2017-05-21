package object

import (
	"../ast"
	"../token"
	"math"
)

func newID(val string) *ast.Identifier {
	return &ast.Identifier{Token: token.New(token.ID, val), Value: val}
}

var (
	OBJECT_MODEL = &Model{
		Parent:     nil,
		Properties: []*ast.Identifier{},
		Methods:    map[*ast.Identifier]Object{},
		Id:         -1,
	}

	VECTOR_MODEL = &Model{
		Parent:     OBJECT_MODEL,
		Properties: []*ast.Identifier{newID("x"), newID("y")},
		Methods:    map[*ast.Identifier]Object{},
		Id:         -2,
	}
)

func InitialiseBuiltinModels() bool {
	OBJECT_MODEL.Methods = map[*ast.Identifier]Object{
		newID("type"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 0 {
					return newError("no arguments expected to object.type")
				}

				thisHash := this.(*Hash)

				return thisHash.Model
			},
		},
		newID("parent"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 0 {
					return newError("no arguments expected to object.type")
				}

				thisHash := this.(*Hash)

				if thisHash.Model.Parent != nil {
					return thisHash.Model.Parent
				} else {
					return &Null{}
				}
			},
		},
	}

	VECTOR_MODEL.Methods = map[*ast.Identifier]Object{
		newID("_new"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 0 {
					return newError("no arguments expected to vec._new")
				}

				thisHash := this.(*Hash)

				if _, ok := thisHash.Get("x").(*Number); !ok {
					return newError("the 'x' property of a vec must be a number")
				}

				if _, ok := thisHash.Get("y").(*Number); !ok {
					return newError("the 'y' property of a vec must be a number")
				}

				return thisHash
			},
		},
		newID("_plus"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 1 {
					return newError("expected exactly one argument to vec._plus")
				}

				left := this.(*Hash)

				right, ok := args[0].(*Hash)
				if !ok {
					return newError("expected another hash as the only argument to vec._plus")
				}

				if !left.Model.Equals(right.Model) {
					return newError("expected the first argument of vec._plus to be another hash")
				}

				leftX := left.Get("x").(*Number).Value
				leftY := left.Get("y").(*Number).Value

				rightX := right.Get("x").(*Number).Value
				rightY := right.Get("y").(*Number).Value

				return VECTOR_MODEL.Instantiate([]Object{
					&Number{Value: leftX + rightX},
					&Number{Value: leftY + rightY},
				})
			},
		},
		newID("_mul"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 1 {
					return newError("expected exactly one argument to vec._plus")
				}

				left := this.(*Hash)

				right, ok := args[0].(*Hash)
				if !ok {
					return newError("expected another hash as the only argument to vec._plus")
				}

				if !left.Model.Equals(right.Model) {
					return newError("expected the first argument of vec._plus to be another vector")
				}

				leftX := left.Get("x").(*Number).Value
				leftY := left.Get("y").(*Number).Value

				rightX := right.Get("x").(*Number).Value
				rightY := right.Get("y").(*Number).Value

				return VECTOR_MODEL.Instantiate([]Object{
					&Number{Value: leftX + rightX},
					&Number{Value: leftY + rightY},
				})
			},
		},
		newID("len"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 0 {
					return newError("no arguments expected to vec._new")
				}

				hash := this.(*Hash)

				x := hash.Get("x").(*Number).Value
				y := hash.Get("y").(*Number).Value

				return &Number{Value: math.Sqrt(x*x + y*y)}
			},
		},
		newID("translate"): &Builtin{
			Fn: func(this Object, args ...Object) Object {
				if len(args) != 1 {
					return newError("expected one or two arguments to vec.translate")
				}

				other, ok := args[0].(*Hash)
				if !ok {
					return newError("expected the first argument of vec.translate to be a hash")
				}

				hash := this.(*Hash)

				if !other.Model.Equals(hash.Model) {
					return newError("expected the first argument of vec.translate to be another vector")
				}

				thisX := hash.Get("x").(*Number).Value
				thisY := hash.Get("y").(*Number).Value

				otherX := other.Get("x").(*Number).Value
				otherY := other.Get("y").(*Number).Value

				newX := &Number{Value: thisX + otherX}
				newY := &Number{Value: thisY + otherY}

				hash.Set("x", newX)
				hash.Set("y", newY)

				return VECTOR_MODEL.Instantiate([]Object{newX, newY})
			},
		},
	}

	return true
}

var DefaultModels = map[string]*Model{
	"object": OBJECT_MODEL,
	"vec":    VECTOR_MODEL,
}

var _ = InitialiseBuiltinModels()
