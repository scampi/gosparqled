package sparql

import "testing"

// Parses the query and asserts that there is no error
func parse(t *testing.T, query string) *Sparql {
    s := &Sparql{ Buffer : query }
    s.Init()
    if err := s.Parse(); err != nil {
        t.Errorf("Failed to parse query\n%v", err)
    }
    return s
}

func TestFunctionCall(t *testing.T) {
    parse(t, "select ( <aaa>(?test) as ?e ) { ?s ?p ?o }")
    parse(t, `
            prefix : <aaa>
            select ( :(?test) as ?e ) {
                ?s ?p ?o
            }
        `)
}

func TestBuiltinCall(t *testing.T) {
    parse(t, "select ( str(?test) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( floor(?test) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( contains(?test, \"bla\") as ?e ) { ?s ?p ?o }")
    parse(t, "select ( bound(?test) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( bnode(?test) && bnode() as ?e ) { ?s ?p ?o }")
    parse(t, "select ( now() as ?e ) { ?s ?p ?o }")
    parse(t, "select ( concat(\"l\", \"u\", ?tt, \"f\", \"y\") as ?e ) { ?s ?p ?o }")
    parse(t, "select ( regex(?a, ?b) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( regex(?a, ?b, \"i\") as ?e ) { ?s ?p ?o }")
    parse(t, "select ( if(?i, ?t, ?e) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( exists { ?i ?t ?e } as ?e ) { ?s ?p ?o }")
}

func TestComments(t *testing.T) {
    parse(t, `   # this is sparta
    select  #blabla
            * { ?s # sub!
            ?p # pred !
            ?o # obj  !!
        }
    `)
}

func TestExpressions(t *testing.T) {
    parse(t, "select ( ?s as ?e ) { ?s ?p ?o }")
    parse(t, "select ( !?s -?s + ?s as ?e ) { ?s ?p ?o }")
    parse(t, "select ( ( ?a + ?b ) * 2 as ?e ) { ?s ?p ?o }")
    parse(t, "select ( ?s / 2 + 42 as ?e ) { ?s ?p ?o }")
    parse(t, "select ( ?s = ?o && ?p in ( ?o1, ?o2 ) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( <aaa> + \"23\" - true + ( 32 = ?s ) as ?e ) { ?s ?p ?o }")
    parse(t, "select ( ?s not in (?o) as ?e ) { ?s ?p ?o }")
}

func TestDescribe(t *testing.T) {
    parse(t, "describe *")
    parse(t, "describe ?p")
    parse(t, "describe <aaa>")
    parse(t, "describe * { ?s ?p ?o }")
}

func TestAsk(t *testing.T) {
    parse(t, "ASK { ?s ?p ?o }")
}

func TestConstruct(t *testing.T) {
    parse(t, "CONSTRUCT { ?a ?b ?c } { ?s ?p ?o }")
}

func TestQName1(t *testing.T) {
    parse(t, "select * from acme:test { ?s acme:p ?o }")
}

func TestQName2(t *testing.T) {
    parse(t, "prefix : <acme.org/> select * { ?s :p ?o }")
}

func TestTriplesSameSubject1(t *testing.T) {
    parse(t, "SELECT * { ?s ?p ?o }")
}

func TestTriplesSameSubject2(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o
        }
    `)
}

func TestTriplesSameSubject3(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o1, ?o2
        }
    `)
}

func TestTriplesSameSubject4(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o1; ?p ?o2, ?o3
        }
    `)
}

func TestTriplesSameSubject5(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s <1> | <2> ?o
        }
    `)
}

func TestTriplesSameSubject6(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s <1> / <2> ?o
        }
    `)
}

func TestTriplesSameSubject7(t *testing.T) {
    parse(t, `
        SELECT * {
            [ <1> ?a; <2> ?b ] <3> ?c
        }
    `)
}

func TestTriplesSameSubject8(t *testing.T) {
    parse(t, `
        SELECT * {
            ( <1> <2> <3> ) ?p ?o
        }
    `)
}

func TestOptional(t *testing.T) {
    parse(t, `
        SELECT * {
            OPTIONAL { <1> <2> <3> }
        }
    `)
}

func TestSubSelect(t *testing.T) {
    parse(t, `
        SELECT * {
            SELECT * {
                ?s ?p ?o
            }
        }
    `)
}

func TestUnion(t *testing.T) {
    parse(t, `
        SELECT * {
            { ?s ?p ?o } UNION { ?a ?b ?c }
        }
    `)
}

func TestLimitOffset(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o
        }
        limit 10
        offset 10
    `)
}

func TestPrefix(t *testing.T) {
    parse(t, ` 
        prefix a: <aaa>
        prefix b: <bbb>
        
        SELECT * {
            ?s ?p ?o
        }
    `)
}

func TestProjections(t *testing.T) {
    parse(t, `
        SELECT ?a {
            ?s ?p ?o
        }
    `)
}

func TestBGP1(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o .
            ?a ?b ?c .
        }
    `)
}

func TestBGP2(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s ?p ?o; ?b ?c .
        }
    `)
}

func TestGraphTerms(t *testing.T) {
    parse(t, `
        SELECT * {
            ?s a <Person>;
               <name> "Jean";
               <age> "42"^^<int>;
               <lang> "fr"@fr-be12;
               ?p true, false;
               ?p _:b1, [];
               ?p ()
        }
    `)
}

