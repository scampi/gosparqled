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

func TestEscapedQuotes(t *testing.T) {
    parse(t, `
        SELECT *
        WHERE {
            ?s <aaa> "a\"b"
        }
        `)
}

func TestTwoTypes(t *testing.T) {
    parse(t, `
        SELECT *
        WHERE {
            ?s a :Person; a <aaa>
        }
        `)
}

func TestObjectNodePath(t *testing.T) {
    parse(t, `
            select * {
                ?s ?p [ <p1> ?o1; ]
            }
        `)
}

func TestTrailingSemiColon(t *testing.T) {
    parse(t, `
            select * {
                ?s <p1> ?o1;
            }
        `)
    parse(t, `
            select * {
                [ <p1> ?o1; ]
            }
        `)
}

func TestTriplesNodePath(t *testing.T) {
    parse(t, `
            select * {
                [
                    <p1> ?o1;
                    <p2> ?o2
                ]
            }
        `)
}

func TestMultipleDatasets(t *testing.T) {
    parse(t, `
            select *
            from <aaa>
            from <bbb> {
                ?s ?p ?o
            }
        `)
    parse(t, `
            construct { ?a ?b ?c }
            from <aaa>
            from <bbb> {
                ?s ?p ?o
            }
        `)
    parse(t, `
            ASK
            from <aaa>
            from <bbb> {
                ?s ?p ?o
            }
        `)
    parse(t, `
            describe ?s
            from <aaa>
            from <bbb> {
                ?s ?p ?o
            }
        `)
}

func TestGraph(t *testing.T) {
    parse(t, `
            select * {
                graph ?g {
                    ?s ?p ?o
                }
                graph <aaa> {
                    ?s ?p ?o
                }
                graph a:b {
                    ?s ?p ?o
                }
            }
        `)
}

func TestMinus(t *testing.T) {
    parse(t, `
            select * {
                minus {
                    ?s ?p ?o
                }
            }
        `)
}

func TestGroupBy(t *testing.T) {
    parse(t, `
            select (count(*) as ?c) {
                ?s ?p ?o
            }
            group by ( bound(?s) as ?sss )
        `)
}

func TestOrderBy(t *testing.T) {
    parse(t, `
            select (count(*) as ?c) {
                ?s ?p ?o
            }
            order by desc (?s + ?p)
        `)
}

func TestAggregate(t *testing.T) {
    parse(t, `
            select (count(*) as ?c) {
                filter ( count(?s) > 5 )
            }
        `)
}

func TestFilterOrBind1(t *testing.T) {
    parse(t, `
            select * {
                filter ( str(?s) = "en" )
                filter bound(?o)
            }
        `)
}

func TestFilterOrBind2(t *testing.T) {
    parse(t, `
            select * {
                ?s ?p ?o
                filter bound(?o)
            }
        `)
}

func TestFilterOrBind3(t *testing.T) {
    parse(t, `
            select * {
                filter bound(?o)
                ?s ?p ?o
            }
        `)
}

func TestFilterOrBind4(t *testing.T) {
    parse(t, `
            select * {
                filter xsd:long(?o)
            }
        `)
}

func TestFilterOrBind5(t *testing.T) {
    parse(t, `
            select * {
                bind(23 + ?price as ?o)
            }
        `)
}

func TestFunctionCall1(t *testing.T) {
    parse(t, "select ( <aaa>(?test) as ?e ) { ?s ?p ?o }")
}

func TestFunctionCall2(t *testing.T) {
    parse(t, `
            prefix : <aaa>
            select ( :doit(?test) as ?e ) {
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
            SELECT (count(*) as ?c) {
                ?s ?p ?o
            }
            group by ?s
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

