package lru

import "container/list"

type Cache struct {
	maxBytes  int64                         // Cache最大能容纳字节数
	nbytes    int64                         // 当前已装载容量
	ll        *list.List                    // Cache数据结构中双端链表
	cache     map[string]*list.Element      // key-value map，用于找到指定 key 对应的 *list.Element
	onRemoved func(key string, value Value) // Callback 机制
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onRemoved func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onRemoved: onRemoved,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // ele 类型：*list.Element
		// update the exists Element
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) // ele.Value 类型：*entry
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// add a new Element
		newNode := c.ll.PushFront(&entry{key, value})
		c.cache[key] = newNode
		c.nbytes += int64(len(key) + value.Len())
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// Cache最前端是最近一次访问的节点，末尾是最早访问的节点
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) // ele.Value 类型：*entry
		return kv.value, true    // 如果直接写成 return 返回的 ok 值为 false
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.ll.Remove(ele)
		c.nbytes -= (int64(len(kv.key)) + int64(kv.value.Len()))
		if c.onRemoved != nil {
			c.onRemoved(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int64 {
	return int64(c.ll.Len())
}
