package token

const (
    ILLEGAL = "ILLEGAL"
    EOF     = ""

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

type TokenType string

type Token struct {
    // used to distinguish between things like integer or brackets, etc
    Type    TokenType
    // holds the literal value of the token (ex 5, 10, etc)
    Literal string
}


