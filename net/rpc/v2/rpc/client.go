package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
)

type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	call.Done <- call
}

var ErrShutdown = errors.New("connection is shut down")

func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	// 在 Client 中封装 net.Dial
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn, opt)
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
		log.Println("rpc client: codec error:", err)
		return nil, err
	}
	// Client 发送给 Server 的格式：| Option | Header1 | Body1 | Header2 | Body2 |...
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error:", err)
		_ = conn.Close()
		return nil, err
	}
	return newClient(f(conn), opt), nil
}

type Client struct {
	seq uint64      // 用于给请求编号，每个请求拥有唯一编号
	cc  codec.Codec // 消息的编解码器，序列化将要发出去的请求，反序列号接收到的响应
	opt *Option

	mu      sync.Mutex       // 支持对 pending 的并发读写
	pending map[uint64]*Call // Client 被保留（未处理）的请求，format: seq-*Call

	sending sync.Mutex   // 确保请求的有序发送，防止出现多个请求报文混淆
	header  codec.Header // 请求的消息头

	closing  bool // user has called Close
	shutdown bool // server has told us to stop
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

func (client *Client) Call(serviceMethod string, args, reply interface{}) error {
	// 同步调用，持续阻塞(<- channel)
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
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
	call.Seq = client.seq
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
