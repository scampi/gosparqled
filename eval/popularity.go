package eval

import (
    "github.com/scampi/gosparqled/autocompletion"
    "log"
    "strconv"
    "time"
    "sort"
    "math"
)

type Measurement struct {
    Min, Max, Length int
    Avg float32
    ElapsedTime time.Duration
    Recs []Recommendation
}

type Recommendation struct {
    item string
    count int
}

type ByCount []Recommendation

func (bc ByCount) Len() int { return len(bc) }
func (bc ByCount) Swap(i, j int)      { bc[i], bc[j] = bc[j], bc[i] }
func (bc ByCount) Less(i, j int) bool { return bc[i].count < bc[j].count }

func Measure(endpoint string, from string, query string, template string) Measurement {
    pofs, elapsedTime := getRecommendations(endpoint, from, query, template)
    min, max, sum := math.MaxInt32, 0, float32(0)
    for _,c := range pofs {
        if c.count < min {
            min = c.count
        }
        if c.count > max {
            max = c.count
        }
        sum += float32(c.count)
    }
    return Measurement{ Min: min, Max: max, Avg: sum / float32(len(pofs)), ElapsedTime: elapsedTime, Length: len(pofs), Recs: pofs }
}

func getRecommendations(endpoint string, from string, query string, template string) ([]Recommendation, time.Duration) {
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
    // get the POF bindings and rank them
    bindings := GetBindings(body)
    counts := make(map[string]int, len(bindings))
    for _,v := range bindings {
        count,_ := strconv.Atoi(v["count"]["value"])
        counts[v["POF"]["value"]] += count
    }
    log.Printf("Results: %v\n", counts)
    var pofs ByCount
    for k,v := range counts {
        pofs = append(pofs, Recommendation{ item: k, count: v })
    }
    sort.Sort(sort.Reverse(pofs))
    min := int(math.Min(10, float64(len(pofs))))
    top := pofs[:min]
    log.Printf("TOP10: %v\n", top)
    // get the popularity of each recommended item
    values := "values ?POF { "
    for _,r := range top {
        values += "<" + r.item + "> "
    }
    values += "}"
    tmpl := "SELECT ?POF (count(?POF) as ?count) FROM <" + from + "> WHERE {" +
         values + `
         {{range .Tps}}
             {{.S}} {{.P}} {{.O}} .
         {{end}}
         }
        `
    scope = autocompletion.NewScopeWithTemplate(tmpl)
    s = &autocompletion.Sparql{ Buffer : query, Scope : scope }
    s.Init()
    err = s.Parse()
    if err != nil {
        log.Fatal(err)
    }
    s.Execute()
    body2, _ := ExecuteQuery(endpoint, s.RecommendationQuery())
    defer body2.Close()
    var popularity []Recommendation
    bindings = GetBindings(body2)
    for _,v := range bindings {
        count,_ := strconv.Atoi(v["count"]["value"])
        popularity = append(popularity, Recommendation{ item: v["POF"]["value"], count: count })
    }
    log.Printf("Popularity=%v\n", bindings)
    return popularity, elapsedTime
}

