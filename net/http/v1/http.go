package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandleFunc)

	http.HandleFunc("/hel/", helloHandleFunc)

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandleFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, req.URL.Path)
}

func helloHandleFunc(w http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header { // req.Header map[string][]string
		fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
	}
}
