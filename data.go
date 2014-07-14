package main

import (
    "os"
    "github.com/scampi/gosparqled/eval/data"
)

func main() {
    data.Clean("http://ssdtest.index.sindice.net:4747/sparql", "http://sindice.com/usewod/dbpedia-3-3", data.Load(os.Args[1]), "clean"+os.Args[1])
}

