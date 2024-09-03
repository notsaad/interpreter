// the lexer works to read in source code, and output tokens that represent that source code
package lexer

import (
    "testing"
    "skibidi/token"
)

// in a real, production quality interpreter we would want to attach filename and line number to the token so that error messages are more descriptive
func TestNextToken(t *testing.T) {
    // walrus operator is shorthand variable declaration that only works for vars inside functions
    // it infers the variable type from the assigned value
    // eg: 'x := 10' == 'var x int = 10'
    
    // back ticks are for raw string literals
    input := `=+(){},;`

    tests := []struct {
        expectedType    token.TokenType
        expectedLiteral string
    }{
        {token.ASSIGN, "="}
        {token.PLUS, "+"}
        {token.LPAREN, "("}
        {token.RPAREN, ")"}
        {token.LBRACE, "{"}
        {token.RBRACE, "}"}
        {token.COMMA, ","}
        {token.SEMICOLON, ";"}
        {token.EOF, ""}
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.NextToken()

        if tok.Type != tt.expected {
            // like throwing an exception (Fatalf)
            t.Fatalf("tests[%d] - token type wrong. expected: %q, got: %q", i, tt.expectedType, tok.Type)
        }

        if tok.Literal != tt.expectedLiteral {
            // %q is a format verb (like in c), specifically for strings in this case
            t.Fatalf("tests[%d] - token literal wrong. expected: %q, got: %q",i, tt.expectedLiteral, tok.Literal)
        }

    }

}

