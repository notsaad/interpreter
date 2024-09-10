// the purpose of a parser is to translate user code into a data structure that represents that user input
// this style of parser is based upon Vaughn Pratt top down operator precedence parsing
// this style of parsing was generated by Pratt as an easier to understand alternative to grammar based parsing

package parser

import (
    "skibidi/ast"
    "skibidi/lexer"
    "skibidi/token"
    "fmt"
    "strconv"
)

// the precedences in the Skibidi programming language
const (
    // iota gives numbers to these values (think enum in c)
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

// seperate into prefix and infix operators because they are treated completely differently
// postfix operators are not supported in skibidi purely for simplicity sake
type (
    prefixParseFn   func() ast.Expression
    infixParseFn    func(ast.Expression) ast.Expression
)

type Parser struct {
    l           *lexer.Lexer // a pointer to an instance of the lexer (where we call nextToken())
    curToken    token.Token // these two act like two 'pointers' to the curr and upcoming tokens
    peekToken   token.Token
    errors []string

    prefixParseFns  map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)

    // read two tokens, so curToken and peekToken are both set
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier {
        Token: p.curToken,
        Value: p.curToken.Literal,
    }
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
    return p.errors
}

// works with the expectPeek function to make debugging easier in cases of errors
func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got: %s", t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.curTokenIs(token.EOF){
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
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement {
        Token: p.curToken,
    }

    stmt.Expression = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    // checks if we have a parsing function set up for the current token type
    prefix := p.prefixParseFns[p.curToken.Type]

    if prefix == nil {
        return nil
    }
    
    // if there is a parsing function associated with the current token type in the prefix position
    leftExp := prefix()
    
    return leftExp

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

    stmt := &ast.ReturnStatement{ 
        Token: p.curToken,
    }
    
    p.nextToken()

    // TODO: we're skipping the expression(s) until we see a semicolon

    for !p.curTokenIs(token.SEMICOLON){
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }

    // TODO: we're skipping the expressions until we encounter a semicolon

    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekToken.Type == t
}

// this is an assertion function that is common among all types of parsers
// primary purpose is to enforce the correctness of the order of tokens by checking the type of the next token
// it checks the type of the next token and if it is the expected type, only then does it advance the tokens via a nextToken call
func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    lit := &ast.IntegerLiteral{
        Token: p.curToken,
    }
    // first parameter is the string s
    // second parameter is the base value of the given value (0, 2-36)
    // third parameter is the bitSize (integer type that is returned) - 64 for int64, etc
    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

    if err != nil {
        msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }

    lit.Value = value

    return lit

}

