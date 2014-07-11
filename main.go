package main

import (
    "github.com/gopherjs/gopherjs/js"
    "github.com/scampi/gosparqled/autocompletion"
    "text/template"
)

var tmpl string = `
SELECT DISTINCT ?POF
WHERE {
{{range .Tps}}
    {{.S}} {{.P}} {{.O}} .
{{end}}
{{if .Keyword}}
    FILTER regex(?POF, "{{.Keyword}}", "i")
{{end}}
}
LIMIT 10
`
var tp, _ = template.New("rec").Parse(tmpl)

func RecommendationQuery(query string, callback func(string)) {
    go func(query string) {
        s := &autocompletion.Sparql{ Buffer : query, Bgp : &autocompletion.Bgp{Template : tp} }
        s.Init()
        err := s.Parse()
        if err == nil {
            s.Execute()
            callback(s.RecommendationQuery())
        } else {
            callback(query + "\n" + err.Error())
        }
    }(query)
}

func main() {
    js.Global.Set("autocompletion", map[string]interface{}{
        "RecommendationQuery": RecommendationQuery,
    })
}
