package ast

// ast -> abstract syntax tree (abstract because some syntax like semicolons or brackets may be left out)

import (
    "skibidi/token"
    "bytes"
)

type Node interface {
    TokenLiteral() string
    String() string // main purpose of this is to be able to print nodes for debugging and comparing node values
}

func (p *Program) String() string {

    var out bytes.Buffer // creates a buffer

    for _, s:= range p.Statements {
        out.WriteString(s.String()) // writes the return value of each statements String() method to it
    }
    return out.String() // returns the buffer as a string
}

type Statement interface {
    // interfaces the default node interface into 
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type PrefixExpression struct {
    Token       token.Token // the token value of the prefix expression (eg ! or -)
    Operator    string
    Right       Expression // expression to the right (lol)
}

type InfixExpression struct {
    Token       token.Token
    Left        Expression // expression to the left (lol)
    Operator    string
    Right       Expression // expression to the right (lol)
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

// program node is going to be the root node of every AST our parser makes
type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

func (il *IntegerLiteral) expressionNode()  {}

func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

type ExpressionStatement struct {
    Token       token.Token // the first token of the expression
    Expression  Expression
}

type LetStatement struct {
    Token   token.Token // the token.LET token
    Name    *Identifier
    Value   Expression
}

func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")

    return out.String()
}

type ReturnStatement struct {
    Token       token.Token // the 'return' token
    ReturnValue Expression // contains the expression that is to be returned
}

func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    out.WriteString(";")

    return out.String()
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }
    return ""
}

// holds the identifier of the binding (x in 'let x = 5;')
type Identifier struct {
    Token   token.Token // the token.IDENT token
    Value   string
}

// expression nodes would be the ones that hold stuff like the '5' in 'let x = 5'
// expression nodes could also be identifiers in situations like 'let x = 5'
func (i *Identifier) expressionNode()           {}
func (i *Identifier) TokenLiteral() string      { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

func (pe *PrefixExpression) expressionNode() {

}

func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    // adding the parentheses is so we can easily tell what actually belongs to the prefix expressionNode (in the case of debugging or smt)
    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()

}

func (ie *InfixExpression) expressionNode() {
    
}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    // adding the parentheses is so we can easily tell what actually belongs to the infix expressionNode (in the case of debugging or smt)
    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()

}

// very easy implementation for a boolean value node for the ast
type Boolean struct {
    Token   token.Token
    Value   bool
}

func (b *Boolean) expressionNode() {

}

func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}

func (b *Boolean) String() string {
    return b.Token.Literal
}

type IfExpression struct {
    Token       token.Token // the 'if' token
    Condition   Expression // holds the value of the if statement so to speak
    Consequence *BlockStatement // if the condition is true
    Alternative *BlockStatement // else
}

func (fe *IfExpression) expressionNode() {

}

func (fe *IfExpression) TokenLiteral() string {
    return fe.Token.Literal
}

func (fe *IfExpression) String() string {

    var out bytes.Buffer

    out.WriteString("if")
    out.WriteString(fe.Condition.String())
    out.WriteString(" ")
    out.WriteString(fe.Consequence.String())

    if fe.Alternative != nil {
        out.WriteString("else ")
        out.WriteString(fe.Alternative.String())
    }

    return out.String()
}

type BlockStatement struct {
    Token       token.Token // the { token
    Statements  []Statement
}

func (bs *BlockStatement) statementNode() {

}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

