package v4

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// httpGetter send HTTP request to specific Server, and get the cache value
type httpGetter struct {
	baseURL string
}

func (getter *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v",
		getter.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	log.Printf("httpGetter send request to: %v", u)
	// 依据指定的 URL 发出请求，等待 Server 响应并回传 cache value
	response, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", response.Status)
	}

	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body error:%v", err)
	}
	return res, nil
}
