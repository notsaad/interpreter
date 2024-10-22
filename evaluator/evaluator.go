package evaluator

import (
    "skibidi/ast"
    "skibidi/object"
    "fmt"
)

var (
    NULL  = &object.Null{}
    TRUE  = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

// we use object.Objects as a generic type which is then evaluated to the right type by the object.go file

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        
        // statements
    case *ast.Program:
        return evalProgram(node)
    case *ast.ExpressionStatement:
        return Eval(node.Expression) // recursive evaluation call
    case *ast.BlockStatement:
        return evalBlockStatement(node)

    // expressions
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.Boolean:
        return boolToBooleanObj(node.Value)
    case *ast.PrefixExpression:
        right := Eval(node.Right)
        return evalPrefixExpression(node.Operator, right)
    case *ast.InfixExpression:
        left := Eval(node.Left)
        right := Eval(node.Right)
        return evalInfixExpression(node.Operator, left, right)
    case *ast.IfExpression:
        return evalIfExpression(node)
    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue)
        return &object.ReturnValue{Value: val}

    }
    return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
    var result object.Object
    for _, statement := range stmts {
        result = Eval(statement)

        if returnValue, ok := result.(*object.ReturnValue); ok {
            return returnValue.Value
        }

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

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
    case TRUE:
        return FALSE
    case FALSE:
        return TRUE
    case NULL:
        return TRUE
    default:
        return FALSE
    }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
    // if the value to the right of the minus is not an integer then just return null
    if right.Type() != object.INTEGER_OBJ {
        return newError("unknown operator: -%s", right.Type())
    }
    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
    case "!":
        return evalBangOperatorExpression(right)
    case "-":
        return evalMinusPrefixOperatorExpression(right)
    default:
        return newError("unknown operator: %s%s", operator, right.Type())
    }
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right)
    case operator == "==":
        return boolToBooleanObj(left == right)
    case operator == "!=":
        return boolToBooleanObj(left != right)
    case left.Type() != right.Type():
        return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
    default:
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value


    switch operator {
    case "+":
        return &object.Integer{Value: leftVal + rightVal}
    case "-":
        return &object.Integer{Value: leftVal - rightVal}
    case "*":
        return &object.Integer{Value: leftVal * rightVal}
    case "/":
        return &object.Integer{Value: leftVal / rightVal}
    case "<":
        return boolToBooleanObj(leftVal < rightVal)
    case ">":
        return boolToBooleanObj(leftVal > rightVal)
    case "==":
        return boolToBooleanObj(leftVal == rightVal)
    case "!=":
        return boolToBooleanObj(leftVal != rightVal)
    default:
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }

}

func evalIfExpression(ie *ast.IfExpression) object.Object {
    condition := Eval(ie.Condition)
    if isTruthy(condition) {
        return Eval(ie.Consequence)
    } else if (ie.Alternative != nil) {
        return Eval(ie.Alternative)
    } else{
        return NULL
    }
}

func isTruthy(obj object.Object) bool {
    switch obj {
    case NULL:
        return false
    case TRUE:
        return true
    case FALSE:
        return false
    default:
        return true
    }
}

// a more specialized fn to evaluate block statements
func evalProgram(program *ast.Program) object.Object {
    var result object.Object
    
    for _, stmt := range program.Statements {
        result = Eval(stmt)

        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value
        case *object.Error:
            return result
        }

    }
    return result
}

// so the inner most blocked return statement is returned, and the rest of the block is not evaluated
func evalBlockStatement(block *ast.BlockStatement) object.Object {
    var result object.Object

    for _, stmt := range block.Statements {
        result = Eval(stmt)

        if result != nil{
            rt := result.Type()
            if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
                return result
            }
        }

    }
    return result
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

