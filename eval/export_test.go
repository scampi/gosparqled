package eval

import "testing"

func TestJaccard(t *testing.T) {
    a := []Recommendation{ Recommendation{ Item: "aaa", }, Recommendation{ Item: "bbb" } }
    b := []Recommendation{ Recommendation{ Item: "aaa", }, Recommendation{ Item: "ccc" } }
    j := jaccard(a, b)
    if j != float64(1)/float64(3) {
        t.Errorf("Expected %v, but got %v", float64(1)/float64(3), j)
    }
}

