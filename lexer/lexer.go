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

    l.skipWhitespace()

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
        tok.Type = token.EOF
    // default case in case the token is not recoginzed
    // used mainly to determine if the irregular input is a known keyword of the language
    // or if it is a custom identifier made by the user
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            return tok
            
        } else if isDigit(l.ch){
            // should read the entirety of the number and assign it
            tok.Literal = l.readDigits()
            tok.Type = token.INT
        } else {
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}



// found in a lot of parsers, sometimes is called eatWhitespace / consumeWhitespace
func (l *Lexer) skipWhitespace() {
    // skip characters as long as they are white space
    // this means whitespace does not matter in this langauge
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch){
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readDigits() string {
    position := l.position
    for isDigit(l.ch){
        l.readChar()
    }
    return l.input[position:l.position]
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return ('0' <= ch && ch <= '9')
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
