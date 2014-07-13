package data

import (
    "os"
    "bufio"
    "log"
    "github.com/scampi/gosparqled/eval"
    "strings"
)

func Load(file string) []string {
    fi, err := os.Open(file)
    if err != nil { log.Fatal(err) }
    defer fi.Close()
    s := bufio.NewScanner(fi)
    // function split on lines equal to "###"
    s.Split(func(data []byte, atEOF bool) (int, []byte, error) {
        advance := 0
        tmpData := data
        for {
            line, token, err := bufio.ScanLines(tmpData, atEOF)
            if err != nil {
                log.Fatal(err)
            }
            if !atEOF && line == 0 { // get more input
                return 0, nil, nil
            }
            advance += line
            tmpData = tmpData[line:]
            if string(token) == "###" { // delimiter
                data = data[0:advance-line]
                break
            } else if len(tmpData) == 0 { // scanned whole array
                data = data[0:advance]
                break
            }
        }
        return advance, data, nil
    })
    var queries []string
    for s.Scan() {
        queries = append(queries, s.Text())
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
        body := eval.ExecuteQuery(endpoint, ask)
        defer body.Close()
        if (eval.Ask(body)) {
            w.WriteString(query)
            w.WriteString("###\n")
        }
    }
}

