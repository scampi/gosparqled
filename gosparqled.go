package main

import (
    "github.com/gopherjs/gopherjs/js"
    "github.com/scampi/gosparqled/autocompletion"
)

// Scope as a global variable so that the text/template is created only once
var scope = autocompletion.NewScope()

// RecommendationQuery returns a SPARQL query for retrieving recommendations.
// If the input query does not have a Point Of Focus, an empty string is returned
func RecommendationQuery(query string, callback func(string, autocompletion.Type, string)) {
    go func(query string) {
        s := &autocompletion.Sparql{ Buffer : query, Scope : scope, Skip : &autocompletion.Skip{} }
        s.Init()
        autocompletion.Reset(s)
        err := s.Parse()
        if err == nil {
            s.Execute()
            callback(s.RecommendationQuery(), s.RecommendationType(), "")
        } else {
            callback(query, autocompletion.NONE, "Failed to process query\n" + err.Error())
        }
    }(query)
}

func main() {
    js.Global.Set("autocompletion", map[string]interface{}{
        "RecommendationQuery": RecommendationQuery,
        "PATH": autocompletion.PATH,
    })
}
