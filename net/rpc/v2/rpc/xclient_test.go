package rpc

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/go-examples-with-tests/net/rpc/v2/discover"
)

func runServer(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")

	server := NewServer()
	_ = server.Register(&foo)

	addrCh <- l.Addr().String()

	server.Accept(l)
}

func TestXClient(t *testing.T) {
	addrCh1 := make(chan string)
	addrCh2 := make(chan string)
	addrCh3 := make(chan string)

	go runServer(addrCh1)
	go runServer(addrCh2)
	go runServer(addrCh3)

	rpcAddr := make([]string, 0)
	rpcAddr = append(rpcAddr, <-addrCh1)
	rpcAddr = append(rpcAddr, <-addrCh2)
	rpcAddr = append(rpcAddr, <-addrCh3)

	d := discover.NewMultiServersDiscovery([]string{
		"tcp@" + rpcAddr[0],
		"tcp@" + rpcAddr[1],
		"tcp@" + rpcAddr[2]},
	)

	xclient := NewXClient(d, discover.RandomSelect, nil)
	defer func() {
		_ = xclient.Close()
	}()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
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
	wg.Wait()
}
