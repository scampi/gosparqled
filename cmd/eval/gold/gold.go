package main

import (
    "os"
    "bufio"
    "github.com/golang/glog"
    "github.com/scampi/gosparqled/eval"
    "github.com/scampi/gosparqled/eval/data"
    "fmt"
    "flag"
)

var queries = flag.String("queries", "", "The path to the queries file")
var endpoint = flag.String("endpoint", "", "The SPARQL endpoint")
var output = flag.String("output", "results", "The path to the output file")
var graph = flag.String("graph", "", "The named graph to get the count of each recommendation")

func missingOption(option string) {
    fmt.Println("Missing option -" + option)
    flag.Usage()
    os.Exit(1)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if *queries == "" { missingOption("queries") }
    if *endpoint == "" { missingOption("endpoint") }
    if *output == "" { missingOption("output") }
    if *graph == "" { missingOption("graph") }

    tmpl := "SELECT ?POF (COUNT(?POF) AS ?count) FROM <" + *graph + ">" +
        ` WHERE {
         {{range .Tps}}
             {{.S}} {{.P}} {{.O}} .
         {{end}}
         {{if .Keyword}}
             FILTER regex(?POF, "{{.Keyword}}", "i")
         {{end}}
             FILTER(?POF != <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>)
         }
         ORDER BY DESC(?count)
         LIMIT 10
    `
    file := data.Load(*queries)
    glog.Infof("Processing file [%s]", *queries)

    fi, err := os.Create(*output)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    w := bufio.NewWriter(fi)
    defer w.Flush()

    for _,query := range file {
        glog.Infof("\tProcessing query [%s]", query)
        for _,pof := range data.POFs(query) {
            glog.Infof("\t\tProcessing [%s]", pof)
            gold := eval.Gold(*endpoint, *graph, pof, tmpl)
            w.WriteString(fmt.Sprintf("%v\n", gold))
        }
    }
}

