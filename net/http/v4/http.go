package v4

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-examples-with-tests/net/http/v4/consistenthash"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

// 承载节点 HTTP 通信的核心数据结构，是分布式缓存系统的一个节点
// 作为通信的服务端，也就是接收客户端的 HTTP 缓存请求
type HTTPPool struct {
	self     string
	basePath string // API前缀，默认是 defaultBasePath

	mu          sync.Mutex
	peers       *consistenthash.Map    // hash(key) 选择 http://10.0.0.2:8008
	httpGetters map[string]*httpGetter // http://10.0.0.2:8008 --> *httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP serve client HTTP request, and response the cache value
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	// get cache value
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
	w.Write([]byte("\r\n"))
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 在本地为远端分布式节点建立模型，也就是将自身看作是 Client，从远端节点获取 cache value
	p.peers = consistenthash.New(defaultReplicas, nil)
	// add all peers
	log.Printf("HTTPPool add peers: %v", peers)
	// http://localhost:8001, http://localhost:8002, http://localhost:8003...
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter)

	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	peer := p.peers.Get(key)
	log.Printf("get peer in peers: %v, p.self: %v", peer, p.self)
	if peer != "" && peer != p.self {
		// peer != p.self 时，发起远端 HTTP 请求；如果是，则直接在本地请求
		return p.httpGetters[peer], true
	}
	return nil, false
}
