package main

import (
	"fmt"
	"log"
	"net/rpc"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

const serverAddr = ""

func main() {
	client, err := rpc.DialHTTP("tcp", serverAddr+":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}

	args := &v1.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith err:", err)
	}
	fmt.Printf("Arith.Multiply: %d*%d=%d", args.A, args.B, reply)
}
