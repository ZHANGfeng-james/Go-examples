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

	v := strings.SplitN(url[len("/standard/"):], "/", 2)
	request := &protopb.Response{
		Name: v[1],
		Age:  18,
	}
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
