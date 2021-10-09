RPC——远程过程调用，是一种**计算机通信协议**，允许调用**不同进程空间**的程序。RPC 的客户端和服务器可以在一台机器上，也可以在不同的机器上。程序员使用时，就像调用本地程序一样，**无需关注内部的实现细节**。

**不同的应用程序**之间的**通信方式**有很多，比如浏览器和服务器之间广泛使用的**基于 HTTP 协议的 RESTfull API 标准**。与 RPC 相比，RESTfull API 有相对统一的标准，因而更通用，兼容性更好，支持不同的语言。HTTP 协议是**基于文本的**，一般具备**更好的可读性**。但是**缺点**也很明显：

- RESTfull 接口需要额外的定义，无论是客户端还是服务端，都需要额外的代码来处理，而 RPC 调用则**更接近于直接调用**。
- 基于 HTTP 协议的 RESTfull 报文冗余，承载了过多的无效信息，而 RPC 通常使用**自定义的协议格式**，减少冗余报文。
- RPC 可以采用**更高效的序列化协议**，将文本转为二进制传输，获得更高的性能。
- 因为 RPC 的灵活性，所以更容易扩展和集成诸如注册中心、负载均衡等功能。

从底层网络传输的内容来看，实际上就是本文讨论的主题是 RPC，其本身就是**基于传输二进制数据**的 TCP 协议的应用层协议。与之不同的是，HTTP 的底层也是 TCP 协议的应用层协议（**基于文本**）。

**RPC 框架需要解决什么问题？为什么需要 RPC 框架**？

我们可以想象下**两台机器上，两个应用程序之间需要通信**，那么首先，需要确定采用的**传输协议**是什么？如果这两个应用程序位于不同的机器，那么一般会选择 TCP 协议或者 HTTP 协议；那如果两个应用程序位于相同的机器，也可以选择 Unix Socket 协议。传输协议确定之后，还需要确定**报文的编码格式**，比如采用最常用的 JSON 或者 XML，那如果报文比较大，还可能会选择 protobuf 等其他的编码方式，甚至编码之后，再进行压缩。接收端获取报文则需要相反的过程，先解压再解码。

> Protobuf 是一种数据表示格式，是 Google 出品的一种描述数据内容的方式和格式，传输效率高。

解决了传输协议和报文编码的问题，接下来还需要解决一系列的**可用性问题**，例如，连接超时了怎么办？是否支持异步请求和并发？

如果服务端的实例很多，客户端并不关心这些实例的地址和部署位置，只关心自己能否获取到期待的结果，那就引出了**注册中心(registry)和负载均衡(load balance)的问题**。简单地说，即客户端和服务端互相不感知对方的存在，服务端启动时将自己注册到注册中心，客户端调用时，从注册中心获取到所有可用的实例，选择一个来调用。这样服务端和客户端只需要感知注册中心的存在就够了。注册中心通常还需要实现服务动态添加、删除，使用心跳确保服务处于可用状态等功能。

再进一步，假设服务端是不同的团队提供的，如果没有**统一的 RPC 框架**，各个团队的服务提供方就需要各自实现一套消息编解码、连接池、收发线程、超时处理等“业务之外”的重复技术劳动，造成整体的低效。因此，**“业务之外”的这部分公共的能力，即是 RPC 框架所需要具备的能力**。

> RPC 框架本质上是要解决端之间的数据通信问题。

Go 语言广泛地应用于**云计算**和**微服务**，**成熟的 RPC 框架和微服务框架**汗牛充栋。`grpc`、`rpcx`、`go-micro` 等都是非常成熟的框架。一般而言，**RPC 是微服务框架的一个子集**，**微服务框架**可以自己实现 **RPC 部分**，当然，也可以选择不同的 RPC 框架作为通信基座。

考虑性能和功能，上述成熟的框架代码量都比较庞大，而且通常和第三方库，例如 `protobuf`、`etcd`、`zookeeper` 等有比较深的耦合，难以直观地窥视**框架的本质**。GeeRPC 的目的是以最少的代码，**实现 RPC 框架中最为重要的部分**，帮助大家理解 RPC 框架在设计时需要考虑什么。代码简洁是第一位的，功能是第二位的。

