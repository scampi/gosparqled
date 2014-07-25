package sparql

import "testing"

func parse(t *testing.T, query string) {
    s := &Sparql{ Buffer : query }
    s.Init()
    if err := s.Parse(); err != nil {
        t.Errorf("Failed to parse query\n%v", err)
    }
}

func TestQName(t *testing.T) {
    parse(t, "select * from acme:test { ?s acme:p ?o }")
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

