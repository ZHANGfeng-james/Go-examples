package registry

import (
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-examples-with-tests/database/v1/log"
)

// 作为一个 registry，需要装载哪些字段，才能实现这个功能模型
type GeeRegistry struct {
	servers map[string]*ServerItem
	timeout time.Duration // 服务需删除的时间间隔，用于保持 registry 干净、可用
	mu      sync.Mutex
}

type ServerItem struct {
	Addr           string    // 服务实例 rpcAddr
	aliveStartTime time.Time // 服务实例存活开始时间
}

const (
	defaultRegistryPath = "/_geerpc_/registry"
	defaultTimeout      = time.Minute * 5 // 5分钟之内，没有任何心跳包，表示服务不再存活
)

var DefaultGeeRegistry = NewRPCRegistry(defaultTimeout)

func NewRPCRegistry(duration time.Duration) *GeeRegistry {
	return &GeeRegistry{
		servers: make(map[string]*ServerItem),
		timeout: duration,
	}
}

func (registry *GeeRegistry) putServer(rpcAddr string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	server, ok := registry.servers[rpcAddr]
	if !ok {
		registry.servers[rpcAddr] = &ServerItem{Addr: rpcAddr, aliveStartTime: time.Now()}
	} else {
		// if exists, update alive time to keep alive 表示服务还活着！
		server.aliveStartTime = time.Now()
	}
}

func (registry *GeeRegistry) aliveServers() []string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	servers := make([]string, 0)
	for rpcAddr, serverItem := range registry.servers {
		if registry.timeout == 0 || serverItem.aliveStartTime.Add(registry.timeout).After(time.Now()) {
			servers = append(servers, rpcAddr)
		} else {
			delete(registry.servers, rpcAddr)
		}
	}
	sort.Strings(servers)
	log.Infof("registry alive server all, rpcAddr:%v", servers)
	return servers
}

func (registry *GeeRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// keep it simple, server is in req.Header
		w.Header().Set("X-Geerpc-Servers", strings.Join(registry.aliveServers(), ","))
	case http.MethodPost:
		addr := req.Header.Get("X-Geerpc-Server")
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Infof("registry receive heartbeat from:%s", addr)
		registry.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (registry *GeeRegistry) HandleHTTP(registryPath string) {
	http.Handle(registryPath, registry)
}

func HandleHTTP() {
	DefaultGeeRegistry.HandleHTTP(defaultRegistryPath)
}
