// this lexer only supports ASCII characters
// the purpose of a lexer is to convert user written skibidi code into tokens
package lexer

import "skibidi/token"

type Lexer struct {
    input           string
    position        int // current pos in input (points to current char)
    readPosition    int // current reading position in input (points to after current char)
    ch              byte // current char under examination
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    switch l.ch{
    case '=':
        tok = newToken(token.ASSIGN, l.ch)
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '0':
        tok.Literal = ""
        tok = newToken(token.EOF, l.ch)
    }

    l.readChar()
    return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

// just moves the curr and next character along
// essentially a method for a Lexer, which is the receiver type
func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input){
        // 0 is the ascii code for the "NUL" character, signifying EOF
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition++
}