因此，**GeeRPC 选择从零实现 Go 语言官方的标准库 `net/rpc`**，并在此基础上，新增了协议交换(protocol exchange)、注册中心(registry)、服务发现(service discovery)、负载均衡(load balance)、超时处理(timeout processing)等特性。分七天完成，最终代码约 1000 行。

> 从上面这句内容：“GeeRPC 选择从零实现 Go 语言官方的标准库 net/rpc”，我大概知道了本文的目标。

# 1 服务端与消息编码

一个典型的 RPC 调用如下：

~~~go
err = client.Call("Arith.Multiply", args, &reply)
~~~

客户端发送的请求包括服务名 `Arith`，对应服务下的某个方法 `Multiply`，以及发送给这个方法的入参。紧接着的是返回值：reply，以及调用的状态反馈 err。

我们将请求的参数和返回值抽象在 Body 中，剩余的信息放在 Header 中，那么就可以抽象出数据结构 Header：

~~~go
type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}
~~~

上面说的 Header 和 Body 部分就是对于一个 HTTP 通信来说的，将一个消息划分为相同的结构。ServiceMethod 是服务名和方法名，通常与 Go 语言中的结构体类型名和方法名相映射。Seq 是请求的序列号，也可以认为是某个请求的 ID，用来区分不同的请求，是有 Client 端给定。

最终的传输内容格式设计成：

~~~bash
| Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
~~~

进一步抽象出对消息体进行编解码的接口 Codec，**抽象出接口**是为了**实现不同的 Codec 实例**：

~~~go
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
~~~

GobCodec 作为 Codec 的一种，需要实现 4 个方法：

~~~go
package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}
}

func (gob *GobCodec) ReadHeader(h *Header) error {
	return gob.dec.Decode(h)
}

func (gob *GobCodec) ReadBody(body interface{}) error {
	return gob.dec.Decode(body)
}

func (gob *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
        // buf --> conn
		_ = gob.buf.Flush()
		if err != nil {
			_ = gob.Close()
		}
	}()

	// write header --> conn
	if err := gob.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	// write body --> conn
	if err := gob.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (gob *GobCodec) Close() error {
	return gob.conn.Close()
}
~~~

上面整个过程实现了**消息的序列化和反序列化**，也就是通过 encoding/gob 实现了 Encode/Decode 过程。

客户端与服务端的通信需要**协商一些内容**，例如 HTTP 报文，**分为 Header 和 Body 部分**，body 的格式和长度通过 Header 中**的 Content-Type 和 Content-Length 指定**，服务端通过解析 Header 就能够知道如何从 body 中读取需要的信息。对于 RPC 协议来说，这部分协商是需要自主设计的。

为了提升性能，一般在报文的最开始会规划固定的字节，来协商相关的信息。比如：第 1 个字节用来表示序列化方式，第 2 个字节表示压缩方式，第 3～6字节表示 header 的长度，7～10字节表示 body 的长度。对于 GeeRPC 来说，目前需要协商的唯一一项内容是**消息的编解码方式**：

~~~go
package rpc

import "github.com/go-examples-with-tests/net/rpc/v2/codec"

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // 标记这是 geerpc 的 request
	CodecType   codec.Type // client 还可使用其他的 codec 用于编码 body 部分
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType, // 默认情况下，RPC 服务端使用 gob codec 解码
}
~~~

一般来说，涉及协议协商的这部分信息，需要设计**固定的字节来传输**。但是为了实现上更简单，GeeRPC 客户端固定采用 JSON 编码 Option，后续的 header 和 body 的编码方式由 Option 中的 CodeType 指定：

~~~bash
| Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
~~~

**在一次连接（net.Conn）中**，Option 固定在报文的最开始，Header 和 Body 可以有多个，即报文可能是这样的：

~~~bash
| Option | Header1 | Body1 | Header2 | Body2 |...
~~~

接下来就要去实现 Server 的部分：

~~~go
type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept() // net.Listener --> Accept() --> net.Conn
		if err != nil {
			log.Println("rpc server: accept error, ", err)
			return
		}
		go server.ServeConn(conn) // 开启 goroutine 处理 net.Conn
	}
}

