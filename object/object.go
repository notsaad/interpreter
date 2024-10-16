package object

import (
    "fmt"
)

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ    = "NULL"
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

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}

func (n *Null) Inspect() string {
    return "null"
}
