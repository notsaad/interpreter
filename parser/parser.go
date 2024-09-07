// the purpose of a parser is to translate user code into a data structure that represents that user input
package parser

import (
    "skibidi/ast"
    "skibidi/lexer"
    "skibidi/token"
)

type Parser struct {
    l           *lexer.Lexer // a pointer to an instance of the lexer (where we call nextToken())
    curToken    token.Token // these two act like two 'pointers' to the curr and upcoming tokens
    peekToken   token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}

    // read two tokens, so curToken and peekToken are both set
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement
