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
