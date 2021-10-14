package main

import (
	"context"
	"log"
	"net"

	"github.com/go-examples-with-tests/net/rpc/v4/protopb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	protopb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *protopb.HelloRequest) (*protopb.HelloReply, error) {
	// 如果 Client 给定的 HelloRequest 是 nil，则表示 Client 根本没有传入
	if req.Name == nil {
		return &protopb.HelloReply{Message: "nobody"}, nil
	}
	// 此处处理的情况，包括 HelloRequest 传入的 req.Name 是空字符串的情况
	log.Printf("Received: %s", req.GetName())
	return &protopb.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen:%s", port)
	}

	// 创建一个 gRPC Server 实例
	s := grpc.NewServer()
	protopb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("falied to serve:%v", err)
	}
}
