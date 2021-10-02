package main

import (
	"fmt"
	"log"
	"net/http"

	v4 "github.com/go-examples-with-tests/net/http/v4"
)

func main() {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	v4.NewGroup("scores", 2<<10, v4.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key ", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := ":9999"
	peers := v4.NewHTTPPool(addr)
	log.Println("geecache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
