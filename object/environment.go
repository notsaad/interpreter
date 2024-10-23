package object

/*
the environment is used to keep track of values by associating them with a name
in this example, the environment is just a wrapper of a standard Go hashmap
*/

func NewEnvironment() *Environment {
    s := make(map[string]Object)
    return &Environment{store: s}
}

type Environment struct {
    store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
    e.store[name] = val
    return val
}

