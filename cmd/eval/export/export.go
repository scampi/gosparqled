package main

import (
    "github.com/golang/glog"
    "fmt"
    "os"
    "github.com/scampi/gosparqled/eval"
    "flag"
    "strings"
)

var output = flag.String("output", "", "The path to the output Latex file")
var gold = flag.String("gold", "", "The path to the Gold results for the evaluated queries")
var result = flag.String("result", "", "The path to the results directory. Can be multi-valued, with values separated by commas.")

func missingOption(option string) {
    fmt.Println("Missing option -" + option)
    flag.Usage()
    os.Exit(1)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if *output == "" { missingOption("output") }
    if *gold == "" { missingOption("gold") }
    if *result == "" { missingOption("result") }

    eval.Export(*output, *gold, strings.Split(*result, ","))
}