func Accept(lis net.Listener) {
    // net.Listener 从哪里来？ l, err := net.Listen("tcp", ":0")
	DefaultServer.Accept(lis)
}
~~~

创建了 Server 结构体，实现了 Accept 方法，使用 net.Listener 作为参数，for 循环中等待 Socket 连接建立，并开启 goroutine 处理，处理过程交给 ServeConn 方法。

如果想要启动服务，过程是非常简单的，传入 net.Listener 实例即可，TCP 协议和 UNIX 协议都支持：

~~~go
listener, _ := net.Listen("tcp", ":9999")
geerpc.Accept(listener)
~~~

紧接着实现 ServeConn 方法：

~~~go
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	// | Option | Header1 | Body1 | Header2 | Body2 |...
	// option 方面使用 JSON 格式编码，最先解析的是 json 格式的 Option
	//FIXME json.NewDecoder(conn).Decode(&opt) 的工作原理？
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: Options error, ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn))
}

var invalidRequest = struct{}{}

// 相当于是 Option 后续的 Header 和 Body 部分，一次 net.Conn 仅需要传输一次 Option
type request struct {
	h            *codec.Header // Header
	argv, replyv reflect.Value // Body
}

type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}

// f(conn) 得到的是一个 codec.Codec 编解码器
func (server *Server) serveCodec(cc codec.Codec) {
	// 注意，此处使用的是 *sync.Mutex 和 *sync.WaitGroup
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		// 读取请求
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}

			req.h.Error = err.Error() // string 格式的
			// 回复请求
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		// 处理请求 和 回复请求
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
    // | Option | {Header1 | Body1} | {Header2 | Body2} |...
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	// 通过 cc.ReadBody 修改 req.argv 的值，req.argv 在当前是作为一个 string 类型
	req.argv = reflect.New(reflect.TypeOf("")) // *string --> reflect.Value
    // 作为一个 codec，ReadHeader 和 ReadBody 时，需要标记已读取的字节序号
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
    // sending sync.Mutex 避免发送数据过程中并发导致数据混乱
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// 处理请求部分，仅打印 request 内容
	log.Println(req.h, req.argv.Elem())
	// 处理请求后，为 reply 设置值
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerepc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
~~~

Client 在请求 Server 时，其通信格式是：`| Option | Header1 | Body1 | Header2 | Body2 |...`

Server 在接收到 Client 请求后，会依次解析出 Option，紧接着是 Header1 和 Body1。整个处理逻辑依次是：读取请求、处理请求和回复请求。

测试程序：

~~~go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
	"github.com/go-examples-with-tests/net/rpc/v2/rpc"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0") // Server 端程序是以 Listen 开始
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	rpc.Accept(l) // 接收 net.Listener
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr) // Client 端程序是以 Dial 开始
	defer func() {
		_ = conn.Close()
	}()

	time.Sleep(5 * time.Second)

	// 写入 Option
	// | Option | Header1 | Body1 | Header2 | Body2 |... 写入 json 格式的 option
	_ = json.NewEncoder(conn).Encode(rpc.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i), // Client 端忽略 Error
		}
		// write head and body
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))

		_ = cc.ReadHeader(h)

		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
~~~

执行结果如下：

~~~bash
ant@MacBook-Pro v2 % go run main.go
2021/10/07 15:50:16 start rpc server on [::]:59378
2021/10/07 15:50:21 &{Foo.Sum 0 } geerpc req 0
2021/10/07 15:50:21 reply: geerepc resp 0
2021/10/07 15:50:21 &{Foo.Sum 1 } geerpc req 1
2021/10/07 15:50:21 reply: geerepc resp 1
2021/10/07 15:50:21 &{Foo.Sum 2 } geerpc req 2
2021/10/07 15:50:21 reply: geerepc resp 2
2021/10/07 15:50:21 &{Foo.Sum 3 } geerpc req 3
2021/10/07 15:50:21 reply: geerepc resp 3
2021/10/07 15:50:21 &{Foo.Sum 4 } geerpc req 4
2021/10/07 15:50:21 reply: geerepc resp 4
~~~

