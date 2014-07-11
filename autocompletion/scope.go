package autocompletion

import (
    "strings"
    "text/template"
    "bytes"
)

type triplePattern struct {
    S, P, O string
}

type Bgp struct {
    triplePattern
    Tps []triplePattern
    scope map[string]bool
    Template *template.Template
    Keyword string
}

func (b *Bgp) setKeyword(keyword string) {
    if len(keyword) != 0 {
        b.Keyword = keyword
    }
}

func (b *Bgp) setSubject(s string) {
    s = strings.TrimSpace(s)
    if (len(s) == 0) {
        return
    }
    b.S = s
}

func (b *Bgp) setPredicate(p string) {
    p = strings.TrimSpace(p)
    if (len(p) == 0) {
        return
    }
    b.P = p
}

func (b *Bgp) setObject(o string) {
    o = strings.TrimSpace(o)
    if (len(o) == 0) {
        return
    }
    b.O = o
}

func (b *Bgp) addTriplePattern() {
    tp := triplePattern{ S : b.S, P : b.P, O : b.O }
    b.Tps = append(b.Tps, tp)
}

func (b *Bgp) trimToScope() {
    b.scope = map[string]bool{ "?POF" : true }
    size := 0
    for size != len(b.scope) {
        size = len(b.scope)
        for _, tp := range b.Tps {
            if (tp.in(b.scope)) {
                tp.addToScope(b.scope)
            }
        }
    }
    var scoped []triplePattern
    for _,tp := range b.Tps {
        if (tp.in(b.scope)) {
            scoped = append(scoped, tp)
        }
    }
    b.Tps = scoped
}

func (tp *triplePattern) in(scope map[string]bool) bool {
    if scope[tp.S] || scope[tp.P] || scope[tp.O] {
        return true
    }
    return false
}

func (tp *triplePattern) addToScope(scope map[string]bool) {
    scope[tp.S] = true
    scope[tp.P] = true
    scope[tp.O] = true
}

func (b *Bgp) RecommendationQuery() string {
    b.trimToScope()
    var out bytes.Buffer
    b.Template.Execute(&out, b)
    return out.String()
}

