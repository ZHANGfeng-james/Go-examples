package rpc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
)

type Client struct {
	seq uint64      // 用于给请求编号，每个请求拥有唯一编号
	cc  codec.Codec // 消息的编解码器，序列化将要发出去的请求，反序列号接收到的响应
	opt *Option

	mu      sync.Mutex       // 支持对 pending 的并发读写
	pending map[uint64]*Call // Client 被保留所有已发出去的的请求，format: seq-*Call

	sending sync.Mutex   // 确保请求的有序发送，防止出现多个请求报文混淆
	header  codec.Header // 请求的消息头

	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

var ErrShutdown = errors.New("connection is shut down")

func DialWithAddr(rpcAddr string, opts ...*Option) (client *Client, err error) {
	params := strings.Split(rpcAddr, "@")
	if len(params) != 2 {
		return nil, fmt.Errorf("rpc client, wrong format address: %s, expect: protocol@addr", rpcAddr)
	}

	network, address := params[0], params[1]
	switch network {
	case "tcp":
		return Dial(network, address, opts...)
	case "http":
		return DialHTTP(network, address, opts...)
	default:
		return nil, fmt.Errorf("rpc client, wrong network protocol:%s, expect:tcp/http...", network)
	}
}

func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewClient, network, address, opts...)
}

func DialHTTP(network, address string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewClientHTTP, network, address, opts...)
}

func NewClientHTTP(conn net.Conn, opt *Option) (*Client, error) {
	Info("NewClientHTTP write to net.Conn")

	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", defaultRPCPath))
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	Info("resp.Status:%s", resp.Status)
	if err == nil && resp.Status == connected {
		return NewClient(conn, opt)
	}
	if err == nil {
		err = errors.New("unexpected HTTP status:" + resp.Status)
	}
	return nil, err
}

type newClientFunc func(conn net.Conn, opt *Option) (*Client, error)

func dialTimeout(f newClientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	// 在 Client 中封装 net.Dial
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeout)
	if err != nil {
		Info("net.DialTimeout:", err)
		return nil, err
	}
	// Dial --> NewClient --> newClient --> client.receive() 此处的 defer 在最后执行
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	ch := make(chan clientResult)
	go func() {
		client, err := f(conn, opt)
		ch <- clientResult{
			client: client,
			err:    err,
		}
	}()
	if opt.ConnectionTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}

	select {
	case <-time.After(opt.ConnectionTimeout):
		Info("timeout")
		return nil, fmt.Errorf("rpc client: connect server timeout, expect with %s", opt.ConnectionTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}

type clientResult struct {
	client *Client
	err    error
}

func parseOptions(opts ...*Option) (*Option, error) {
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	return opt, nil
}

func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		Info("rpc client: codec error:", err)
		return nil, err
	}
	// Client 发送给 Server 的格式：| Option | Header1 | Body1 | Header2 | Body2 |...
	// 也就是让 Server 知道 Client 当前的协议格式，一种协商措施
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		Info("rpc client: options error:", err)
		_ = conn.Close()
		return nil, err
	}
	return newClient(f(conn), opt), nil
}

func newClient(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1,
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}

	// 如何避免 goroutine 泄漏？
	go client.receive() // 启动接收，那 send 在哪执行？

	return client
}

func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			// 退出 for 循环
			break
		}
		// h.Seq 就是 Client 发送给 Server 的 sequence
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			//FIXME 什么情况下会出现？如果入参是 nil，会发送什么情况，io.Reader 如何解析？
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			// 本次调用结束时，通知调用方
			call.done()
		default:
			// 填充 call.Reply
			err := client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			// 本次调用结束时，通知调用方
			call.done()
		}
	}

	client.terminateCalls(err)
}

func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]

	delete(client.pending, seq)
	return call
}

func (client *Client) terminateCalls(err error) {
	// 当有多个 defer 语句时，其执行顺序类似入栈后出栈
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()

	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

func (client *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 同步调用，持续阻塞(<- channel)
	call := client.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-ctx.Done():
		client.removeCall(call.Seq)
		return errors.New("rpc client: call failed, " + ctx.Err().Error())
	case call := <-call.Done: // 用于接收 call.Error
		return call.Error
	}
}

func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 1)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	// Call 数据结构封装了一次 Client 的调用
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	// 异步执行，调用 Go 后立即返回
	return call
}

func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()

	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	// Client 封装的 Call 发送到 Server 端
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq // 更新本次注册的 *Call 实例中的 seq 字段
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.closing && !client.shutdown
}
