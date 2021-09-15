package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
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
		obj := map[string]interface{}{
			"name":     "geektutu",
			"password": 1234,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(obj); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "404 page not found")
	}
}
