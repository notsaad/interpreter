// the lexer works to read in source code, and output tokens that represent that source code
package lexer

import (
    "testing"
    "skibidi/token"
    "fmt"
)
 
// in a real, production quality interpreter we would want to attach filename and line number to the token so that error messages are more descriptive
func TestNextToken(t *testing.T) {
    // walrus operator is shorthand variable declaration that only works for vars inside functions
    // it infers the variable type from the assigned value
    // eg: 'x := 10' == 'var x int = 10'

    // back ticks are for raw string literals
    input := `let five = 5;
    let ten = 10;
    let add = fn(x, y) {
    x + y;
    };
    let result = add(five, ten);
    `

    tests := []struct {
        expectedType    token.TokenType
        expectedLiteral string
    }{
        {token.LET, "let"},
        {token.IDENT, "five"},
        {token.ASSIGN, "="},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        {token.LET, "let"},
        {token.IDENT, "ten"},
        {token.ASSIGN, "="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        {token.LET, "let"},
        {token.IDENT, "add"},
        {token.ASSIGN, "="},
        {token.FUNCTION, "fn"},
        {token.LPAREN, "("},
        {token.IDENT, "x"},
        {token.COMMA, ","},
        {token.IDENT, "y"},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.IDENT, "x"},
        {token.PLUS, "+"},
        {token.IDENT, "y"},
        {token.SEMICOLON, ";"},
        {token.RBRACE, "}"},
        {token.SEMICOLON, ";"},
        {token.LET, "let"},
        {token.IDENT, "result"},
        {token.ASSIGN, "="},
        {token.IDENT, "add"},
        {token.LPAREN, "("},
        {token.IDENT, "five"},
        {token.COMMA, ","},
        {token.IDENT, "ten"},
        {token.RPAREN, ")"},
        {token.SEMICOLON, ";"},
        {token.EOF, ""},
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.NextToken()
        fmt.Print(tok.Literal)
        if tok.Type != tt.expectedType {
            // like throwing an exception (Fatalf)
            t.Fatalf("tests[%d] - token type wrong. expected: %q, got: %q", i, tt.expectedType, tok.Type)
        }

        if tok.Literal != tt.expectedLiteral {
            // %q is a format verb (like in c), specifically for strings in this case
            t.Fatalf("tests[%d] - token literal wrong. expected: %q, got: %q",i, tt.expectedLiteral, tok.Literal)
        }

    }

}

