> 要把有限的精力用在自己认为最有价值的事情上，不要被负面情绪裹挟浪费时间。纠结于一些小事会被泥沼缠住脚；甩甩头不理，前方海阔天空。

Go 的 bytes 包，其功能是处理 `[]byte` 类型实例的相关函数；其功能和 strings 包类似，可以做类比（analogous）。

`func Compare(a, b []byte) int`：其表达的含义可以通过下面的代码说明

~~~go
func Compare() {
	// Interpret Compare's result by comparing it to zero.
	var a, b []byte
	if bytes.Compare(a, b) < 0 {
		// a less b
	}
	if bytes.Compare(a, b) <= 0 {
		// a less or equal b
	}
	if bytes.Compare(a, b) > 0 {
		// a greater b
	}
	if bytes.Compare(a, b) >= 0 {
		// a greater or equal b
	}

	// Prefer Equal to Compare for equality comparisons.
	if bytes.Equal(a, b) {
		// a equal b
	}
	if !bytes.Equal(a, b) {
		// a not equal b
	}
}
~~~

[]byte 类型值之间的比较，可以看做是**基于字节值的比较**，毕竟 []byte 中存放的是**字节序列**。

~~~go
package gobytes

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
)

func TestCompare(t *testing.T) {
	// Binary search to find a matching byte slice.
	var needle []byte = []byte("d")

	var haystack [][]byte // Assume sorted
	haystack = append(haystack, []byte("abc"))
	haystack = append(haystack, []byte("abd"))
	haystack = append(haystack, []byte("d"))
	haystack = append(haystack, []byte("efg"))

	i := sort.Search(len(haystack), func(i int) bool {
		// Return haystack[i] >= needle.
		return bytes.Compare(haystack[i], needle) >= 0
	})
	if i < len(haystack) && bytes.Equal(haystack[i], needle) {
		fmt.Println("Found it! index:", i)
	}
}
Found it! index: 2
~~~

`func Count(s, sep []byte) int`：在 s 中寻找和 sep 相同的字节序列，输出其个数。如果 sep 是空串，其结果是 s 中 `UTF-8-encoded code points` 个数 + 1

~~~go
func TestCount(t *testing.T) {
	value := "中国\xE4\xB8\xAD"
	fmt.Println(bytes.Count([]byte(value), nil))
	for index, ele := range []rune(value) {
		fmt.Printf("rune %d#, %q\n", index, ele)
	}

	Count("中国人来自中国", "中国")
	Count("中国人来自中国", "中国人")
	Count("众里寻他千百度，慕容回首，那人却在灯火阑珊处！", "百度")

	Count("five", "ve")
	Count("five", "v")
	Count("five", "")
}

func Count(s, sep string) {
	fmt.Println(bytes.Count([]byte(s), []byte(sep)))
}
4
rune 0#, '中'
rune 1#, '国'
rune 2#, '中'
2
1
1
1
1
5
~~~

`func Fields(s []byte) [][]byte`：将 s 视作为 `UTF-8-encoded code points` 序列，将 s 按照 unicode.IsSpace 进行拆分，其中 sep 可以是一个或者多个 unicode.IsSpace 字符。如果目标 s 中只包含有 unicode.IsSpace，结果返回的是一个空的切片。

~~~go
func TestFields(t *testing.T) {
	Fields("  foo bar  baz   ")
	Fields("  foo bar baz")
	Fields("\tfoo \rbar    \nbaz")
	Fields("\t\r\n")
}

func Fields(s string) {
	result := bytes.Fields([]byte(s))
	fmt.Printf("%T, Fields are %q\n", result, result)
}
[][]uint8, Fields are ["foo" "bar" "baz"]
[][]uint8, Fields are ["foo" "bar" "baz"]
[][]uint8, Fields are ["foo" "bar" "baz"]
[][]uint8, Fields are []
~~~

`func FieldsFunc(s []byte, f func(rune) bool) [][]byte`：将 s 视作为 `UTF-8-encoded code points` 序列，将 s 按照 f 的要求进行拆分。如果 f 返回的是 true 就表示当前的 rune 是一个拆分字符，反之则不是。如果 s 的字符内容均满足 f，或者 s 为空，则返回空的切片。

~~~go
func TestFieldsFunc(t *testing.T) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Printf("Fields are %q\n", bytes.FieldsFunc([]byte("  foo1;bar2,baz3..."), f))
}
Fields are ["foo1" "bar2" "baz3"]
~~~

`func HasPrefix(s, prefix []byte) bool`：判断 s 是否以 prefix 开头。很特别的，如果 prefix 为 nil 或者为空，返回结果默认就是 true。

~~~go
func TestHasPrefix(t *testing.T) {
	fmt.Println(bytes.HasPrefix([]byte("Michoi"), []byte("mi")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), []byte("中")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), []byte("")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), nil))

	fmt.Println(bytes.HasPrefix([]byte(""), nil))
	fmt.Println(bytes.HasPrefix(nil, nil))
}
false
true
true
true
true
true
~~~

