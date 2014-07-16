package data

import (
    "os"
    "bufio"
    "log"
    "github.com/scampi/gosparqled/eval"
    "strings"
    "regexp"
)

func Load(file string) []string {
    fi, err := os.Open(file)
    if err != nil { log.Fatal(err) }
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

func Clean(endpoint string, graph string, queries []string, out string) {
    fi, err := os.Create(out)
    if err != nil { log.Fatal(err) }
    defer fi.Close()
    w := bufio.NewWriter(fi)
    defer w.Flush()
    for _,query := range queries {
        ask := strings.Replace(query, "SELECT *", "ASK FROM <" + graph + "> ", 1)
        body, _ := eval.ExecuteQuery(endpoint, ask)
        defer body.Close()
        if (eval.Ask(body)) {
            w.WriteString(query)
            w.WriteString("###\n")
        }
    }
}

func POFs(query string) []string {
    reg, _ := regexp.Compile("<[^>]*>")
    m := reg.FindAllStringIndex(query, -1)
    if m == nil {
        log.Fatal("No match for " + query)
    }
    pofs := make([]string, len(m))
    for i,ind := range m {
        pofs[i] = query[:ind[0]] + " < " + query[ind[1]:]
    }
    return pofs
}

