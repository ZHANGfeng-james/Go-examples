package singleflight

import "sync"

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

type call struct {
	wg sync.WaitGroup

	// 用于保存相同 key，通过 HTTP 获取到的缓存结果
	val interface{}
	err error
}

func (g *Group) Do(key string, fun func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 延迟加载，懒加载
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock() // 此处必须加上解锁！
		c.wg.Wait()   // 如果没有任何 goroutine 执行了 c.wg.Add(1)，会有什么影响？
		return c.val, c.err
	}

	c := new(call) // *call 类型
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fun()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
