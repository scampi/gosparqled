package eval

import (
    "fmt"
    "testing"
)

func TestTopK(t *testing.T) {
    tmpl := `
        SELECT ?POF
        WHERE {
        {{range .Tps}}
            {{.S}} {{.P}} {{.O}} .
        {{end}}
        {{if .Keyword}}
            FILTER regex(?POF, "{{.Keyword}}", "i")
        {{end}}
        }
        LIMIT 100
    `
    fmt.Println(Measure("http://aemet.linkeddata.es/sparql", "select * { ?s < }", tmpl))
}
