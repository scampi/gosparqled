package eval

import (
    "net/http"
    "net/url"
    "github.com/golang/glog"
    "encoding/json"
    "io"
    "time"
)

// Binding maps a variable name to its solution
type Binding map[string]string

// executeQuery executes the SPARQL query over the endpoint
// and returns the io.ReadClosers body if successful
func executeQuery(endpoint string, query string) (io.ReadCloser, time.Duration) {
    time.Sleep(time.Second)
    q := endpoint + "?format=application/json&query=" + url.QueryEscape(query)
    glog.Infof("Execute request: [%s]", q)
    start := time.Now()
    resp, err := http.Get(q)
    if err != nil {
        glog.Fatal(err)
    }
    return resp.Body, time.Since(start)
}

// GetBindings returns the list of Bindings for the query over the endpoint.
func GetBindings(endpoint string, query string) ([]map[string]Binding, time.Duration) {
    body, et := executeQuery(endpoint, query)
    defer body.Close()
    dec := json.NewDecoder(body)
    var res = new(struct{Results struct{Bindings []map[string]Binding}})
    if err := dec.Decode(&res); err != nil {
        glog.Fatal(err)
    }
    return res.Results.Bindings, et
}

// Ask returns the result of the ASK query over the endpoint.
func Ask(endpoint string, query string) (bool, time.Duration) {
    body, et := executeQuery(endpoint, query)
    defer body.Close()
    dec := json.NewDecoder(body)
    var res = new(struct{Boolean bool})
    if err := dec.Decode(&res); err != nil {
        glog.Fatal(err)
    }
    return res.Boolean, et
}

