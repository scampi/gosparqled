package eval

import (
    "github.com/scampi/gosparqled/autocompletion"
    "log"
)

func Measure(endpoint string, query string, template string) (uint, uint, float32) {
    pofs := getRecommendations(endpoint, query, template)
    counts := make(map[string]uint, len(pofs))
    for _,pof := range pofs {
        counts[pof] += 1
    }
    min, max, sum := ^uint(0), uint(0), float32(0)
    for _,c := range counts {
        if c < min {
            min = c
        }
        if c > max {
            max = c
        }
        sum += float32(c)
    }
    return min, max, sum / float32(len(counts))
}

func getRecommendations(endpoint string, query string, template string) []string {
    // retrieve the recommendations
    var scope *autocompletion.Scope
    if len(template) == 0 {
        scope = autocompletion.NewScope()
    } else {
        scope = autocompletion.NewScopeWithTemplate(template)
    }
    s := &autocompletion.Sparql{ Buffer : query, Scope : scope }
    s.Init()
    err := s.Parse()
    if err != nil {
        log.Fatal(err)
    }
    s.Execute()
    body := ExecuteQuery(endpoint, s.RecommendationQuery())
    defer body.Close()
    // get the POF bindings
    bindings := GetBindings(body)
    pofs := make([]string, len(bindings))
    for i,v := range bindings {
        pofs[i] = v["POF"]["value"]
    }
    return pofs
}

