

完成 2 件事情：


1. 首先我们实现最简单的基于 gRPC 的调用：仅仅使用的是 protobuf 数据通信格式，基座是 HTTP 通信；
2. 标准的 gRPC 通信案例；

# 1 基于 protobuf 数据格式的 HTTP 通信

我们实现最简单的基于 gRPC 的调用：仅仅使用的是 protobuf 数据通信格式，基座是 HTTP 通信。其潜在的实现就是，Client 和 Server 端的通信仍然是简单的 HTTP，但其数据格式是不再基于字符，而是 protoc 编译得到的**二进制数据**。

定义数据传输格式：

~~~go
syntax="proto3";

package protopb;

option go_package="../protopb";

message Request{
    
}

message Response{
    string name = 1;
    int32 age = 2;
}

// no service
~~~

在 `.proto` 文件所在的目录下执行 `protoc --go_out=. *.proto` 可生成对应的 `.pb.go` 文件。

接下来定义客户端代码：

~~~go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/go-examples-with-tests/net/rpc/v3/protopb"
	"google.golang.org/protobuf/proto"
)

var baseURL = "http://localhost:9090/standard/rpc"

func main() {
	fmt.Println("client run...")

	query := "Katyusha"
	url := fmt.Sprintf("%v/%v", baseURL, url.QueryEscape(query))
	log.Printf("client request:%s", url)
	// Client Request 的 name 信息是通过 URL 传到 Server 端
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http.Get %s; error:%s", url, err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body error:%v", err)
		return
	}

	response := &protopb.Response{}
	if err := proto.Unmarshal(bytes, response); err != nil {
		log.Printf("decode request protobuf:%v", err)
		return
	}

	log.Printf("response: Name:%s; Age:%d", response.Name, response.Age)
}
~~~

紧接着是服务端代码：

~~~go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-examples-with-tests/net/rpc/v3/protopb"
	"google.golang.org/protobuf/proto"
)

type server struct{}

func (server *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 从 request 中读取到请求数据
	fmt.Println("server get Reqeust")
	url := r.URL.Path
	log.Printf("Server get client request:%s", url)

    // "http://localhost:9090/standard/rpc/Katyusha" url:/standard/rpc/Katyusha
	v := strings.SplitN(url[len("/standard/"):], "/", 2)
	request := &protopb.Response{
		Name: v[1],
		Age:  18,
	}
    // protobuf 编码
	bytes, err := proto.Marshal(request)
	if err != nil {
		log.Printf("decode request protobuf:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}

func main() {
	fmt.Println("server...")
	// 启动 http 服务
	http.ListenAndServe(":9090", &server{})
}
~~~

整个通信的基础依然是 HTTP，但是通信数据是 protobuf 格式的。

# 2 gRPC

使用 proto 文件定义服务以及 message：

~~~go
syntax="proto3";

option go_package="github.com/go-examples-with-tests/net/rpc/v4/protopb";

package protopb;

message HelloRequest{
    string name = 1;
}

message HelloReply{
    string message = 1;
}

// The greeting service definition
service Greeter {
    // send a greeting
    rpc SayHello(HelloRequest) returns (HelloReply){}
}
~~~

其中 `go_package` 是导包路径，对应的 protoc 编译指令：`protoc -I. --go_out=plugins=grpc:$GOPATH/src helloworld.proto` 对应会在 `$GOPATH/src` 目录下生成 `helloworld.pb.go` 文件。

特别注意，此处**我们定义了一个 service，也就是一个服务**。

服务端代码实现：

~~~go
package main

import (
	"context"
	"log"
	"net"

	"github.com/go-examples-with-tests/net/rpc/v4/protopb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	protopb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *protopb.HelloRequest) (*protopb.HelloReply, error) {
	log.Printf("Received: %s", req.GetName())
	return &protopb.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen:%s", port)
	}

	// 创建一个 gRPC Server 实例
	s := grpc.NewServer()
    // *grpc.Server 和 service 关联
	protopb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("falied to serve:%v", err)
	}
}
~~~

在 .proto 文件中定义的 service，在 Server 端的用处：在创建 *grpc.Server 之后，需要让这个实例和 service 关联起来，也就是 Client 的请求会让这个 *grpc.Server 处理。

客户端代码实现：

~~~go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v4/protopb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	// 监听指定 port，获得一个 *grpc.ClientConn 实例
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial addr: %s error", port)
	}
	defer conn.Close()

	// 使用这个 *grpc.ClientConn 实例，创建指定的 GreeterClient 实例
	c := protopb.NewGreeterClient(conn)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 调用 GreeterClient 实例的方法，和在本地调用方法是一样的，这就是 RPC 带来的便捷
	reply, err := c.SayHello(ctx, &protopb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet:%v", err)
	}
	log.Printf("Get reply:%s", reply.GetMessage())
}
~~~

`c := protopb.NewGreeterClient(conn)` 相当于创建了客户端 Stub。在生成的 `.pb.go` 文件中，通过 `NewGreeterClient` 函数就能创建和 server 对应的 Client 实例，通过这个实例就能请求对应的 RPC 方法。

