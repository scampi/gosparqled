package main

import (
    "os"
    "github.com/scampi/gosparqled/eval"
)

func main() {
    eval.Export(os.Args[1], os.Args[2], os.Args[3:])
}
