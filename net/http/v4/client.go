package v4

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/go-examples-with-tests/net/http/v4/cachepb"
	"google.golang.org/protobuf/proto"
)

// httpGetter send HTTP request to specific Server, and get the cache value
type httpGetter struct {
	baseURL string
}

func (getter *httpGetter) Get(in *cachepb.Request, out *cachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v",
		getter.baseURL, url.QueryEscape(in.Group), url.QueryEscape(in.Key))

	// http://localhost:8003/_geecache/scores/Katyusha
	log.Printf("httpGetter send request to: %v", u)
	// 依据指定的 URL 发出请求，等待 Server 响应并回传 cache value
	response, err := http.Get(u)
	if err != nil {
		log.Println("client Get: " + err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", response.Status)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading body error:%v", err)
	}

	// 通信数据的格式是 protobuf，解码
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
