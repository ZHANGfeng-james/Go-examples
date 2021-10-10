package discover

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota // 随机选择
	RoundRobinSelect                   // 轮训
)

type Discover interface {
	GetAll() ([]string, error)
	Update([]string) error
	Get(mode SelectMode) (string, error)
	Refresh() error // refresh from remote registry
}

type MultiServersDiscovery struct {
	servers []string

	mu    sync.RWMutex
	rand  *rand.Rand
	index int
}

func NewMultiServersDiscovery(servers []string) *MultiServersDiscovery {
	instance := &MultiServersDiscovery{
		servers: servers,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	instance.index = instance.rand.Intn(math.MaxInt32 - 1) // 避免每个实例都从 0 开始
	return instance
}

func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	len := len(d.servers)
	if len == 0 {
		return "", errors.New("rpc discovery, no avaiable server")
	}

	switch mode {
	case RandomSelect:
		// 随机选取
		n := d.rand.Intn(len)
		return d.servers[n], nil
	case RoundRobinSelect:
		// 轮询方式选取
		n := d.index % len // servers could be update, so mode len to ensure safety
		d.index = (d.index + 1) % len
		// 特别是 servers 增加的时候，可能会出现访问越界
		return d.servers[n], nil
	default:
		return "", errors.New("rpc discovery, unknown select mode")
	}
}

func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	//FIXME d.mu.Unlock()
	defer d.mu.RUnlock()

	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)

	return servers, nil
}
