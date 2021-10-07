Go 语言标准库 net/rpc 提供了通过**网络**或**其他 I/O 连接**对**一个对象导出方法**的**访问**。

具体的实现是：服务端注册一个**对象**，使它**作为一个服务**被暴露，**服务的名字**是**该对象的类型名**。注册之后，对象的导出方法就可以被远程访问。服务端可以注册多个不同类型的对象，但注册具有相同类型的多个对象是错误的。

只有满足如下标准的方法才能用于远程访问：

1. 方法所属的类型是可导出的，比如下述类型 `T`；
2. 方法是可导出的；
3. 方法有 2 个参数，都是可导出类型或内建类型；
4. 方法的第二个参数是指针；
5. 方法只有一个 error 接口类型的返回值。

符合上述标准的函数类型是这样的：

~~~go
func (t *T)MethodName(argType T1, replyType *T2) error
~~~

其中 T1 和 T2 是可以被 encoding/gob 包序列化的，上述限制条件即使是使用不同的编解码器仍然适用。

方法的第一个参数代表调用者提供的参数；第二个参数表示返回给调用者的参数。方法的返回值，如果非 nil，将被作为**字符串**回传，在客户端看来就和 errors.New 创建的一样。如果返回了错误，回复的参数不会被发送到客户端。

服务端可能会在单个连接上调用 ServeConn **管理请求**，更典型的是，它会创建一个网络监听器然后调用 Accept，或者对于 HTTP 监听器，调用 HandleHTTP 和 http.Serve。

想要使用服务的客户端会创建一个链接，然后用该连接调用 NewClient。更方便的函数 Dial 会在一个原始的连接上依次执行上述 2 个步骤。生成的 Client 类型值有两个方法：Call 和 Go，它们的参数为要调用的服务和方法、一个包含参数的指针、一个用于接收结果的指针。Call 方法会等待远端调用完成**（同步阻塞式）**，而 Go 方法**异步地**发送调用请求并使用返回的 Call 结构体类型的 Done 通道字段传递完成信号。

除非显式地设置编解码器，net/rpc 包**默认**使用 encoding/gob 包来传输数据。

我先来看看标准库 net/rpc 都有什么内容？

~~~go
ant@MacBook-Pro Go-examples-with-tests % go doc net/rpc |grep "^func"
func Accept(lis net.Listener)
func HandleHTTP()
func Register(rcvr interface{}) error
func RegisterName(name string, rcvr interface{}) error
func ServeCodec(codec ServerCodec)
func ServeConn(conn io.ReadWriteCloser)
func ServeRequest(codec ServerCodec) error
~~~

重要的结构体类型：

~~~go
ant@MacBook-Pro Go-examples-with-tests % go doc net/rpc |grep "^type" |grep "struct"
type Call struct{ ... }
type Client struct{ ... }
type Request struct{ ... }
type Response struct{ ... }
type Server struct{ ... }

ant@MacBook-Pro Go-examples-with-tests % go doc net/rpc |grep "^type" |grep "interface"
type ClientCodec interface{ ... }
type ServerCodec interface{ ... }
~~~

标准库中的 net/rpc 和 gRPC 有什么关系？

从 Go 标准库 net/rpc 的描述来看，这个库实现了基本的 RPC 功能！

下面是一个简单的应用例子：

~~~go
package v1

import (
	"errors"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Que, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quotient *Quotient) error {
	if args.B == 0 {
		// string --> errors.New
		return errors.New("divided by zero")
	}
	quotient.Que = args.A / args.B
	quotient.Rem = args.A % args.B
	time.Sleep(2 * time.Second) // 模拟 Server 端处理延时
	return nil
}
~~~

Serve 端：

~~~go
package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

func main() {
	arith := new(v1.Arith)
	rpc.Register(arith)

	rpc.HandleHTTP() // HTTP 监听器

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(listener, nil) // 调用的是 http 包下的函数
}
~~~

客户端请求服务：**同步方式**

~~~go
package main

import (
	"fmt"
	"log"
	"net/rpc"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

const serverAddr = ""

func main() {
	client, err := rpc.DialHTTP("tcp", serverAddr+":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}

	args := &v1.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith err:", err)
	}
	fmt.Printf("Arith.Multiply: %d*%d=%d", args.A, args.B, reply)
}
~~~

客户端请求服务：异步方式

~~~go
package main

import (
	"context"
	"fmt"
	"log"
	"net/rpc"
	"time"

	v1 "github.com/go-examples-with-tests/net/rpc/v1"
)

const serverAddr = ""

func main() {
	client, err := rpc.DialHTTP("tcp", serverAddr+":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}

    // 增加超时处理机制
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := &v1.Args{A: 7, B: 8}
	quotient := new(v1.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)

	select {
	case <-timeout.Done():
		fmt.Println("timeout")
	case <-divCall.Done:
		fmt.Printf("Arith.Divide: %d / %d = %d.%d\n", args.A, args.B, quotient.Que, quotient.Rem)
	}

	fmt.Println("OVER")
}
~~~

服务端的实现应为客户端提供简单、类型安全的包装。

标准库的 net/rpc 默认使用的是 encoding/gob 编解码，支持 TCP 和 HTTP 数据传输方式。由于其他语言不支持 gob 编解码方式，由此使用 Go 语言标准的 net/rpc 实现的服务不支持跨语言通信。
