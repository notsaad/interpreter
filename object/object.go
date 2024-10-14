import (
    "fmt"
)

package object

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
)

// every value in the source code will be represented as an object for simplicity
// the objectType field in the interface is an struct, this is bc every value needs a diff internal representation (bools, chars, ints, etc)
type Object interface {
    Type() ObjectType // object type specific integer, char, etc
    Inspect() string
}

// integer object type, very simple
type Integer struct {
    Value int64
}

func (i *Integer) Inspect() string{
    return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
    return INTEGER_OBJ
}
