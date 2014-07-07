package main

import (
    "github.com/gopherjs/gopherjs/js"
    "github.com/hoisie/mustache"
    "github.com/scampi/gosparqled/autocompletion"
)

func RecommendationQuery(query string) string {
    tmpl := `
SELECT DISTINCT ?POF
WHERE {
{{#Tps}}
    {{{S}}} {{{P}}} {{{O}}} .
{{/Tps}}
}
LIMIT 100
    `
    tp, _ := mustache.ParseString(tmpl)
    s := &autocompletion.Sparql{ Buffer : query, Bgp : autocompletion.Bgp{Template : tp} }
    s.Init()
    if err := s.Parse(); err == nil {
        s.Execute()
        return s.RecommendationQuery()
    }
    return ""
}

func main() {
    js.Global.Set("autocompletion", map[string]interface{}{
        "RecommendationQuery": RecommendationQuery,
    })
}
