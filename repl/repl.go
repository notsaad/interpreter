package repl

import (
    "bufio"
    "fmt"
    "io"
    "skibidi/lexer"
    "skibidi/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }
        // read in input line by line
        line := scanner.Text()
        // send each line to the lexer to get broken down to Tokens
        l := lexer.New(line)
        // print each of the received tokens one by one
        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
            fmt.Printf("%+v\n", tok)
        }
    }
}
