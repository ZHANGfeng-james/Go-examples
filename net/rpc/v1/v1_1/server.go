package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

func main() {
	arith := new(v1.Arith)
	rpc.Register(arith)

	rpc.HandleHTTP() // HTTP 监听器

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(listener, nil)
}
