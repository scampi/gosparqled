// Package eval provides methods for evaluating the SPARQL recommendations.
package eval

import (
    "github.com/scampi/gosparqled/autocompletion"
    "log"
    "strconv"
    "time"
    "sort"
    "math"
)

// Measurement contains various measures about recommendations for a SPARQL query
type Measurement struct {
    // Min is the lowest number of occurences of a recommendation
    // Max is the highest number of occurrences of a recommendation
    // Length is the total number of recommendations
    Min, Max, Length int
    // Avg is the average number of occurrences in the list of recommendations
    Avg float32
    // The time spent on retrieving the list of recommendations
    ElapsedTime time.Duration
    // The ranked list of recommendations
    Recs []Recommendation
}

// Recommendation is a pair consisting in the label of the recommendation along
// with the number of occurences in the dataset
type Recommendation struct {
    Item string
    Count int
}

// byCount sorts the list of recommendation by Recommendation#Count
type byCount []Recommendation
func (bc byCount) Len() int { return len(bc) }
func (bc byCount) Swap(i, j int)      { bc[i], bc[j] = bc[j], bc[i] }
func (bc byCount) Less(i, j int) bool { return bc[i].Count < bc[j].Count }

// Gold retrieves the gold standard list of recommendations.
// Endpoint is the address of the SPARQL endpoint, from is the named graph,
// query is the SPARQL query with the POF, and template is the text/template
// used for generating the Recommendation query.
func Gold(endpoint string, from string, query string, template string) []Recommendation {
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
    bindings, _ := GetBindings(endpoint, s.RecommendationQuery())
    var recs []Recommendation
    for _,v := range bindings {
        count,_ := strconv.Atoi(v["count"]["value"])
        recs = append(recs, Recommendation{ Item : v["POF"]["value"], Count : count })
    }
    return recs
}

// Measure measures the list of recommendations retrieved for the given query.
// Endpoint is the address of the SPARQL endpoint, from is the named graph,
// query is the SPARQL query with the POF, and template is the text/template
// used for generating the Recommendation query.
func Measure(endpoint string, from string, query string, template string) Measurement {
    pofs, elapsedTime := getRecommendations(endpoint, from, query, template)
    min, max, sum := math.MaxInt32, 0, float32(0)
    for _,c := range pofs {
        if c.Count < min {
            min = c.Count
        }
        if c.Count > max {
            max = c.Count
        }
        sum += float32(c.Count)
    }
    return Measurement{ Min: min, Max: max, Avg: sum / float32(len(pofs)), ElapsedTime: elapsedTime, Length: len(pofs), Recs: pofs }
}

// getRecommendations retrives the list of recommendation for the given query
// and rank it based on Recommendation#count. It returns the list of
// recommendations along with the time it took to retrieve them.
// Endpoint is the address of the SPARQL endpoint, from is the named graph,
// query is the SPARQL query with the POF, and template is the text/template
// used for generating the Recommendation query.
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
    bindings, elapsedTime := GetBindings(endpoint, s.RecommendationQuery())
    // get the POF bindings and rank them
    counts := make(map[string]int, len(bindings))
    for _,v := range bindings {
        count,_ := strconv.Atoi(v["count"]["value"])
        counts[v["POF"]["value"]] += count
    }
    log.Printf("Results: %v\n", counts)
    var pofs byCount
    for k,v := range counts {
        pofs = append(pofs, Recommendation{ Item: k, Count: v })
    }
    sort.Sort(sort.Reverse(pofs))
    min := int(math.Min(10, float64(len(pofs))))
    top := pofs[:min]
    log.Printf("TOP10: %v\n", top)
    // get the total number of occurrences of each recommended item
    values := "values ?POF { "
    for _,r := range top {
        values += "<" + r.Item + "> "
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
    bindings, _ = GetBindings(endpoint, s.RecommendationQuery())
    var popularity []Recommendation
    for _,v := range bindings {
        count,_ := strconv.Atoi(v["count"]["value"])
        popularity = append(popularity, Recommendation{ Item: v["POF"]["value"], Count: count })
    }
    log.Printf("Popularity=%v\n", bindings)
    return popularity, elapsedTime
}

