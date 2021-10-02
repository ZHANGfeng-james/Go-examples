package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            // 扩充虚拟节点
	keys     []int          // 所有节点根据 hash 值排序（0~2^31-1）
	hashMap  map[int]string // 节点 hash 值和节点名称的映射
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 向 Map 中新增节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 一个缓存节点扩展成 m.replicas 个节点，相当于增加了虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key // key 可认为是某个节点，也就是对一个的某个缓存系统
		}
	}
	// 让所有节点按照 hash 的大小依次顺序排列在整个 0 ~ 2^31-1 组成的环上
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return "" // 表示没有匹配到任何缓存系统节点
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// circle index
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
