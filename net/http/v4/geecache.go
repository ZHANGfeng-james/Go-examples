package v4

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache //FIXME 此处为什么不能是 *cache？什么时候使用指针，什么时候使用普通类型？

	picker PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) //FIXME 此处为什么存的是 *Group？
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error) // 接口型函数，实现了Getter接口

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}

	if _, exist := groups[name]; exist {
		panic("group " + name + " exists")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	// defer mu.RUnlock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) RegistePeers(picker PeerPicker) {
	if g.picker != nil {
		panic("RegistePeers called more than once")
	}
	g.picker = picker
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.picker != nil {
		if peer, ok := g.picker.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer ", err)
		}
	}
	// 调用 getter，用户自定义获取数据方式
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	cache, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: cloneBytes(cache)}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用 getter 从数据源获取数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
