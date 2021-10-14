package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v4/protopb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	// 监听指定 port，获得一个 *grpc.ClientConn 实例
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial addr: %s error", port)
	}
	defer conn.Close()

	// 使用这个 *grpc.ClientConn 实例，创建指定的 GreeterClient 实例
	c := protopb.NewGreeterClient(conn)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 调用 GreeterClient 实例的方法，和在本地调用方法是一样的，这就是 RPC 带来的便捷
	reply, err := c.SayHello(ctx, &protopb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet:%v", err)
	}
	log.Printf("Get reply:%s", reply.GetMessage())
}
