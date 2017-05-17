package object

import (
	"../ast"
	"fmt"
	"strings"
)

var nextModelId int64 = 0

type Model struct {
	Parent     *Model
	ParentArgs []ast.Expression
	Properties []*ast.Identifier
	Methods    map[ast.Identifier]*Function
	Id         int64
}

func NewModel() *Model {
	return NewModelWithParent(OBJECT_MODEL)
}

func NewModelWithParent(parent *Model) *Model {
	nextModelId += 1

	model := &Model{
		Parent:     parent,
		Properties: []*ast.Identifier{},
		Methods:    make(map[ast.Identifier]*Function),
		Id:         nextModelId,
	}

	return model
}

func (m *Model) Instantiate(args []Object) Object {
	hash := NewHash(m)

	for i, prop := range m.Properties {
		hash.Set(prop.Value, args[i])
	}

	return hash
}

func (m *Model) GetMethod(name string) (*MethodInstance, bool) {
	for k, v := range m.Methods {
		if k.Value == name {
			return &MethodInstance{Function: v}, true
		}
	}

	if m.Parent != nil {
		if parentMethod, ok := m.Parent.GetMethod(name); ok {
			return parentMethod, true
		}
	}

	return nil, false
}

func (m *Model) Type() ObjectType { return MODEL_OBJ }

func (m *Model) Inspect() string {
	props := []string{}
	for _, p := range m.Properties {
		props = append(props, p.Value)
	}

	if m.Parent != nil {
		return fmt.Sprintf("model (%v) : (%v)", strings.Join(props, ", "), m.Parent.Inspect())
	} else {
		return fmt.Sprintf("model (%v)", strings.Join(props, ", "))
	}
}

func (m *Model) Equals(other Object) bool {
	switch other := other.(type) {
	case *Model:
		return m.Id == other.Id
	default:
		return false
	}
}
