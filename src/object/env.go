package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	env := &Environment{store: s, outer: nil}
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	for k, v := range DefaultModels {
		if k == name {
			return v, true
		}
	}

	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Declare(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Assign(name string, val Object) Object {
	if e.outer != nil {
		if _, ok := e.outer.store[name]; ok {
			return e.outer.Assign(name, val)
		}
	}

	return e.Declare(name, val)
}
