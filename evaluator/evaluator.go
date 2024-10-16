package evaluator

import (
    "skibidi/ast"
    "skibidi/object"
)

var (
    NULL  = &object.Null{}
    TRUE  = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        
        // statements
    case *ast.Program:
        return evalStatements(node.Statements)
    case *ast.ExpressionStatement:
        return Eval(node.Expression) // recursive evaluation call

    // expressions
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.Boolean:
        return boolToBooleanObj(node.Value)

    }
    return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
    var result object.Object
    for _, statement := range stmts {
        result = Eval(statement)
    }
    return result
}

func boolToBooleanObj(input bool) *object.Boolean {
    if input {
        return TRUE
    } else {
        return FALSE
    }
}
