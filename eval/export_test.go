package eval

import "testing"

func TestJaccard(t *testing.T) {
    a := []Recommendation{ Recommendation{ item: "aaa", }, Recommendation{ item: "bbb" } }
    b := []Recommendation{ Recommendation{ item: "aaa", }, Recommendation{ item: "ccc" } }
    j := jaccard(a, b)
    if j != float64(1)/float64(3) {
        t.Errorf("Expected %v, but got %v", float64(1)/float64(3), j)
    }
}

