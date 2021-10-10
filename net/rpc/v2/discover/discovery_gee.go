package discover

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-examples-with-tests/database/v1/log"
)

type GeeRegistryDiscovery struct {
	*MultiServersDiscovery
	registry string // Registry中心地址

	timeout          time.Duration // 更新服务列表时间间隔
	lastUpdateServer time.Time     // 最后一次更新服务列表的时间戳
}

const defaultUpdateTimeout = 10 * time.Second

func NewGeeRegistryDiscovery(registryAddr string, duration time.Duration) *GeeRegistryDiscovery {
	if duration == 0 {
		duration = defaultUpdateTimeout
	}

	return &GeeRegistryDiscovery{
		MultiServersDiscovery: NewMultiServersDiscovery(make([]string, 0)),
		timeout:               duration,
		registry:              registryAddr,
	}
}

func (d *GeeRegistryDiscovery) Refresh() error { // refresh from remote registry
	// GET 请求和 Registry 通信，获取所有可用服务列表，并更新本地缓存
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdateServer.Add(d.timeout).After(time.Now()) {
		return nil
	}

	log.Infof("rpc registry: refresh servers from registry:%s", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Infof("rpc registry refresh err:", err)
		return err
	}

	servers := strings.Split(resp.Header.Get("X-Geerpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, s := range servers {
		if strings.TrimSpace(s) != "" {
			d.servers = append(d.servers, strings.TrimSpace(s))
		}
	}
	d.lastUpdateServer = time.Now()
	return nil
}

func (d *GeeRegistryDiscovery) Update(servers []string) error {
	//FIXME 功能未知！
	d.mu.Lock()
	defer d.mu.Unlock()

	d.servers = servers
	d.lastUpdateServer = time.Now()
	return nil
}

func (d *GeeRegistryDiscovery) GetAll() ([]string, error) {
	// 判断是否超时，若未超时，则直接返回本地缓存；否则，执行 Refresh 后返回本地缓存
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}

func (d *GeeRegistryDiscovery) Get(mode SelectMode) (string, error) {
	// 同 GetAll
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}
