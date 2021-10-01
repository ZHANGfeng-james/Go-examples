package main

import (
	"fmt"
	"log"
	"net/http"
)

func test() {
	engine := &Engine{}
	log.Fatal(http.ListenAndServe(":9999", engine))
}

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	switch path {
	case "/":
		fmt.Fprintf(w, req.RequestURI)
	case "/hello":
		for key, value := range req.Header { // req.Header map[string][]string
			fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
		}
	default:
		fmt.Fprintf(w, "404 page not found")
	}
}
