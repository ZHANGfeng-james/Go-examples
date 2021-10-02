

Go 中内置函数 copy：

~~~go
// The copy built-in function copies elements from a source slice into a
// destination slice. (As a special case, it also will copy bytes from a
// string to a slice of bytes.) The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
func copy(dst, src []Type) int

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int
~~~

copy 函数中的入参类型是 Type，其真实类型可以是任何的 Go 类型，并不局限是 int 类型。

copy 函数的注释中有这么一段：函数返回值表示已拷贝的元素个数，其值等于 len(dst) 和 len(src) 较小的值。对应的使用：

~~~go
func cloneBytes(src []byte) []byte {
	//FIXME dest并没有初始化 make([]byte, len(v.b))，如果没有初始化会出现问题
	var dest []byte
    dest = make([]byte, len(src)) // 替换成 dest = make([]byte, 0)
	copy(dest, src)
	return dest
}
~~~

替换成 `dest = make([]byte, 0)`，将**不会发生任何拷贝动作**，或者是仅声明 dest 变量。**实际拷贝的字节数，是和 len(dst) 和 len(src) 的较小值相等的**。