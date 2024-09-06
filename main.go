package main

import (
    "fmt"
    "os"
    "os/user"
    "skibidi/repl"
)

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! This is the skibidi programming language!\n", user.Username)
    fmt.Printf("Feel free to type in commands:\n")
    repl.Start(os.Stdin, os.Stdout)
}
