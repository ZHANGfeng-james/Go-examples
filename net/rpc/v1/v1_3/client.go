package main

import (
	"context"
	"fmt"
	"log"
	"net/rpc"
	"time"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

const serverAddr = ""

func main() {
	client, err := rpc.DialHTTP("tcp", serverAddr+":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := &v1.Args{A: 7, B: 8}
	quotient := new(v1.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)

	select {
	case <-timeout.Done():
		fmt.Println("timeout")
	case <-divCall.Done:
		fmt.Printf("Arith.Divide: %d / %d = %d.%d\n", args.A, args.B, quotient.Que, quotient.Rem)
	}

	fmt.Println("OVER")
}
