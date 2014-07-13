package eval

import (
    "github.com/scampi/gosparqled/autocompletion"
    "net/http"
    "net/url"
    "log"
    "encoding/json"
    "io"
)

type param struct {
    key string
    value string
}

type Pof struct {
    POF string `json:"value"`
}

func MeasureWithoutDistinct(endpoint string, query string, template string) (uint, uint, float32) {
    pofs := getRecommendations(endpoint, query, template)
    counts := make(map[string]uint, len(pofs))
    for _,pof := range pofs {
        counts[pof.POF] += 1
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

func ExecuteQuery(endpoint string, params ...*param) io.ReadCloser {
    q := endpoint
    for i, p := range params {
        if i == 0 {
            q += "?"
        } else {
            q += "&"
        }
        q += p.key + "=" + url.QueryEscape(p.value)
    }
    resp, err := http.Get(q)
    if err != nil {
        log.Fatal(err)
    }
    return resp.Body
}

func getRecommendations(endpoint string, query string, template string) []Pof {
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
    body := ExecuteQuery(endpoint, &param{ key: "format", value: "application/json" }, &param{ key: "query", value: s.RecommendationQuery() })
    defer body.Close()
    // get the POF bindings
    dec := json.NewDecoder(body)
    var res = new(struct{Results struct{Bindings []map[string]Pof}})
    if err = dec.Decode(&res); err != nil {
        log.Fatal(err)
    }
    pofs := make([]Pof, len(res.Results.Bindings))
    for i,v := range res.Results.Bindings {
        pofs[i] = v["POF"]
    }
    return pofs
}