`func HasSuffix(s, suffix []byte) bool`：判断 s 是否以 suffix 结尾

~~~go
func TestHasSuffix(t *testing.T) {
	fmt.Println(bytes.HasSuffix([]byte("Michoi"), []byte("oi")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), []byte("中")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), []byte("")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), nil))

	fmt.Println(bytes.HasSuffix([]byte(""), nil))
	fmt.Println(bytes.HasSuffix(nil, nil))
}
true
false
true
true
true
true
~~~

`func Join(s [][]byte, sep []byte) []byte`：将 s 中的 []byte，使用 sep 连接起来

~~~go
func TestJoin(t *testing.T) {
	var s [][]byte
	s = append(s, []byte("M"))
	s = append(s, []byte("i"))
	s = append(s, []byte("c"))
	s = append(s, []byte("h"))
	s = append(s, []byte("o"))
	s = append(s, []byte("i"))
	fmt.Printf("%q\n", bytes.Join(s, []byte("-")))
}
"M-i-c-h-o-i"
~~~

`func Map(mapping func(r rune) rune, s []byte) []byte`：将 s 按照指定映射关系做映射，如果 mapping 的结果为负数，则丢弃该 rune

~~~go
func TestMap(t *testing.T) {
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Printf("%s", bytes.Map(rot13, []byte("'Twas brillig and the slithy gopher...")))
}
'Gjnf oevyyvt naq gur fyvgul tbcure...
~~~

再比如：

~~~go
func TestMap(t *testing.T) {
	rot13 := func(r rune) rune {
		// F 之后的字符输出正值，F 之前的字符输出负值
		switch {
		case r >= 'F' && r <= 'Z':
			return 'A'
		case r >= 'f' && r <= 'z':
			return 'a'
		case r < 'F':
			return -1
		case r < 'f':
			return -1
		}

		return r
	}
	fmt.Printf("%s", bytes.Map(rot13, []byte("'agAG...")))
}
aA
~~~

`func Repeat(b []byte, count int) []byte`：输出 count 次 b 的内容。如果 count 为负数，或者其最终的结果造成内存溢出，会**导致 panic！**

~~~go
func TestRepeat(t *testing.T) {
	fmt.Printf("ba%s", bytes.Repeat([]byte("na"), -2))
}
banana
~~~

`func ReplaceAll(s, old, new []byte) []byte`：在 s 中，使用 new 替换所有的 old。如果 old 是空的，会从 s 的开始直到其结尾，都添加上一个 new 实例。比如 `HiLink` 一共是 6 个 rune 字符，最终替换上了 7 个 new 字符。

~~~go
func TestReplaceAll(t *testing.T) {
	s := "oink oink oink"
	old := "oink"
	new := "moo"
	ReplaceAll(s, old, new)

	s = "HiLink"
	old = ""
	new = "-"
	ReplaceAll(s, old, new)
}

func ReplaceAll(s, old, new string) {
	fmt.Printf("%s\n", bytes.ReplaceAll([]byte(s), []byte(old), []byte(new)))
}
moo moo moo
-H-i-L-i-n-k-
~~~

`func Replace(s, old, new []byte, n int) []byte`：和 ReplaceAll 类似，不同之处是限制了替换次数。如果 count 值为负数，则其调用的结果和 ReplaceAll 相同。

~~~go
func TestReplaceCount(t *testing.T) {
	s := "oink oink oink"
	old := "k"
	new := "ky"
	count := 2
	Replace(s, old, new, count)

	old = "oink"
	new = "moo"
	count = -1
	Replace(s, old, new, count)
}

func Replace(s, old, new string, count int) {
	fmt.Printf("%s\n", bytes.Replace([]byte(s), []byte(old), []byte(new), count))
}
oinky oinky oink
moo moo moo
~~~

`func Runes(s []byte) []rune`：将 s 视作为 `UTF-8-encoded code points`，函数返回 `Unicode code points` 切片

~~~go

func TestRunes(t *testing.T) {
	s := "中Michoi国\xcc, \u4E2D"

	runes := bytes.Runes([]byte(s))
	for i, ele := range runes {
		fmt.Printf("%d# %q; ", i, ele)
	}
	fmt.Println()
}
0# '中'; 1# 'M'; 2# 'i'; 3# 'c'; 4# 'h'; 5# 'o'; 6# 'i'; 7# '国'; 8# '�'; 9# ','; 10# ' '; 11# '中';
~~~

`func Title(s []byte) []byte`：将字符串转化为其文本标题格式，比如单次的首字母大写。

~~~go
func TestTitle(t *testing.T) {
	fmt.Printf("%s", bytes.Title([]byte("her royal highness")))
}
Her Royal Highness
~~~

