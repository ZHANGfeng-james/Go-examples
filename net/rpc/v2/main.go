package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/rpc"
)

type Foo int

type Args struct {
	Num1 int
	Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	if err := rpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	rpc.Accept(l) // 接收 net.Listener
}

func main() {
	log.SetFlags(0)

	addr := make(chan string)
	go startServer(addr)

	client, err := rpc.Dial("tcp", <-addr, &rpc.Option{HandleTimeout: 2 * time.Second})
	if client == nil {
		log.Println(err)
		return
	}

	log.Println("client is normal")

	defer func() {
		// 原先是 net.Conn
		_ = client.Close()
	}()

	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			// defer cancel()

			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				rpc.Info("call Foo.Sum error:", err)
			} else {
				log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
			}
		}(i)
	}

	wg.Wait()
}

func test() {

	addrCh := make(chan string)
	go startServer(addrCh)

	client, _ := rpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(2 * time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
