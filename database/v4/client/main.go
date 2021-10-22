package main

import (
	"context"
	"log"
	"time"

	"github.com/go-examples-with-tests/database/v4/pb"
	"google.golang.org/grpc"
)

var port = ":8081"

func main() {
	// 监听指定 port，获得一个 *grpc.ClientConn 实例
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial addr: %s error", port)
	}
	defer conn.Close()

	c := pb.NewCacheClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	limit := int64(1)
	msg := "Katyusha"
	reply, err := c.ListUsers(ctx, &pb.ListUsersRequest{
		Limit: &limit,
		Msg:   &msg,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("reply.Count:%d", reply.GetCount())
	items := reply.GetItems()
	if items != nil && len(items) > 0 {
		for _, item := range items {
			log.Printf("%s, %s, %s, %s", item.Nickname, item.Password, item.Phone, item.Email)
		}
	}
}
