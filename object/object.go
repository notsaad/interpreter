package object

import (
    "fmt"
    "skibidi/ast"
    "bytes"
    "strings"
)

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ    = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
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

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
    return rv.Value.Inspect()
}

type Error struct {
    Message string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJ
}

func (e *Error) Inspect() string {
    return "ERROR: " + e.Message
}

// self explanatory, the parts that make up a functions structure (at least the parts we care about)
type Function struct {
    Parameters  []*ast.Identifier
    Body        *ast.BlockStatement
    Env         *Environment
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
    var out bytes.Buffer // basically just a string variable

    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")

    return out.String()

}

