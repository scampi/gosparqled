package data

import "testing"

func TestPOFs(t *testing.T) {
    pofs := POFs(`
        select * {
            ?s a <Person> .
            ?s <name> ?name
        }
    `)
    if len(pofs) != 2 {
        t.Error("Wrong size")
    }
    pof1 := `
        select * {
            ?s a  <  .
            ?s <name> ?name
        }
    `
    pof2 := `
        select * {
            ?s a <Person> .
            ?s  <  ?name
        }
    `
    if pofs[0] != pof1 {
        t.Errorf("Expected %s but got %s", pof1, pofs[0])
    }
    if pofs[1] != pof2 {
        t.Errorf("Expected %s but got %s", pof2, pofs[1])
    }
}

