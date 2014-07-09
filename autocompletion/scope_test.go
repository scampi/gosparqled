package autocompletion

import (
    "testing"
    "github.com/hoisie/mustache"
)

type templateData struct {
    Tps []triplePattern
}

func (td *templateData) add(s string, p string, o string) {
    td.Tps = append(td.Tps, triplePattern{ S : s, P : p, O : o })
}

func parse(t *testing.T, query string, expected *templateData) {
    tp, _ := mustache.ParseFile("/home/stecam/documents/prog/go/src/github.com/scampi/gosparqled/autocompletion/template.mustache")
    s := &Sparql{ Buffer : query, Bgp : &Bgp{Template : tp} }
    s.Init()
    if err := s.Parse(); err != nil {
        t.Errorf("Failed to parse query\n%v", err)
    }
    s.Execute()
    actual := s.RecommendationQuery()
    if actual != tp.Render(expected) {
        t.Errorf("Expected %v\nbut got %v\n", tp.Render(expected), actual)
    }
}

func TestEditor(t *testing.T) {
    td := &templateData{}
    td.add("?sub", "a", "<http://schema.org/MusicGroup>")
    td.add("?sub", "?POF", "?FillVar")
    parse(t, `
SELECT * WHERE {
  ?sub a <http://schema.org/MusicGroup> .
  ?sub < 
}
LIMIT 10
    `, td)
}

func TestSubject(t *testing.T) {
    td := &templateData{}
    td.add("?POF", "?p", "?o1")
    td.add("?POF", "a", "?o2")
    parse(t, `
        select * {
            < ?p ?o1; a ?o2 .
            ?o ?op ?oo .
            ?a ?b ?c
        }
    `, td)
}

func TestPredicate(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?POF", "?FillVar")
    td.add("?s", "a", "?o")
    td.add("?o", "?op", "?oo")
    parse(t, `
        select * {
            ?s < ; a ?o .
            ?o ?op ?oo .
            ?a ?b ?c
        }
    `, td)
}

func TestObject(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?p", "?POF")
    td.add("?s", "a", "?o")
    td.add("?o", "?op", "?oo")
    parse(t, `
        select * {
            ?s ?p < ; a ?o .
            ?o ?op ?oo .
            ?a ?b ?c
        }
    `, td)
}

func TestTerms(t *testing.T) {
    td := &templateData{}
    td.add("?s", "<p1>", "?POF")
    td.add("?s", "a", "?o")
    td.add("?o", "?p", "\"test\"")
    parse(t, `
        select * {
            ?s <p1> < ; a ?o .
            ?o ?p "test" .
        }
    `, td)
}

func TestOptional1(t *testing.T) {
    td := &templateData{}
    td.add("?s", "<p1>", "?o")
    td.add("?o", "?POF", "?FillVar")
    parse(t, `
        select * {
            ?s <p1> ?o .
            OPTIONAL { ?o < } .
        }
    `, td)
}

func TestOptional2(t *testing.T) {
    td := &templateData{}
    td.add("?o", "?POF", "?FillVar")
    td.add("?s", "<p1>", "?o")
    parse(t, `
        select * {
            OPTIONAL { ?o < } .
            ?s <p1> ?o .
        }
    `, td)
}

func TestOptional3(t *testing.T) {
    td := &templateData{}
    td.add("?s", "<p1>", "?o")
    td.add("?o", "?POF", "?FillVar")
    td.add("?s", "<p1>", "?o")
    parse(t, `
        select * {
            ?s <p1> ?o .
            OPTIONAL { ?o < } .
            ?s <p1> ?o .
        }
    `, td)
}

