package repl

import (
    "bufio"
    "fmt"
    "io"
    "skibidi/lexer"
    "skibidi/parser"
    "skibidi/evaluator"
)

const SKIBIDI_ASCII = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠄⠒⠒⠀⠀⠒⠂⠠⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⢄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠑⢄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠⠀⠀⣐⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⠀⠡⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡐⠀⠀⠈⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠁⠀⠀⠐⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢀⠄⠂⠠⢀⡀⠀⠀⠀⠀⠔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠨⡄⠀⠀⠀⠀⢀⠠⠐⠒⠐⠄
⢠⠁⠀⠀⠀⠀⠀⠈⠉⠉⠉⠀⠀⠀⠀⠀⣀⣤⣶⣿⣿⣿⣿⣿⣿⣶⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⡤⡀⠁⠉⠀⠀⠀⠀⠀⠈
⢸⠀⠀⠀⠀⠀⠀⠀⠀⠐⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⢱⠰⠀⠀⠀⠀⠀⠀⠀⢨
⠘⠄⠤⠀⠀⠀⠀⠀⠀⡆⠀⠀⢀⣼⣿⣿⠿⠿⠛⠻⠛⠛⠛⠙⠛⠛⠋⠩⠝⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⢘⣿⡶⣶⠲⢶⣴⣦⡄⠁
⠀⠀⠇⠀⠀⠀⠠⠔⢈⠁⠁⠀⣾⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢃⡹⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀⠀⣿⣝⡾⣑⢎⡷⣏⡇⠀
⠀⠀⠏⠛⠓⠻⠷⠿⢿⠀⠀⠀⣿⣿⠇⠀⠀⠀⠀⣀⣠⣄⣤⣄⣀⣀⣀⠀⠸⣽⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⣽⡾⣝⣳⢎⡿⡥⠂⠀
⠀⠀⢰⠀⠀⠀⠀⠀⢘⠀⠀⠀⢿⣿⢰⣤⣴⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⣿⡽⣯⡽⢾⣹⡓⠀⠀
⠀⠀⢸⠀⠀⠀⠀⠀⠨⠀⡀⢠⡘⣿⠸⣿⣿⡟⣋⢿⣿⡿⠙⣿⣿⡿⣽⢻⣿⡿⢽⣷⡇⠀⠀⠀⠀⠀⠀⢠⣿⣟⣷⡻⣝⣧⢹⠀⠀
⠀⠀⠘⡀⠀⠀⠀⠀⠀⡄⠙⠀⣿⣿⠀⠈⠛⠂⠡⠾⠟⠁⠀⡨⠟⠓⣖⣋⡁⣤⣿⢿⡇⠀⠀⠀⠀⠀⠀⢸⣿⣞⣷⡻⣽⡺⡥⠀⠀
⠀⠀⠀⡄⠀⠀⠀⠀⠀⡇⠀⠀⢹⣯⡦⣤⣦⣤⡄⣃⠀⠀⠀⢴⣛⣶⢏⣾⣿⣿⣿⣻⠁⠀⠀⠀⠀⠀⠀⣸⣿⢾⣳⣟⡷⡽⡁⠀⠀
⠀⠀⠀⠇⠀⠀⠀⠀⠀⢡⢸⣡⡍⢳⣿⣳⣭⡷⡤⠝⣤⣁⣀⣺⣿⢻⢾⣭⣿⣿⡿⠃⠀⠀⠀⠰⣭⣱⢀⣿⣿⣻⡽⣾⣝⣷⠁⠀⠀
⠀⠀⠀⢰⠀⠀⢀⢀⠀⢈⡆⠁⠀⠀⠸⣷⣻⢿⣿⣶⣼⣿⣿⣿⣯⣿⣿⣿⣿⣿⠇⠀⠀⠀⠀⠀⠉⠀⣼⣿⡷⣿⣽⣳⢯⠸⠀⠀⠀
⠀⠀⠀⠈⡄⠀⠀⢢⠳⡠⢼⠀⠀⠀⠀⢿⣶⣂⢿⣿⣏⢙⠈⠙⢻⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣽⡷⣯⣟⠯⡆⠀⠀⠀
⠀⠀⠀⠀⠰⠀⠀⠀⡑⢝⢦⣇⠀⠀⠀⠘⣿⣿⣦⡹⢿⣶⣶⣶⢿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⢠⣿⣿⣟⣾⣿⣻⣞⢷⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠇⠀⠶⡈⡎⣷⡾⣷⡀⠀⠀⣿⣿⣿⣿⣆⠈⣀⣈⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⡿⣷⣿⡎⡆⠀⠀⠀⠀
⠀⠀⠀⠀⢀⠈⠮⠵⠵⠎⠵⠛⠛⠛⠛⠻⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠷⠶⠶⠶⠾⠿⠿⠿⠿⠿⠿⠽⣿⣟⢗⠜⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣅⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠁⠀⠀⠀⠀⠀⠀⢀⠀⣀⣀⣀⣀⡈⠎⠁⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠈⠀⠀⠉⠉⠉⠀⠀⠀⠀⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠂⠐⠀⠈⠀⠉⠉⠉⠉⠀⠀⠀⠀⠐⢰⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⢄⡰⢠⢒⡰⢂⡖⣐⠺⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⠠⢄⢢⡱⣘⢦⡜⣣⢮⡵⣫⡼⣧⢹⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢱⡰⢶⡖⣦⢤⣤⣤⣤⣤⣄⣀⣀⣀⣀⣀⣀⣀⣄⣤⣔⣦⣳⢦⣟⣮⣷⣻⣽⣞⣿⡽⣯⣿⣟⣿⠫⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢇⠢⡙⢤⠛⣼⢳⣻⢮⣝⡯⣽⣹⢮⡽⣭⣛⢮⣗⡻⣞⡽⣯⣻⠷⣯⢷⣻⢾⣽⣛⣷⣻⣞⢧⠁⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠘⡵⣹⢦⣛⢦⣏⣷⣻⢮⡽⣶⣛⣮⢗⡷⣫⣟⡼⣝⣧⠿⣵⢯⡿⣭⢿⣭⣟⡾⣽⡞⣷⢏⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠸⣟⣯⣟⡿⣞⡷⣯⢿⡽⣞⣳⣭⢿⣹⣗⡾⣝⡾⣞⣻⡽⣾⡽⢯⣷⣛⡾⣽⣳⢿⢙⠊⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠪⠿⣽⣯⢿⡽⣯⢿⣽⣳⢯⣟⡷⡾⣝⣻⡼⣯⢷⣻⣗⣿⣻⢾⣽⣻⠷⡫⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠑⠩⢟⢿⣽⣟⣾⣽⣟⣾⣽⡿⣽⣷⣻⣽⣯⣷⣻⣾⡽⡛⠎⠃⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀⠉⠉⠉⠙⠛⠛⡙⢩⠩⣅⠫⢖⡭⣿⢺⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠀⠣⡐⢍⠎⡼⣻⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠄⠱⡈⢎⡱⣝⣏⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠠⢑⠢⣑⢮⢿⣠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠡⠌⡒⡡⢞⣯⢯⠱⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡐⢢⠁⢧⡙⣮⣟⣾⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢂⠍⢦⡙⣮⢿⡈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡌⡘⢦⣙⣾⣋⠁⠀  ,⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

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
        p := parser.New(l)

        program := p.ParseProgram()
        if len(p.Errors()) != 0 {
            printParserErrors(out, p.Errors())
            continue
        }

        evaluated := evaluator.Eval(program)

        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
    }
}

func printParserErrors(out io.Writer, errors []string){
    io.WriteString(out, SKIBIDI_ASCII)
    io.WriteString(out, "We ran into an error here!\n")
    io.WriteString(out, "Parser Errors:\n")
    for _, msg := range errors {
        io.WriteString(out, "\t" + msg + "\n")
    }
}
