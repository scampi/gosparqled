package main

import (
    "os"
    "bufio"
    "github.com/golang/glog"
    "github.com/scampi/gosparqled/eval"
    "github.com/scampi/gosparqled/eval/data"
    "fmt"
    "flag"
    "strconv"
)

var queries = flag.String("queries", "", "The path to the queries file")
var limit = flag.Int("limit", 10, "The value of the LIMIT clause")
var endpoint = flag.String("endpoint", "", "The SPARQL endpoint")
var output = flag.String("output", "results", "The path to the output file")
var recGraph = flag.String("rec-graph", "", "The named graph to get recommendations from")
var countGraph = flag.String("count-graph", "", "The named graph to get the count of each recommendation")

func missingOption(option string) {
    fmt.Println("Missing option -" + option)
    flag.Usage()
    os.Exit(1)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if *queries == "" { missingOption("queries") }
    if *endpoint == "" { missingOption("endpoing") }
    if *output == "" { missingOption("output") }
    if *recGraph == "" { missingOption("rec-graph") }
    if *countGraph == "" { missingOption("count-graph") }

    tmpl := "SELECT ?POF ?count FROM <" + *recGraph + ">" +
        ` WHERE {
         {{range .Tps}}
             {{.S}} {{.P}} {{.O}} .
         {{end}}
         {{if .Keyword}}
             FILTER regex(?POF, "{{.Keyword}}", "i")
         {{end}}
             BIND(1 as ?count)
             FILTER(?POF != <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>)
         }
    `
    tmpl += "LIMIT " + strconv.Itoa(*limit)
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
            measure := eval.Measure(*endpoint, *countGraph, pof, tmpl)
            w.WriteString(fmt.Sprintf("%v %v %v %v %v %v\n", measure.Min, measure.Max, measure.Avg, measure.Length, measure.ElapsedTime, measure.Recs))
        }
    }
}

