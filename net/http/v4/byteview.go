package v4

type ByteView struct {
	b []byte
}

// Len 实现 lru.go 中的 Value 接口
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 做了一次深拷贝，防止反馈给用户后修改缓存值（Update入口做统一管理）
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(src []byte) []byte {
	//FIXME dest并没有初始化 make([]byte, len(v.b))，如果没有初始化会出现问题
	var dest []byte
	dest = make([]byte, len(src))
	copy(dest, src)
	return dest
}
