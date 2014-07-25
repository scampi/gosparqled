package autocompletion

import (
    "testing"
    "bytes"
)

type templateData struct {
    Tps []triplePattern
    Keyword, Pof, Prefix string
}

func (td *templateData) add(s string, p string, o string) {
    td.Tps = append(td.Tps, triplePattern{ S : s, P : p, O : o })
}

func (td *templateData) setPrefix(prefix string) {
    td.Prefix = prefix
}

func (td *templateData) setKeyword(kw string) {
    td.Keyword = kw
}

// Get the RecommendationQuery from query and compare it against the expected one
func parse(t *testing.T, query string, expected *templateData, rType Type) {
    parseWithPof(t, query, expected, "?POF", rType)
}

// Like parse but use pof as the POF variable expression
func parseWithPof(t *testing.T, query string, expected *templateData, pof string, rType Type) {
    s := &Sparql{ Buffer : query, Scope : NewScope() }
    s.Init()
    if err := s.Parse(); err != nil {
        t.Errorf("Failed to parse query\n%v", err)
    }
    s.Execute()
    actual := s.RecommendationQuery()
    var out bytes.Buffer
    expected.Pof = pof
    s.template.Execute(&out, expected)
    expectedString := out.String()
    if actual != expectedString {
        t.Errorf("Expected %v\nbut got %v\n", expectedString, actual)
    }
    aType := s.RecommendationType()
    if rType != aType {
        t.Errorf("Expected Recommendation type to be [%v] but got [%v]\n", rType, aType)
    }
}

func TestPrefixRecommendation1(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?POF", "?FillVar")
    td.setPrefix("aaa")
    parse(t, `
        PREFIX a: <aaa>
        SELECT *
        WHERE {
            ?s a:< 
        }
        `, td, PREDICATE)
}

func TestPrefixRecommendation2(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?POF", "?FillVar")
    td.setPrefix("aaa")
    parse(t, `
        PREFIX : <aaa>
        SELECT *
        WHERE {
            ?s :< 
        }
        `, td, PREDICATE)
}

func TestPrefixRecommendation3(t *testing.T) {
    td := &templateData{}
    td.add("?s", "a", "?POF")
    td.setPrefix("bbb")
    parse(t, `
        PREFIX a: <aaa>
        PREFIX b: <bbb>
        PREFIX c: <ccc>
        SELECT *
        WHERE {
            ?s a b:< 
        }
        `, td, CLASS)
}

func TestPrefixRecommendation4(t *testing.T) {
    td := &templateData{}
    td.add("?s", ":bbb", "?o")
    td.add("?s", "?POF", "?FillVar")
    parse(t, `
        PREFIX : <aaa>
        SELECT * WHERE {
          ?s :bbb ?o; <
        } 
        LIMIT 10
        `, td, PREDICATE)
}

func TestPath1(t *testing.T) {
    td := &templateData{}
    td.add("?sFillVar2", "?POF3", "?FillVar")
    td.add("?s", "?POF1", "?sFillVar1")
    td.add("?sFillVar1", "?POF2", "?sFillVar2")
    parseWithPof(t, `
        SELECT *
        WHERE {
          ?s 3/< 
        }
        `, td, pathPof(3), PATH)
}

func TestPath2(t *testing.T) {
    td := &templateData{}
    td.add("?s", "a", "<aaa>")
    td.add("?sFillVar2", "?POF3", "?FillVar")
    td.add("?s", "?POF1", "?sFillVar1")
    td.add("?sFillVar1", "?POF2", "?sFillVar2")
    parseWithPof(t, `
        SELECT *
        WHERE {
          ?s a <aaa>; 3/< 
        }
        `, td, pathPof(3), PATH)
}

func TestEval1(t *testing.T) {
    td := &templateData{}
    td.add("?v0", "a", "?POF")
    td.add("?v1", "<http://dbpedia.org/ontology/developer>", "?v0")
    td.add("?v1", "a", "<http://dbpedia.org/ontology/Software>")
    parse(t, `
        SELECT *
        WHERE {
          ?v0 a  <  .
          ?v1 <http://dbpedia.org/ontology/developer> ?v0 .
          ?v1 a <http://dbpedia.org/ontology/Software> .
        }
        `, td, CLASS)
}

func TestEval2(t *testing.T) {
    td := &templateData{}
    td.add("?v0", "a", "?POF")
    td.add("?v0", "<http://dbpedia.org/ontology/director>", "?v1")
    td.add("?v0", "<http://xmlns.com/foaf/0.1/name>", "?v2")
    td.add("?v0", "<http://dbpedia.org/property/imdbId>", "?v3")
    td.add("?v1", "<http://dbpedia.org/property/dateOfBirth>", "?v4")
    parse(t, `
        SELECT *
        WHERE {
            ?v0 a  <  .
            ?v0 <http://dbpedia.org/ontology/director> ?v1 .
            ?v0 <http://xmlns.com/foaf/0.1/name> ?v2 .
            ?v0 <http://dbpedia.org/property/imdbId> ?v3 .
            ?v1 <http://dbpedia.org/property/dateOfBirth> ?v4 .
        }
        `, td, CLASS)
}

func TestEval3(t *testing.T) {
    td := &templateData{}
    td.add("?v0", "a", "?POF")
    td.add("?v0", "<http://dbpedia.org/ontology/birthdate>", "?v1")
    td.add("?v0", "<http://xmlns.com/foaf/0.1/name>", "?v2")
    td.add("?v0", "<http://dbpedia.org/property/abstract>", "?v3")
    parse(t, `
        SELECT *
        WHERE {
            ?v0 a  <  ;<http://dbpedia.org/ontology/birthdate> ?v1 ;<http://xmlns.com/foaf/0.1/name> ?v2 ;<http://dbpedia.org/property/abstract> ?v3 .
        }
        `, td, CLASS)
}

func TestKeyword1(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?POF", "?FillVar")
    td.setKeyword("test")
    parse(t, `
        SELECT * WHERE {
          ?s test< 
        }
        LIMIT 10
    `, td, PREDICATE)
}

func TestKeyword2(t *testing.T) {
    td := &templateData{}
    td.add("?s", "?p", "?POF")
    td.setKeyword("test")
    parse(t, `
        SELECT * WHERE {
          ?s ?p test< 
        }
        LIMIT 10
    `, td, OBJECT)
}

func TestKeyword3(t *testing.T) {
    td := &templateData{}
    td.add("?s", "a", "?POF")
    td.setKeyword("Person-1")
    parse(t, `
        SELECT * WHERE {
          ?s a Person-1< 
        }
        LIMIT 10
    `, td, CLASS)
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
    `, td, PREDICATE)
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
    `, td, SUBJECT)
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
    `, td, PREDICATE)
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
    `, td, OBJECT)
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
    `, td, OBJECT)
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
    `, td, PREDICATE)
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
    `, td, PREDICATE)
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
    `, td, PREDICATE)
}

