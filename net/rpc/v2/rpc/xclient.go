package rpc

import (
	"context"
	"reflect"
	"sync"

	"github.com/go-examples-with-tests/net/rpc/v2/discover"
)

type XClient struct {
	discover   discover.Discover
	selectMode discover.SelectMode
	opt        *Option
	clients    map[string]*Client // 保留所有已和指定string(server addr)创建 net.Conn 的 *Client 实例
	mu         sync.Mutex
}

func NewXClient(d discover.Discover, mode discover.SelectMode, opt *Option) *XClient {
	return &XClient{
		discover:   d,
		selectMode: mode,
		opt:        opt,
		clients:    make(map[string]*Client),
	}
}

func (xclient *XClient) Close() error {
	xclient.mu.Lock()
	defer xclient.mu.Unlock()

	for addr, client := range xclient.clients {
		//FIXME I have no idea how to deal with error, just ignore it.
		_ = client.Close()
		delete(xclient.clients, addr)
	}
	return nil
}

// 向服务列表的某个 Server 发起请求，基于某种 discover.SelectMode
func (xclient *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 依据 mode 挑选出 server 端地址，取名为：rpcAddr
	rpcAddr, err := xclient.discover.Get(xclient.selectMode)
	Infof("XClient Call select rpcAddr:%s", rpcAddr)

	if err != nil {
		return err
	}
	// 调用 call
	return xclient.call(rpcAddr, ctx, serviceMethod, args, reply)
}

func (xclient *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 依据 addr 调用 dial，并返回一个 *Client
	client, err := xclient.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

func (xclient *XClient) dial(rpcAddr string) (*Client, error) {
	xclient.mu.Lock()
	defer xclient.mu.Unlock()

	client, ok := xclient.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		// 存在但不可用
		_ = client.Close()
		delete(xclient.clients, rpcAddr)
		client = nil
	}

	if client == nil {
		var err error
		client, err = DialWithAddr(rpcAddr, xclient.opt)
		if err != nil {
			return nil, err
		}
		xclient.clients[rpcAddr] = client
	}

	return client, nil
}

// 向服务列表的所有 Server 发起请求
func (xclient *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	clients, err := xclient.discover.GetAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error

	replyDone := reply == nil

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < len(clients); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 并发请求时，确保能够拷贝结果正常，每个 goroutine 创建新的 clonedReply
			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}

			err = xclient.Call(ctx, serviceMethod, args, clonedReply)
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel() // 如果出现一次请求异常，结束所有请求
			}

			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
			}

			mu.Unlock()
		}()
	}

	wg.Wait()
	return nil
}
