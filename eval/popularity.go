package eval

import (
    "github.com/scampi/gosparqled/autocompletion"
    "log"
    "strconv"
    "time"
)

type Recommendation struct {
    item string
    count uint
}

func Measure(endpoint string, query string, template string) (uint, uint, float32, time.Duration) {
    pofs, elapsedTime := getRecommendations(endpoint, query, template)
    min, max, sum := ^uint(0), uint(0), float32(0)
    for _,rec := range pofs {
        if rec.count < min {
            min = rec.count
        }
        if rec.count > max {
            max = rec.count
        }
        sum += float32(rec.count)
    }
    return min, max, sum / float32(len(pofs)), elapsedTime
}

func getRecommendations(endpoint string, query string, template string) ([]Recommendation, time.Duration) {
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
    body, elapsedTime := ExecuteQuery(endpoint, s.RecommendationQuery())
    defer body.Close()
    // get the POF bindings
    bindings := GetBindings(body)
    pofs := make([]Recommendation, len(bindings))
    for i,v := range bindings {
        pofs[i].item = v["POF"]["value"]
        count, _ := strconv.Atoi(v["count"]["value"])
        pofs[i].count = uint(count)
    }
    return pofs, elapsedTime
}

