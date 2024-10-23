package object

/*
the environment is used to keep track of values by associating them with a name
in this example, the environment is just a wrapper of a standard Go hashmap
*/

// now refers to an outer environment as well which allows an environment to be enclosed, same sort of principle as variable scopes
// meaning if something is not found in the current environment, the outer environment is checked (the inner scope extends the outer one)
func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}

func NewEnvironment() *Environment {
    s := make(map[string]Object)
    return &Environment{store: s, outer: nil}
}

type Environment struct {
    store map[string]Object
    outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    if !ok && e.outer != nil {
        obj, ok = e.outer.Get(name)
    }
    return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
    e.store[name] = val
    return val
}