Client 在发出请求时，需要在消息的头部添加 Option 内容，但对于 Server 来说，写入的反馈就不需要 Option 内容了。上述测试用例中，并发 RPC 请求 Server 端的数据，其中 Option 仅在创建了 net.Conn 之后发送一次，后续不再发送。而 Header 会多次发送，对应 Body 也会多次发送。

> 此处的疑问是：Body 部分的 req.argv 是如何被解析出来的？

# 2 支持并发与异步的客户端

在上一节内容中，主要是实现了服务端程序，也就是说，客户端能够发起网络请求，并能获取到 Server 返回的响应。

那本节内容实际上就是实现的是 **net/rpc 标准库的 Client** 的基本功能：**发出请求**和**接收反馈**。也就是说，经过本节内容，就可以实现大致和 net/rpc 相同的功能。

先来看看在实现客户端后的测试程序：

~~~go
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/rpc"
)

func startServer(addr chan string) {
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

	client, _ := rpc.Dial("tcp", <-addr)
	defer func() {
		// 原先是 net.Conn
		_ = client.Close()
	}()

	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup // 实现并发控制
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
            // 从调用形式来看，是和 net/rpc 一样的
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
~~~

从测试程序来看，这个 Client 端程序不再看到 net.Conn，没有了关于 Option 的写入，也没有了对服务端反馈消息的解析。这个调用过程和 net/rpc 是一样的形式。

抽象出 Client 结构体：

~~~go
type Client struct {
	seq uint64      // 用于给请求编号，每个请求拥有唯一编号
	cc  codec.Codec // 消息的编解码器，序列化将要发出去的请求，反序列化接收到的响应
	opt *Option

	mu      sync.Mutex       // 支持对 pending 的并发读写
	pending map[uint64]*Call // Client 被保留（未处理）的请求，format: seq-*Call

	sending sync.Mutex   // 确保请求的有序发送，防止出现多个请求报文混淆
	header  codec.Header // 请求的消息头

	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}
~~~

Client 的字段比较复杂：

- cc 是消息的编解码器，和服务端类似，用来序列化将要发送出去的请求，以及反序列化接收到的响应。
- sending 是一个互斥锁，和服务端类似，为了**保证请求的有序发送**，即**防止出现多个请求报文混淆**。
- header 是每个请求的消息头，header 只有在请求发送时才需要，而请求发送是互斥的，因此每个客户端只需要一个，声明在 Client 结构体中可以复用。
- seq 用于**给发送的请求编号**，每个请求拥有唯一编号。
- pending 存储未处理完的请求，键是编号（seq 的值），值是 Call 实例。
- closing 和 shutdown 任意一个值置为 true，则表示 Client 处于不可用的状态，但有些许的差别，closing 是**用户主动关闭的**，即调用 `Close` 方法，而 shutdown 置为 true 一般是**有错误发生**。

启动 Client，以及创建 Client 实例：

~~~go
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
    // 也就是让 Server 知道 Client 当前的协议格式，一种协商措施
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error:", err)
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
~~~

创建 Client 实例时，首先需要完成一开始的协议交换，即发送 Option 信息给服务端。协商好消息的编解码方式之后，再创建一个 goroutine 调用接收 Request。

Client 调用一次 RPC 请求，抽象成一个 Call 实例：

~~~go
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
~~~

Call 中的所有字段，承载了一次 RPC 调用所需要的全部信息。在结构体中增加了 Done，为了支持异步调用，在获取到 RPC 反馈后，会调用 done 通知调用方。

接下来是一系列和 Call 相关的方法：

~~~go
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

func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.closing || client.shutdown { // client 当前已被关闭
		return 0, ErrShutdown
	}
    
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}
~~~

上面 Client 和 Call 相关的方法，实际上就是 client.pending 维护的 `seq-*Call` 的映射关系。

对一个客户端来说，**接收响应**、**发送请求**是最重要的 2 个功能。`client.cc.ReadHeader` 和 `client.cc.ReadBody` 持续读取 net.Conn 中的数据：

~~~go
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
~~~

 接收到的响应有三种情况：

