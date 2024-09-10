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
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
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
    Token       token.Token // the 'return' toen
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
