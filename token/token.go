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

var keywords = map[string]TokenType{
    "fn": FUNCTION,
    "let": LET,
}

// checks the keywords table to see if the identifier is a known keyword (like var)
// if it isn't then it returns token.IDENT, signifying that it is a custom user identifier
func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return IDENT
}


