// Package data provides methods for handling query files.
package data

import (
    "os"
    "bufio"
    "github.com/golang/glog"
    "github.com/scampi/gosparqled/eval"
    "strings"
    "regexp"
)

// Load reads a file with SPARQL queries separated by a line with "###".
func Load(file string) []string {
    fi, err := os.Open(file)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    s := bufio.NewScanner(fi)
    var queries []string
    query := ""
    for s.Scan() {
        if s.Text() == "###" {
            queries = append(queries, query)
            query = ""
        } else {
            query += s.Text()
        }
    }
    return queries
}

// Clean loads the queries and executes them over the endpoint.
// Only those that have a solution are written to out.
func Clean(endpoint string, graph string, queries string, out string) {
    fi, err := os.Create(out)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    w := bufio.NewWriter(fi)
    defer w.Flush()
    for _,query := range Load(queries) {
        ask := strings.Replace(query, "SELECT *", "ASK FROM <" + graph + "> ", 1)
        res, _ := eval.Ask(endpoint, ask)
        if res {
            w.WriteString(query)
            w.WriteString("###\n")
        }
    }
}

// POFs returns a list of Recommendation queries.
// Each URI in the given query is transformed into a Point Of Focus.
func POFs(query string) []string {
    reg, _ := regexp.Compile("<[^>]*>")
    m := reg.FindAllStringIndex(query, -1)
    if m == nil {
        glog.Fatal("No match for " + query)
    }
    pofs := make([]string, len(m))
    for i,ind := range m {
        pofs[i] = query[:ind[0]] + " < " + query[ind[1]:]
    }
    return pofs
}

