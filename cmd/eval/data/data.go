package main

import (
    "os"
    "github.com/scampi/gosparqled/eval/data"
    "fmt"
)

func main() {
    for _,query := range data.Load(os.Args[1]) {
        fmt.Printf("%s\n\n", query)
    }
}

