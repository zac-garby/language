package object

import "../ast"

func newID(val string) *ast.Identifier {
	return &ast.Identifier{Value: val}
}

var (
	OBJECT_MODEL = &Model{
		Properties: []*ast.Identifier{},
		Methods:    map[ast.Identifier]*Function{},
		Id:         -1,
	}
)

var DefaultModels = map[string]*Model{
	"Object": OBJECT_MODEL,
}