- call 不存在，可能是请求没有发送完整（Client 先发送的是 Header，紧接着发送了 Body，可能是 Body 发送出错），或者因为其他原因被取消，但是服务端仍旧处理了。
- call 存在，但服务端处理出错，即 h.Error 不为空。
- call 存在，服务端处理正常，那么需要从 body 中读取 Reply 的值。

接下来是**发送功能**：

~~~go
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
~~~

最后是入口功能：

~~~go
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
~~~

Call 结构体中的 Done 实际上就是用来**支持异步调用**的。

测试程序输出结果：

~~~bash
ant@MacBook-Pro v2 % go run main.go
start rpc server on [::]:59950
&{Foo.Sum 5 } geerpc req 2
&{Foo.Sum 2 } geerpc req 4
&{Foo.Sum 4 } geerpc req 1
&{Foo.Sum 1 } geerpc req 0
&{Foo.Sum 3 } geerpc req 3
reply: geerepc resp 3
reply: geerepc resp 2
reply: geerepc resp 4
reply: geerepc resp 1
reply: geerepc resp 5
~~~

# 3 服务注册

RPC 框架的一个基础能力是：像调用本地程序一样调用远程服务。基于前 2 节的内容，对于 Go 来说，这个问题就变成了**如何将结构体的方法映射为服务**。

对 net/rpc 而言，一个函数需要能够被远程调用，需要满足如下 5 个条件：

1. 方法所属的类型是可导出的，比如下述类型 `T`；
2. 方法是可导出的；
3. 方法有 2 个参数，都是可导出类型或内建类型；
4. 方法的第二个参数是指针；
5. 方法只有一个 error 接口类型的返回值。

假设客户端发过来一个请求，包含 ServiceMethod 和 Argv：

~~~bash
{
    "ServiceMethod"： "T.MethodName"
    "Argv"："0101110101..." // 序列化之后的字节流
}
~~~

通过 `T.MethodName` 可以确定调用的是类型 T 的 MethodName，如果**硬编码**实现这个功能，可能是这样的：

~~~go
...
switch req.ServiceMethod {
    case "T.MethodName":
    	t := new(t)
        reply := new(T2)
    
        var argv T1
        gob.NewDecoder(conn).Decode(&argv)
    
        err := t.MethodName(argv, reply)
        server.sendMessage(reply,err)
    case "Foo.Sum":
   		...
}
...
~~~

也就是说，如果使用硬编码的方式来实现结构体与服务的映射，那么**每暴露一个方法，就需要编写等量的代码**。那么有没有什么方式，能够**将这个映射过程自动化**？

~~~go
func main() {
	var wg sync.WaitGroup
	// sync.WaitGroup 中定义的是 *sync.WaitGroup 为接收者的方法
	typ := reflect.TypeOf(&wg)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i) // reflect.Method

		argv := make([]string, 0, method.Type.NumIn())     // the type of method, Func
		returns := make([]string, 0, method.Type.NumOut()) // Func

		// 第 0 个入参是 wg 自己
		for j := 1; j < method.Type.NumIn(); j++ {
			argv = append(argv, method.Type.In(j).Name()) // the jth input parameter type name
		}

		for j := 0; j < method.Type.NumOut(); j++ {
			returns = append(returns, method.Type.Out(j).Name())
		}

		log.Printf("func (w *%s) %s(%s) %s",
			typ.Elem().Name(),
			method.Name,
			strings.Join(argv, ","),
			strings.Join(returns, ","))
	}
}
~~~

通过反射，我们能够非常容易地获取某个结构体的所有方法，并且能够通过方法，获取到该方法所有的参数类型与返回值。上述程序的运行结果：

~~~go
ant@MacBook-Pro v2 % go run main.go
2021/10/08 09:49:55 func (w *WaitGroup) Add(int) 
2021/10/08 09:49:55 func (w *WaitGroup) Done() 
2021/10/08 09:49:55 func (w *WaitGroup) Wait() 
~~~

实现服务注册功能（通过结构体名，以及对应的方法名，对应就能调用这个方法，同时附带有入参和输出值），封装结构体方法信息：

