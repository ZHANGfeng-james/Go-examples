package v4

import (
	"sync"

	"github.com/go-examples-with-tests/net/http/v4/lru"
)

type cache struct {
	lock       sync.Mutex // 无需初始化，直接就可使用
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, view ByteView) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil) // lru.Cache函数返回的是 *lru.Cache 类型值
	}
	c.lru.Add(key, view)
}

func (c *cache) get(key string) (v ByteView, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.lru == nil {
		return
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), true // 返回了 ByteView 后，不会被改变吗？
	}
	return
}
