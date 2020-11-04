package main

import (
	"fmt"
	"os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("Hey!\n")
        return
    }
    fmt.Printf("Hey %s!\n", os.Args[1])
    for _, arg := range os.Args {
        println(arg)
    }
}