~~~go
// 一个 method 的所有完整信息
type methodType struct {
	// func (t *T)MethodName(argType T1, replyType *T2) error
	// 一个 method 的所有信息包括：方法名（统一到 Func 这种类型值上），入参，返回值
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}
~~~

与之对应的方法：

~~~go
func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value

    // 创建 reflect.Value 准备接收 request.argv
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem()) // reflect.Type.Elem() --> reflect.Type
	} else {
		argv = reflect.New(m.ArgType).Elem() // reflect.Value.Elem()
	}
	return argv
}

func (m *methodType) newReplyv() reflect.Value {
	// reply must be a pointer type，这是 RPC 协议规定的
	replyv := reflect.New(m.ReplyType.Elem())

	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
        // 初始化
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
        // 初始化
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}
~~~

每一个 methodType 实例包含了一个方法的完整信息，包括：

* method：方法本身
* ArgType：第一个参数的类型，也就是入参；
* ReplyType：第二个参数的类型，也就是出参；
* numCalls：后续统计方法调用次数时会用到。

定义 service 结构体，用于表示**某个结构体信息**：

~~~go
type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

func newService(receiver interface{}) *service {
	s := new(service)

	s.rcvr = reflect.ValueOf(receiver) // 这个结构体的实例

	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(receiver)

	// 判断 struct name 是否是可导出的
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	log.Printf("new Service name:%s", s.name)
	s.registerMethods()

	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type

		// 方法的第一个入参是接收者本身
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}

		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
~~~

RPC 方法的 2 个参数，必须是可导出的，而且还带有一个 error 类型的返回值。

~~~go
func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)

	f := m.method.Func

	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
~~~

在实际调用结构体对应方法时，需要使用 `s.rcvr` 作为第一个参数，也就是方法的接收者。

service 的测试程序：

~~~go
package rpc

import (
	"fmt"
	"reflect"
	"testing"
)

type Foo int

type Args struct {
	num1 int
	num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.num1 + args.num2
	return nil
}

func (f Foo) sum(args Args, reply *int) error {
	*reply = args.num1 + args.num2
	return nil
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func TestNewService(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	_assert(len(s.method) == 1, "wrong service Method, expect 1, but got %d", len(s.method))

	mType := s.method["Sum"]
	_assert(mType != nil, "wrong Method, Sum should not nil")
}

func TestMethodType_Call(t *testing.T) {
	var foo Foo
	s := newService(&foo)

	mType := s.method["Sum"]

	argv := mType.newArgv()
	replyv := mType.newReplyv()
	argv.Set(reflect.ValueOf(Args{num1: 1, num2: 3}))
	err := s.call(mType, argv, replyv)

	_assert(err == nil && *replyv.Interface().(*int) == 4 && mType.numCalls == 1, "failed to call Foo.Sum")
}
~~~

通过反射结构体已经映射为服务，但请求的处理过程还没有完成。从接收到请求到回复还有如下步骤待实现：

1. 根据入参类型，将请求的 body 反序列化；
2. 调用 service.call 完成方法调用；
3. 将 reply 序列化为字节流，构造响应报文。

将服务的注册过程集成到 Server 中：

~~~go
type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}
~~~

为 Server 新增一个字段，表示该 Server 中已注册的所有 Service。

相应的，实现从 Server 中查找指定名称 Server 的方法：

~~~go
func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc service: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svic, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can not find service " + methodName)
		return
	}
	svc = svic.(*service)
	mtype = svc.method[methodName] // 从 service 中获取到指定名称的方法
	if mtype == nil {
		err = errors.New("rpc server: can not find method " + methodName)
	}
	return
}
~~~

另外对于一次 `<Header | Body>`，为 request 新增字段：

~~~go
// 也就是一次连接中，紧接在 Option 之后的部分：Header 和 Body
type request struct {
	h            *codec.Header // Header
	argv, replyv reflect.Value // Body

	mtype *methodType
	svc   *service
}
~~~

因此，从Request中读取到 Header 后就可以解析出对应的 service：

