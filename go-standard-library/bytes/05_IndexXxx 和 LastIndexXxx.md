Go 的 bytes 包，其功能是处理 `[]byte` 类型实例的相关函数；其功能和 strings 包类似，可以做类比（analogous）。

本篇我们来介绍 bytes 包下有关 IndexXxx 和 LastIndexXxx 的函数。

# 1 IndexXxx

总的来说，IndexXxx 函数用于查找 s 中的指定内容，并返回首字节的位置索引值。

`func Index(s, sep []byte) int`：在 s 中搜索指定的 sep 内容（需要完全匹配），并返回（**首次匹配时**）首字节的位置索引。如果没找到，则返回 -1。

~~~go
func TestIndexNormal(t *testing.T) {
	index := bytes.Index([]byte("Michoi"), []byte("hoi"))
	fmt.Printf("index = %d\n", index)

	index = bytes.Index([]byte("chicken"), []byte("Michoi"))
	fmt.Printf("index = %d\n", index)
}
index = 3
index = -1
~~~

`func IndexAny(s []byte, chars string) int`：把 s 当做是 ` UTF-8-encoded Unicode code points` 的序列，方法返回的是 chars 中的**任意一个（这就是 Any 表达的含义）** ` Unicode code points` 首次出现在 s 中的位置索引值。如果 chars 为空，或者 chars 中的 Unicode 码值没有出现在 s 中，则返回 -1。

~~~go
func TestIndexAny(t *testing.T) {
	index := bytes.IndexAny([]byte("Mi国cho中i"), "中国")
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexAny([]byte("chicken"), "aeiouy")
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexAny([]byte("crwth"), "aeiouy")
	fmt.Printf("index = %d\n", index)
}

index = 2
index = 2
index = -1
~~~

`func IndexByte(b []byte, c byte) int`：和上述函数不同的是，IndexByte 匹配的是一个字节 c，并返回首次出现在 b 中的字节索引值。

~~~go
func TestIndexByte(t *testing.T) {
	index := bytes.IndexByte([]byte("Michoi"), byte('i'))
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexByte([]byte("Michoi"), byte('a'))
	fmt.Printf("index = %d\n", index)

	character := 'i'
	fmt.Printf("%T, %q\n", character, character)
}
index = 1
index = -1
int32, 'i'
~~~

`func IndexRune(s []byte, r rune) int`：这个函数需要和 IndexByte 区别开，IndexRune 的入参是 rune 类型，即为一个 Unicode 码的 UTF-8 值。该函数返回的是首个在 s 中出现的 r 的位置索引值。如果 r 是 utf8.RuneError 则返回 s 中首个非 UTF-8 字节位置索引。

~~~go
func TestIndexRune(t *testing.T) {
	index := bytes.IndexRune([]byte{0x01, 0x02, 0xff, 0xee, 0x2f, 0x33}, utf8.RuneError)
	fmt.Printf("index = %d\n", index)
}
index = 2
~~~

`func IndexFunc(s []byte, f func(r rune) bool) int`：把 s 视作为 ` UTF-8-encoded Unicode code points` 的序列，函数返回首个满足 f 的 rune 所在 s 中的位置索引值。

~~~go
func TestIndexFunc(t *testing.T) {
	satisfyFunc := func(r rune) bool {
		fmt.Printf("%q\n", r)
		return unicode.Is(unicode.Han, r)
	}

	index := bytes.IndexFunc([]byte("MichoiЛ中国"), satisfyFunc)
	fmt.Printf("index = %d\n", index)
}
'M'
'i'
'c'
'h'
'o'
'i'
'Л'
'中'
index = 8
~~~

# 2 LastIndexXxx

`func LastIndex(s, sep []byte) int`

`func LastIndexAny(s []byte, chars string) int`

`func LastIndexByte(s []byte, c byte) int`

`func LastIndexFunc(s []byte, f func(r rune) bool) int`

上述这些 LastIndexXxx 方法和对应的 IndexXxx 方法类似，不同之处是前者获取的最后一次匹配成功的位置索引：

~~~go
func TestIndexLast(t *testing.T) {
	index := bytes.LastIndex([]byte("Michoic"), []byte("ic"))
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndex([]byte("chicken"), []byte("Michoi"))
	fmt.Printf("index = %d\n", index)

	index = bytes.LastIndexAny([]byte("Mi国cho中i"), "中国")
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexAny([]byte("chicken"), "aeiouy")
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexAny([]byte("crwth"), "aeiouy")
	fmt.Printf("index = %d\n", index)

	index = bytes.LastIndexByte([]byte("Michoi"), byte('i'))
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexByte([]byte("Michoi"), byte('a'))
	fmt.Printf("index = %d\n", index)

	satisfyFunc := func(r rune) bool {
		fmt.Printf("%q\n", r)
		return unicode.Is(unicode.Han, r)
	}
	index = bytes.LastIndexFunc([]byte("MichoiЛ中国"), satisfyFunc)
	fmt.Printf("index = %d\n", index)
}
index = 5
index = -1
index = 8
index = 5
index = -1
index = 5
index = -1
'国'
index = 11
~~~

