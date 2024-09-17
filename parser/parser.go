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
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.TRUE, p.parseBoolean)
    p.registerPrefix(token.FALSE, p.parseBoolean)
    p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
    p.registerPrefix(token.IF, p.parseIfExpression)
    p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.LPAREN, p.parseCallExpression)

    // read two tokens, so curToken and peekToken are both set
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression {
        Token:      p.curToken,
        Operator:   p.curToken.Literal,
    }

    p.nextToken() // we advance the token we're looking at because of the nature of prefixes
    // only need to verify the type of prefix, everything else is concerned with the next token (expression.Right)

    expression.Right = p.parseExpression(PREFIX)
    return expression

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

func (p *Parser) parseBoolean() ast.Expression {
    return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function for %s found", t)
    p.errors = append(p.errors, msg)
}

// this function is the heart of the Pratt parser
func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.noPrefixParseFnError(p.curToken.Type)
        return nil
    }
    leftExp := prefix()

    // as long as the next token is not a semicolon and the next tokens precedence is not greater
    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        // by repeating iteration based on precedence, even complex expression with brackets and other items are evaluated correctly
        // this is because it keeps going until it finds the 'middle' of the expression and evaluates outwards from there
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }

        p.nextToken()

        leftExp = infix(leftExp)

    }    

    return leftExp
}

// look up table for precedences that references the previously defined const
var precedences = map[token.TokenType] int{
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

func (p *Parser) peekPrecedence() int {
    // look at the precedence of the next token
    if p, ok := precedences[p.peekToken.Type]; ok {
        return p
    // return lowest if not found in lookup table
    }
    return LOWEST
}

func (p *Parser) curPrecedence() int {
    // precedence of cur token
    if p, ok := precedences[p.curToken.Type]; ok {
        return p
    }
    // return lowest if not found in lookup table
    return LOWEST
}

// this parsing function takes the left expression as an argument
// this is so it can use this argument to construct the ast node
// (the infix expression ast node requires a left and right node)
func (p *Parser) parseInfixExpression (left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression {
        Token:      p.curToken,
        Operator:   p.curToken.Literal,
        Left:       left,
    }

    // assign the current tokens precedence to a variable
    precedence := p.curPrecedence()
    p.nextToken()
    // fill the right expresison with that precedence so the evaluation occurs properly
    expression.Right = p.parseExpression(precedence)

    return expression

}

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.nextToken()
    
    exp := p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) {
        return nil
    }
    return exp
}

func (p *Parser) parseIfExpression() ast.Expression {

    expression := &ast.IfExpression{
        Token: p.curToken,
    }

    // the correct syntax for an if statement is "if (" and then something else
    // if the if statement is not followed by an opening parentheses, then abort and throw an error
    if !p.expectPeek(token.LPAREN){
        return nil
    }

    p.nextToken()
    // parses the actual contents of the if statement
    expression.Condition = p.parseExpression(LOWEST)

    // if the if statement is not closed by a bracket, throw an error
    if !p.expectPeek(token.RPAREN){
        return nil
    }

    p.nextToken()

    // { to start the logic after the if statement
    if !p.expectPeek(token.LBRACE) {
        return nil
    }

    expression.Consequence = p.parseBlockStatement()

    if p.peekTokenIs(token.ELSE) {
        p.nextToken()

        if !p.expectPeek(token.LBRACE) {
            return nil
        }

        expression.Alternative = p.parseBlockStatement()
    }

    return expression

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement {
        Token: p.curToken,
    }

    p.nextToken()

    for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        if stmt != nil {
            block.Statements = append(block.Statements, stmt)
        }
        p.nextToken()
    }

    return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
    lit := &ast.FunctionLiteral{
        Token: p.curToken,
    }

    if !p.expectPeek(token.LPAREN) {
        return nil
    }

    lit.Parameters = p.parseFunctionParameters()

    if !p.expectPeek(token.LBRACE){
        return nil
    }

    // the .Body contains the actual content of the function (the stuff inside the curly braces) which is why its parsed like any other block statement
    lit.Body = p.parseBlockStatement()

    return lit

}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    // this is the case where there is no function parameters
    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return identifiers
    }

    p.nextToken()

    ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    identifiers = append(identifiers, ident)

    // as long as there is more parameters to parse, continue
    for p.peekTokenIs(token.COMMA){
        p.nextToken()
        p.nextToken()
        ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        identifiers = append(identifiers, ident)
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return identifiers

}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    exp := &ast.CallExpression {
        Token: p.curToken,
        Function: function,
    }

    exp.Arguments = p.parseCallArguments()

    return exp

}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return args
    }

    p.nextToken()

    args = append(args, p.parseExpression(LOWEST))

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()
        args = append(args, p.parseExpression(LOWEST))
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return args

}

