package object

import (
	"fmt"
	"strings"
)

type Hash struct {
	Pairs map[String]Object
	Model *Model
}

func NewHash(m *Model) *Hash {
	return &Hash{Pairs: make(map[String]Object), Model: m}
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	pairs := []string{}
	for k, v := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k.Value, v.Inspect()))
	}

	return fmt.Sprintf("{%v}", strings.Join(pairs, ", "))
}

func (h *Hash) Equals(other Object) bool {
	switch other := other.(type) {
	case *Hash:
		if !h.Model.Equals(other.Model) {
			return false
		}

		for k, v := range h.Pairs {
			if !v.Equals(other.Get(k.Value)) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (h *Hash) Get(name string) Object {
	if meth, ok := h.Model.GetMethod(name); ok {
		meth.Hash = h
		return meth
	}

	for k, v := range h.Pairs {
		if k.Value == name {
			return v
		}
	}

	return &Null{}
}

func (h *Hash) Set(name string, val Object) {
	if _, exists := h.Model.GetMethod(name); exists {
		return
	}

	h.Pairs[String{Value: name}] = val
}

// Method Instance

type MethodInstance struct {
	Function *Object
	Hash     *Hash
}

func (mi *MethodInstance) Type() ObjectType { return METHOD_INSTANCE_OBJ }
func (mi *MethodInstance) Inspect() string  { return "<method instance>" }
func (mi *MethodInstance) Equals(other Object) bool {
	switch other := other.(type) {
	case *MethodInstance:
		return (*mi.Function).Equals(*other.Function) && mi.Hash.Equals(other.Hash)
	default:
		return false
	}
}
