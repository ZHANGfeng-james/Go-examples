package rpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"
)

type Bar int

func (bar *Bar) Timeout(argv int, replyv *int) error {
	log.Println("Bar timeout run")
	time.Sleep(2 * time.Second)
	return nil
}

func startServer(addr chan string) {
	var bar Bar
	if err := Register(&bar); err != nil {
		log.Fatal("register error:", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()

	log.Println(l.Addr().String())

	Accept(l) // 接收 net.Listener
}

func TestClientCall(t *testing.T) {
	t.Parallel()

	addrCh := make(chan string)
	go startServer(addrCh)

	addr := <-addrCh
	time.Sleep(3 * time.Second)
	t.Run("client timeout control", func(t *testing.T) {
		client, err := Dial("tcp", addr)
		if client == nil {
			log.Println(err)
			return
		}

		log.Println("client is normal")

		defer func() {
			// 原先是 net.Conn
			_ = client.Close()
		}()

		// 用户需要在 1s 内拿到服务端的响应结果
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		var reply int
		err = client.Call(ctx, "Bar.Timeout", 1, &reply)
		if err != nil {
			log.Println("err:", err)
		}
	})

	t.Run("server timeout handle", func(t *testing.T) {
		// 设定 Server 必须在 1 秒内处理结果，否则超时
		client, err := Dial("tcp", addr, &Option{
			HandleTimeout: time.Second,
		})
		if client == nil {
			log.Println(err)
			return
		}

		log.Println("client is normal")

		defer func() {
			// 原先是 net.Conn
			_ = client.Close()
		}()

		var reply int
		err = client.Call(context.Background(), "Bar.Timeout", 1, &reply)
		if err != nil {
			log.Println("err:", err)
		}
	})
}

func TestHttp(t *testing.T) {

}
