gosparqled
==========

Assisted SPARQL Editor written in GO

[![GoDoc](https://godoc.org/github.com/scampi/gosparqled?status.svg)](https://godoc.org/github.com/scampi/gosparqled)

# SPARQL Auto-Completion

Gosparqled provides a library for retrieving context-aware recommendations for a SPARQL query. The library is written in GO and is then translated into JavaScript using [GopherJS](https://github.com/gopherjs/gopherjs).

# Demo

The folder `demo/` shows how gosparqled can be used with other tools such as [YASR](https://github.com/YASGUI/YASR) and [YASQE](https://github.com/YASGUI/YASQE) in order to have a full-fledged SPARQL query editor with the added recommendation feature.

# Auto-Completion

The recommendations are context-aware, i.e., the current state of the SPARQL query is taken into consideration when retrieving the recommendations. In the examples below, the character `<` represents the position in the query to auto-complete. The auto-completion is possible at any position in a triple pattern.

## Class

Recommend possible classes:

```sparql
SELECT * {
    ?s a <
}
``` 

## Predicate

Recommend possible predicates:

```sparql
SELECT * {
    ?s <
}
``` 

## Relation

Recommend possible relations between two variables:

```sparql
SELECT * {
    ?s a :Person .
    ?o a :Document .
    ?s < ?o
}
``` 

## Keyword

Recommend possible terms (e.g., classes or predicates) which URI contains a keyword. Below, it presents classes which contain the word `movie`, case-insensitive:

```sparql
SELECT * {
    ?s a Movie<
}
```

# Building

First, install GopherJS:

```sh
$ go get github.com/gopherjs/gopherjs
```

Run the command below to create the JavaScript library. The `-m` flag minifies the generated JavaScript code.

```sh
$ gopherjs build -m gosparqled.go
```

The following methods can then be called via JavaScript (see `demo/sparqled.js`).

- `RecommendationQuery` in the `autocompletion` namespace

    It takes in the SPARQL query with the character `<` indicating the position in the query to auto-complete. It returns a processed SPARQL query, which can then be sent to the SPARQL endpoint in order to retrieve the recommendations. The recommendations are bound to the variable `?POF`.
