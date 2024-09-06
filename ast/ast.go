package ast

// ast -> abstract syntax tree (abstract because some syntax like semicolons or brackets may be left out)

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
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

type LetStatement struct {
    Token   token.Token // the token.LET token
    Name    *Identifier
    Value   Expression
}

func (ls *LetStatement) statementNode()         {}
func (ls *LetStatement) TokenLiteral() string   { return ls.Token.Literal }

// holds the identifier of the binding (x in 'let x = 5;')
type Identifier struct {
    Token   token.Token // the token.IDENT token
    Value   string
}

// expression nodes would be the ones that hold stuff like the '5' in 'let x = 5'
// expression nodes could also be identifiers in situations like 'let x = 5'
func (i *Identifier) expressionNode()           {}
func (i *Identifier) TokenLiteral() string      { return i.Token.Literal }
