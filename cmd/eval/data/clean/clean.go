package main

import (
    "flag"
    "log"
    "os"
    "bufio"
    "errors"
    "bytes"
    "regexp"
    "sort"
    "strings"
)

var queries = flag.String("queries", "", "The path to the queries file")

func main() {
    flag.Parse()

    if *queries == "" {
        flag.Usage()
        log.Fatal("Missing queries option")
    }
    fi, err := os.Open(*queries)
    if err != nil { log.Fatal(err) }
    r := bufio.NewScanner(fi)
    split := func(data []byte, atEOF bool) (int, []byte, error) {
        ind := bytes.Index(data, []byte("###\n"))
        if ind == -1 {
            if atEOF {
                return 0, nil, errors.New("Invalid input")
            }
            return 0, nil, nil
        }
        return ind + 4, data[:ind-1], nil
    }

    reg := regexp.MustCompile("<[^>]*>")
    r.Split(split)
    duplicates := make(map[string]string)
    for r.Scan() {
        uris := reg.FindAllString(r.Text(), -1)
        sort.Strings(uris)
        key := strings.Join(uris, " ")
        if v, ok := duplicates[key]; ok {
            log.Printf("Duplicate: %v", v)
        } else {
            duplicates[key] = r.Text()
        }
    }
    if r.Err() != nil {
        log.Fatal(r.Err())
    }
}

