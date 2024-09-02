package token

type TokenType string

type Token struct {
    // used to distinguish between things like integer or brackets, etc
    type    TokenType
    // holds the literal value of the token (ex 5, 10, etc)
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"

    // identifiers and literals
    IDENT   = "IDENT"
    INT     = "INT"
    
    // operators
    ASSIGN  = "="
    PLUS    = "+"

    // delimiters
    COMMA   = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"

    // keywords
    FUNCTION = "FUNCTION"
    LET = "LET"

)

