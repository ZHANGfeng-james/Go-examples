package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	v4 "github.com/go-examples-with-tests/net/http/v4"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?") // 什么含义？
	flag.Parse()

	log.Printf("api:%v", api)

	apiAddr := "http://localhost:9999"
	peerMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range peerMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(peerMap[port], addrs, gee)
}

func createGroup() *v4.Group {
	return v4.NewGroup("scores", 2<<10, v4.GetterFunc(func(key string) ([]byte, error) {
		// 从本地DB中取值
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, gee *v4.Group) {
	httppool := v4.NewHTTPPool(addr)
	httppool.Set(addrs...)
	gee.RegistePeers(httppool)

	log.Println("geecache is running at", addr)

	log.Fatal(http.ListenAndServe(addr[7:], httppool))
}

func startAPIServer(apiAddr string, gee *v4.Group) {
	http.Handle("/api", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := gee.Get(key)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Write(view.ByteSlice())
		rw.Write([]byte("\r\n"))
	}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}
