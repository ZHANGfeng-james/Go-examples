package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandleFunc)

	http.HandleFunc("/hello", helloHandlFunc)

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandleFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, req.RequestURI)
}

func helloHandlFunc(w http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header { // req.Header map[string][]string
		fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
	}
}