~~~go
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	// | Option | Header1 | Body1 | Header2 | Body2 |...
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	// 依据 methodName 获取对应的入参类型和返回值类型，并构造实例
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		// 必须要能够修改这个值
		argvi = req.argv.Addr().Interface()
	}
	// 字节流反序列化成变量值
	if err = cc.ReadBody(argvi); err != nil {
		// gob: type mismatch in decoder: want struct type main.Args; got non-struct
		log.Println("rpc server: read body err:", err)
		return req, err
	}

	return req, nil
}
~~~

紧接着就是执行对应的 RPC 方法：

~~~go
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	err := req.svc.call(req.mtype, req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
~~~

具体的：

~~~go
func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)

	// m.method 是 reflect.Method 类型 --> reflect.Value 类型，且其 Kind 是 Func
	f := m.method.Func

	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
~~~

测试代码：

~~~go
package main

import (
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

	client, _ := rpc.Dial("tcp", <-addr)
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

			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}

	wg.Wait()
}
~~~

# 4 超时处理

超时处理是 RPC 框架一个比较基本的能力，如果缺少**超时处理机制**，无论是**服务端**还是**客户端**都容易因为网络或其他错误导致挂起，资源耗尽，这些问题的出现大大降低了**服务的可用性**。因此，我们需要在 RPC 框架中加入超时处理的能力。

纵观整个远程调用的过程，**需要客户端处理超时**的地方有：与服务端建立连接；发送请求到服务端，写报文时；等待服务端处理，等待处理，而服务端因为某些原因已宕机，无法响应；从服务端接收响应结果，读报文。**需要服务端处理超时**的地方有：读取客户端请求报文；调用映射服务的方法，处理报文导致超时异常；发送响应报文，写报文导致超时。

**「客户端处理连接超时」**：

~~~go
package rpc

import (
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // 标记这是 geerpc 的 request
	CodecType   codec.Type // client 还可使用其他的 codec 用于编码 body 部分

	ConnectionTimeout time.Duration
	HandleTimeout     time.Duration
}

var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         codec.GobType, // 默认情况下，RPC 服务端使用 gob codec 解码
	ConnectionTimeout: 3 * time.Second,
}
~~~

为 Option 新增字段：ConnectionTimeout，默认是 3 s。

~~~go
func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewClient, network, address, opts...)
}

type newClientFunc func(conn net.Conn, opt *Option) (*Client, error)

