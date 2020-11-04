package main

import (
	"fmt"
	"os"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("Hey!\n")
        return
    }
    fmt.Printf("Hey %s!\n", os.Args[1])

    fmt.Println("\nInspect")
    fmt.Printf("Args: [ %s ]", strings.Join(os.Args, ", "))
}
