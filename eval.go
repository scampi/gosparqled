package main

import (
    "os"
    "bufio"
    "log"
    "github.com/scampi/gosparqled/eval"
    "github.com/scampi/gosparqled/eval/data"
    "fmt"
)

func main() {
   tmpl := "SELECT ?POF ?count FROM <" + os.Args[5] + ">" +
    `    WHERE {
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
    tmpl += "LIMIT " + os.Args[2]
    queries := data.Load(os.Args[1])
    log.Printf("Processing file [%s]", os.Args[1])

    fi, err := os.Create(os.Args[4])
    if err != nil { log.Fatal(err) }
    defer fi.Close()
    w := bufio.NewWriter(fi)
    defer w.Flush()

    for _,query := range queries {
        log.Printf("\tProcessing query [%s]", query)
        for _,pof := range data.POFs(query) {
            log.Printf("\t\tProcessing [%s]", pof)
            measure := eval.Measure(os.Args[3], os.Args[6], pof, tmpl)
            w.WriteString(fmt.Sprintf("%v %v %v %v %v %v\n", measure.Min, measure.Max, measure.Avg, measure.Length, measure.ElapsedTime, measure.Recs))
        }
    }
}

