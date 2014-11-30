[![GoDoc](https://godoc.org/github.com/scampi/gosparqled?status.svg)](https://godoc.org/github.com/scampi/gosparqled) [![Build Status](https://travis-ci.org/scampi/gosparqled.svg?branch=master)](https://travis-ci.org/scampi/gosparqled)

# Gosparqled - SPARQL Auto-Completion

Gosparqled provides a library for retrieving context-aware recommendations for a SPARQL query. The library is written in GO and is then translated into JavaScript using [GopherJS](https://github.com/gopherjs/gopherjs).

![Gosparqled demo](anim.gif)

# What is _context-aware_ ?

By context-aware recommendation, I mean that all patterns which are connected to the element to recommend participate in the recommendation. For example, the 3<sup>rd</sup> pattern below is left out of the context of the recommendation. The reason is that none of its variables appear in the others. However, the 2<sup>nd</sup> is kept.

```
?s rdf:type %% RECOMMEND HERE! %% .
?s foaf:name ?name .
?o ?p "left out" .
```

# Demo

The folder `demo/` shows how gosparqled can be used with other tools such as [YASR](https://github.com/YASGUI/YASR) and [YASQE](https://github.com/YASGUI/YASQE) in order to have a full-fledged SPARQL query editor with the added recommendation feature. That demo can be tested at [http://scampi.github.io/gosparqled/](http://scampi.github.io/gosparqled/).

# Auto-Completion

The recommendations are context-aware, i.e., the current state of the SPARQL query is taken into consideration when retrieving the recommendations.

In the examples below, the character `<` represents the position in the query to auto-complete by pressing `CTRL + SPACE`. The `<` should not be typed prior to pressing the key combination. The auto-completion is possible at any position in a triple pattern.

## Class

Recommend possible classes:

```sql
SELECT * {
    ?s a <
}
``` 

## Predicate

Recommend possible predicates:

```sql
SELECT * {
    ?s <
}
``` 

## Relation

Recommend possible relations between a `Person` and a `Document`:

```sql
SELECT * {
    ?s a <Person> .
    ?o a <Document> .
    ?s < ?o
}
``` 

## Keyword

Recommend possible terms (e.g., classes or predicates) which URI contains a keyword, case-insensitive. Below, it presents classes which contain the word `movie`:

```sql
SELECT * {
    ?s a Movie<
}
```

## Path

Recommend possible path of a fixed length, written as `X/`. Below, recommendations about paths of lengths 2 between a `Movie` and a `Person` are returned:

```sql
SELECT * {
    ?s a <Movie> .
    ?o a <Person> .
    ?s 2/< ?o
}
```
## Prefix

Recommend possible terms (e.g., classes or predicates) with the given prefix. Below, it presents only the predicates within the `rdfs` prefix:

```sql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

SELECT * {
    ?s rdfs:<
}
```

## Content

Recommend content elements, e.g., literals, URIs, either at the subject or the object location.

- Auto-completion on the object. This returns a sample of the labels, probably literals.

    ```sql
    SELECT * {
        ?s rdfs:label <
    }
    ```

- Auto-completion on the subject. This returns a list of URIs which name is John.

    ```sql
    SELECT * {
        < foaf:name "John" .
    }
    ```

# Building

First, install GopherJS:

```sh
$ go get github.com/gopherjs/gopherjs
```

Run the command below to create the JavaScript library. The `-m` flag minifies the generated JavaScript code.

```sh
$ gopherjs build -m
```

The following methods can then be called via JavaScript (see `demo/autocompletion.js`).

- `RecommendationQuery` in the `autocompletion` namespace

    It takes in the SPARQL query with the character `<` indicating the position in the query to auto-complete. It returns the processed SPARQL query, which can then be sent to the SPARQL endpoint in order to retrieve the possible recommendations. The recommendations are bound to the variable `?POF`.

# Publication

This library is presented in [http://ceur-ws.org/Vol-1272/paper_157.pdf](http://ceur-ws.org/Vol-1272/paper_157.pdf). If you are using this tool, please cite this work.
