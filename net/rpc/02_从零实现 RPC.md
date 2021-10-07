RPC——远程过程调用，是一种**计算机通信协议**，允许调用**不同进程空间**的程序。RPC 的客户端和服务器可以在一台机器上，也可以在不同的机器上。程序员使用时，就像调用本地程序一样，无需关注内部的实现细节。

不同的应用程序之间的通信方式有很多，比如浏览器和服务器之间广泛使用的基于 HTTP 协议的 RESTfull API 标准。与 RPC 相比，RESTfull API 有相对统一的标准，因而更通用，兼容性更好，支持不同的语言。HTTP 协议是**基于文本的**，一般具备**更好的可读性**。但是**缺点**也很明显：

- RESTfull 接口需要额外的定义，无论是客户端还是服务端，都需要额外的代码来处理，而 RPC 调用则**更接近于直接调用**。
- 基于 HTTP 协议的 RESTfull 报文冗余，承载了过多的无效信息，而 RPC 通常使用**自定义的协议格式**，减少冗余报文。
- RPC 可以采用**更高效的序列化协议**，将文本转为二进制传输，获得更高的性能。
- 因为 RPC 的灵活性，所以更容易扩展和集成诸如注册中心、负载均衡等功能。

RPC 框架需要解决什么问题？为什么需要 RPC 框架？

我们可以想象下**两台机器上，两个应用程序之间需要通信**，那么首先，需要确定采用的**传输协议**是什么？如果这两个应用程序位于不同的机器，那么一般会选择 TCP 协议或者 HTTP 协议；那如果两个应用程序位于相同的机器，也可以选择 Unix Socket 协议。传输协议确定之后，还需要确定**报文的编码格式**，比如采用最常用的 JSON 或者 XML，那如果报文比较大，还可能会选择 protobuf 等其他的编码方式，甚至编码之后，再进行压缩。接收端获取报文则需要相反的过程，先解压再解码。

解决了传输协议和报文编码的问题，接下来还需要解决一系列的**可用性问题**，例如，连接超时了怎么办？是否支持异步请求和并发？

如果服务端的实例很多，客户端并不关心这些实例的地址和部署位置，只关心自己能否获取到期待的结果，那就引出了**注册中心(registry)和负载均衡(load balance)的问题**。简单地说，即客户端和服务端互相不感知对方的存在，服务端启动时将自己注册到注册中心，客户端调用时，从注册中心获取到所有可用的实例，选择一个来调用。这样服务端和客户端只需要感知注册中心的存在就够了。注册中心通常还需要实现服务动态添加、删除，使用心跳确保服务处于可用状态等功能。

再进一步，假设服务端是不同的团队提供的，如果没有**统一的 RPC 框架**，各个团队的服务提供方就需要各自实现一套消息编解码、连接池、收发线程、超时处理等“业务之外”的重复技术劳动，造成整体的低效。因此，“业务之外”的这部分公共的能力，即是 RPC 框架所需要具备的能力。

Go 语言广泛地应用于**云计算**和**微服务**，**成熟的 RPC 框架和微服务框架**汗牛充栋。`grpc`、`rpcx`、`go-micro` 等都是非常成熟的框架。一般而言，RPC 是微服务框架的一个子集，微服务框架可以自己实现 RPC 部分，当然，也可以选择不同的 RPC 框架作为通信基座。

考虑性能和功能，上述成熟的框架代码量都比较庞大，而且通常和第三方库，例如 `protobuf`、`etcd`、`zookeeper` 等有比较深的耦合，难以直观地窥视框架的本质。GeeRPC 的目的是以最少的代码，**实现 RPC 框架中最为重要的部分**，帮助大家理解 RPC 框架在设计时需要考虑什么。代码简洁是第一位的，功能是第二位的。

因此，**GeeRPC 选择从零实现 Go 语言官方的标准库 `net/rpc`**，并在此基础上，新增了协议交换(protocol exchange)、注册中心(registry)、服务发现(service discovery)、负载均衡(load balance)、超时处理(timeout processing)等特性。分七天完成，最终代码约 1000 行。

从上面这句内容：“GeeRPC 选择从零实现 Go 语言官方的标准库 net/rpc”，

# 1 服务端与消息编码

一个典型的 RPC 调用如下：

~~~go
err = client.Call("Arith.Multiply", args, &reply)
~~~

客户端发送的请求包括服务名 `Arith`，对应的服务下的某个方法 `Multiply`，以及发送给这个方法的入参。紧接着的是返回值：reply，以及调用的状态反馈 err。

我们将请求和响应中的参数和返回值抽象为 body，剩余的信息放在 header 中，那么就可以抽象出数据结构 Header：

~~~go
type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}
~~~

上面说的 Header 和 Body 部分就是对于一个 HTTP 通信来说的，将一个消息划分为相同的结构。ServiceMethod 是服务名和方法名，通常与 Go 语言中的结构体和方法相映射。Seq 是请求的序列号，也可以认为是某个请求的 ID，用来区分不同的请求。

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

GobCodec 作为 Codec 的一种，需要实现 4 种方法：

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

客户端与服务端的通信需要**协商一些内容**，例如 HTTP 报文，**分为 HEADER 和 Body 部分**，body 的格式和长度通过 HEADER 中**的 Content-Type 和 Content-Length 指定**，服务端通过解析 HEADER 就能够知道如何从 body 中读取需要的信息。对于 RPC 协议来说，这部分协商是需要自主设计的。

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

在一次连接中，Option 固定在报文的最开始，Header 和 Body 可以有多个，即报文可能是这样的：

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
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis) // net.Listener 从哪里来？
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

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
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
		// 处理请求
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
	// | Option | Header1 | Body1 | Header2 | Body2 |...
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	// 通过 cc.ReadBody 修改 req.argv 的值
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println(req.h, req.argv.Elem())

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
			Seq:           uint64(i),
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

Client 在发出请求时，需要在消息的头部添加 Option 内容，但对于 Server 来说，写入的反馈就不需要 Option 内容了。

# 2 支持并发与异步的客户端

在上一节内容中，主要是实现了服务端程序，也就是说，客户端能够发起网络请求，并能获取到 Server 返回的响应。

那本节内容实际上就是实现的是 net/rpc 标准库的 Client 的基本功能：**发出请求**和**接收反馈**。也就是说，经过本节内容，就可以实现大致和 net/rpc 相同的功能。

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











# 3 服务注册







# 4 超时处理







# 5 支持 HTTP 协议







# 6 负载均衡







# 7 服务发现与注册中心