func dialTimeout(f newClientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	// 在 Client 中封装 net.Dial
	log.Println("start to dial...")
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeout)
	if err != nil {
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
		// time.Sleep(6 * time.Second) 此处模拟网络超时
		ch <- clientResult{
			client: client,
			err:    err,
		}
	}()
	if opt.ConnectionTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}

	log.Println("begin to time count...")
	select {
	case <-time.After(opt.ConnectionTimeout):
		return nil, fmt.Errorf("rpc client: connect server timeout, expect with %s", opt.ConnectionTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}

type clientResult struct {
	client *Client
	err    error
}
~~~

修改原先 Dial 函数，新增 dialTimeout 函数，在其中实现连接超时控制。连接超时控制的实现有如下步骤：

1. `net.DialTimeout`：连接函数调用时，附带上超时时间；
2. 使用 time.After，如果在超时时间内 ch 获取 clientResult，说明未超时。

「Client.Call 的**超时机制**」将控制权交给用户：

~~~go
...
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
...
~~~

对应的实现：

~~~go
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
~~~

**「服务器端超时处理」**:

~~~go
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, handleTimeout time.Duration) {
	defer wg.Done()

	called := make(chan struct{})
	sent := make(chan struct{})

	// 如果超时，则此处会导致 goroutine 泄漏
	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if handleTimeout == 0 { // 不做超时控制
		<-called
		<-sent
		return
	}

	select {
	case <-time.After(handleTimeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout, expect within:%s", handleTimeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}
~~~

这里需要确保 `sendResponse` 仅调用一次，因此将整个过程拆分为 `called` 和 `sent` 两个阶段，在这段代码中只会发生如下两种情况：sere

1) called 信道接收到消息，代表处理没有超时，继续执行 sendResponse。
2) `time.After()` 先于 called 接收到消息，说明处理已经超时，called 和 sent 都将被阻塞。在 `case <-time.After(timeout)` 处调用 `sendResponse`。

但此处有一个问题，如果超时，那么在其中的 **goroutine 会导致泄漏**！

测试程序：

~~~go
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
~~~

# 5 支持 HTTP 协议

Web 开发中，我们经常使用 HTTP 协议中的 HEAD/GET/POST 等方式发送请求，等待响应。但**本文所实现的 RPC 的消息格式**与**标准的 HTTP 协议**并不兼容，在这种情况下，就需要**一个协议的转换过程**。HTTP 协议的 **CONNECT 方法**恰好提供了这个能力，CONNECT 一般用于**代理服务**。

假设**浏览器**与**服务器**之间的 HTTPS 通信都是**加密的**，浏览器通过**代理服务器**发起 HTTPS 请求时，由于请求的站点地址和端口号都是加密保存在 HTTPS 请求报文头中的，代理服务器如何知道往哪里发送请求呢？为了解决这个问题，可进行如下步骤：

1. 浏览器通过 HTTP **明文形式**向代理服务器发送**一个 CONNECT 请求**告诉代理服务器**目标地址和端口**；
2. 代理服务器接收到这个请求后，会在对应端口与目标站点**建立一个 TCP 连接**，连接建立成功后返回 HTTP 200 状态码告诉浏览器与该站点的加密通道已经完成；
3. 之后浏览器和服务器开始 HTTPS 握手并交换加密数据，代理服务器**仅需透传浏览器和服务器之间的加密数据包**即可，代理服务器无需解析 HTTPS 报文。

事实上，这个过程其实是通过代理服务器将 HTTP 协议转换为 HTTPS 协议的过程。对于 RPC 服务端来说，需要做的是将 HTTP 协议转换为 RPC 协议；对于客户端来说，需要新增通过 HTTP CONNECT 请求创建连接的逻辑。

**「服务端支持HTTP协议」**

整个通讯过程：

1. 客户端向 RPC 服务器发送 CONNECT 请求；
2. RPC 服务器返回 HTTP 200 状态码表示连接建立成功；
3. 客户端使用创建好的连接（net.Conn）发送 RPC 报文，先发送 Option，再发送请求报文。服务端处理 RPC 请求并响应。

~~~go
const (
	connected        = "200 Connected to Gee RPC"
	defaultRPCPath   = "/_geerpc_"
	defaultDebugPath = "/debug/geerpc"
)

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}

func (server *Server) HandleHTTP() {
	// 启动 HTTP Server 端，同时监听的 path 是：defaultRPCPath 和 defaultDebugPath
	http.Handle(defaultRPCPath, server)
	http.Handle(defaultDebugPath, debugHTTP{server})
}

// 接收到的是 HTTP 协议的内容，自动转化到 ServeHTTP 方法中
func (server *Server) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	Info("request.Method:%s", request.Method)

	if request.Method != http.MethodConnect {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(rw, "405 must CONNECT method\n")
		return
	}

	conn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacker, remote:", request.RemoteAddr, ": ", err.Error())
		return
	}
	_, _ = io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	server.ServeConn(conn)
}
~~~

复用建立的 net.Conn 实例，并转化到 RPC 协议上，和原先是一样的。

**「Client 端实现协议协议转换」**：

~~~go
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
~~~

其关键部分就是，Client 执行 Dial 后，向 net.Conn 写入符合 HTTP 协议的 「起始行/Header/Body」部分，使用的是 CONNECT 方法。同时，接收 Server 端的指定反馈，最后转换到 RPC 协议上，和之前完全一样。

测试程序：

~~~go
func main() {
	addrCh := make(chan string)
	go func(addrCh chan string) {
		var foo Foo
		if err := rpc.Register(&foo); err != nil {
			log.Fatal("register error:", err)
		}

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal("network error:", err)
		}
		log.Println("start rpc server on", l.Addr())
		addrCh <- l.Addr().String()

		rpc.HandleHTTP()
		_ = http.Serve(l, nil) // 启动HTTP服务器
	}(addrCh)

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

	time.Sleep(10 * time.Second)
}
~~~

# 6 负载均衡







# 7 服务发现与注册中心





