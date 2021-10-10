package rpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/discover"
	"github.com/go-examples-with-tests/net/rpc/v2/registry"
)

func runRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9999")
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}

func runServer(wg *sync.WaitGroup) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")

	server := NewServer()
	_ = server.Register(&foo)

	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 5*time.Second)
	wg.Done()

	server.Accept(l)
}

const registryAddr = "http://localhost:9999/_geerpc_/registry"

func TestXClient(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go runRegistry(&wg)
	wg.Wait()

	time.Sleep(2 * time.Second)

	wg.Add(3)
	go runServer(&wg)
	go runServer(&wg)
	go runServer(&wg)
	wg.Wait()

	d := discover.NewGeeRegistryDiscovery(registryAddr, 0)
	xclient := NewXClient(d, discover.RoundRobinSelect, nil)
	defer func() {
		_ = xclient.Close()
	}()

	var work sync.WaitGroup
	for i := 0; i < 5; i++ {
		work.Add(1)

		go func(i int) {
			defer work.Done()
			var reply int
			args := &Args{Num1: i, Num2: i * i}
			// err := xclient.Call(context.Background(), "Foo.Sum", args, &reply)
			// if err != nil {
			// 	log.Printf("%s:%s, err: %v", "Call", "Foo.Sum", err)
			// } else {
			// 	log.Printf("%s %s success: %d + %d = %d", "Call", "Foo.Sum", args.Num1, args.Num2, reply)
			// }

			err := xclient.Broadcast(context.Background(), "Foo.Sum", args, &reply)
			if err != nil {
				log.Printf("%s:%s, err: %v", "Broadcast", "Foo.Sum", err)
			} else {
				log.Printf("%s %s success: %d + %d = %d", "Broadcast", "Foo.Sum", args.Num1, args.Num2, reply)
			}
		}(i)
	}
	work.Wait()

	time.Sleep(10 * time.Second)
}
