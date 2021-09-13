strings 标准库**源代码路径**：`./Go/src/strings/` 目录，其中包含：

1. builder.go、compare.go、reader.go、replace.go、search.go、strings.go
2. builder_test.go、compare_test.go、reader_test.go、replace_test.go、search_test.go、strings_test.go、==example_test.go==、==export_test.go==

strings 包用于操作 UTF-8 编码的字符串！

> 即然操作的是 UTF-8 字符串，那首先要知道 UTF-8 字符编码格式，以及字符串类型。

# 1 strings 导出的方法

`func Compare(a, b string) int`

该函数仅包含的是包字节的对称性，而且该函数通常比 string 类型**内置的比较操作符** `==`、`<`、`>` 更加清楚、快速。

~~~go
// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
//
// Compare is included only for symmetry with package bytes.
// It is usually clearer and always faster to use the built-in
// string comparison operators ==, <, >, and so on.
func Compare(a, b string) int {
	// NOTE(rsc): This function does NOT call the runtime cmpstring function,
	// because we do not want to provide any performance justification for
	// using strings.Compare. Basically no one should use strings.Compare.
	// As the comment above says, it is here only for symmetry with package bytes.
	// If performance is important, the compiler should be changed to recognize
	// the pattern so that all code doing three-way comparisons, not just code
	// using strings.Compare, can benefit.
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return +1
}
~~~

虽然从注释中获取到关于性能的比较，但是其源代码实现中却使用的是内置的 string 类型的比较操作符，根本没区别！

------

`func Contains(s, substr string) bool`

该函数判断 `s` 中是否包含了 `substr`。

特别的示例：`strings.Contains("ss", "")` 其结果为 true；`strings.Contains("", "")` 其结果也为 true。

------

`func ContainsAny(s, chars string) bool`

该函数判断 chars 中的任意 Unicode code point 是否包含在 s 中。

特别的示例：`strings.ContainsAny("ss", "")` 其结果为 false；`strings.ContainsAny("", "")` 其结果也为 false

-----

`func ContainsRune(s string, r rune) bool`

该函数判断 Unicode code point r 是否包含在 s 中。

~~~go
func main() {
	var a rune
	a = 97
	a = 'a'
	fmt.Println(a)

	fmt.Println(strings.ContainsRune("ssss", 's')) // true
	fmt.Println(strings.ContainsRune("a", 97)) // true
}
~~~

----

`func Count(s, substr string) int`

该函数计算出 s 中的 substr 非重叠实例，如果 substring 为空字符串，则返回 s 中 Unicode code point 个数 + 1。

~~~go
package main

func main() {
	fmt.Println(strings.Count("sss", "ss")) // 1
	fmt.Println(strings.Count("fine", "")) // 5
}
~~~

----

`func EqualFold(s, t string) bool`

该函数中将 s 和 t 视为 UTF-8 编码的字符串，判断 s 和 t 在忽略大小写的情况下，是否相等。这是不区分大小写的更通用形式。

------

`func Fields(s string) []string`

该函数将 s 按照**连续的** 1 个或多个空格字符（`unicode.IsSpace(rune)` 判断）拆分，其结果返回的是 `[]string`，或者是空切片（如果 s 中仅包含空格字符）。

~~~go
func main() {
    // ["1" "foo" "bar" "baz" "2"]
	fmt.Printf("Fields are: %q\n", strings.Fields("1   foo bar  baz   2"))
	fmt.Printf("Fields are: %q\n", strings.Fields("   ")) // []
	fmt.Printf("Fields are: %q\n", strings.Fields("")) // []
}
~~~

------

`func FieldsFunc(s string, f func(r rune) bool) []string`

该函数按照满足 `f(r)` 的标准切割 s，并返回 `[]string` 值。如果 s 中的每一个 Unicode code points 都满足 `f(r)` 或者 s 是一个空字符串，将返回一个空的 `[]string`。

~~~go
func main() {
	f := func(c rune) bool {
        // 既不是字母，也不是数字（&&的关系）
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
    // foo1 bar2 baz3
	fmt.Printf("Fields are: %q\n", strings.FieldsFunc("  foo1;bar2,baz3...", f)) 
}
~~~

----

`func HasPrefix(s, prefix string) bool`

该函数判断 s 是否是以 prefix 开头。特别的，如果 prefix 为空字符串，不管 s 是这样的字符串，都返回 true。

~~~go
func main() {
	fmt.Println(strings.HasPrefix("", ""))           // true
	fmt.Println(strings.HasPrefix("gopher", ""))     // true
	fmt.Println(strings.HasPrefix("gopp", " "))      // false
}
~~~

-----

`func HasSuffix(s, suffix string) bool`

该函数判断 s 是否以 suffix 为结尾（后缀）。

~~~go
func main() {
	fmt.Println(strings.HasSuffix("", ""))           // true
	fmt.Println(strings.HasSuffix("gopher", ""))     // true
	fmt.Println(strings.HasSuffix("gopp", " "))      // false
}
~~~

得到的结果和 `strings.HasPrefix()` 相同，特别是 suffix 是空字符串的时候。

----

`func Index(s, substr string) int`

该函数返回 s 中首次出现 substr 的索引位置，如果 s 中不存在 substr，则返回 -1.

~~~go
func main() {
	fmt.Println(strings.Index("", ""))             // 0
	fmt.Println(strings.Index("gopher", ""))       // 0
	fmt.Println(strings.Index("gopp", " "))        // -1

	fmt.Println(strings.Index("中国", "国"))       // 3
}
~~~

如果 substr 为空字符串，则不论 s 是否为空，都将返回 0；对于非 ASCII 的处理，其结果返回的是底层字节数组的位置索引。

`func LastIndex(s, substr string) int`

该函数返回 s 中最后一次出现 substr 的索引位置，如果不存在则返回 -1。

~~~go
func main() {
	fmt.Println(strings.Index("go gopher", "go"))         // 0
	fmt.Println(strings.LastIndex("go gopher", "go"))     // 3
	fmt.Println(strings.LastIndex("go gopher", "rodent")) // -1

	fmt.Println(strings.Index("go gopher", ""))     // 0
	fmt.Println(strings.LastIndex("go gopher", "")) // 9
}
~~~

----

`func IndexAny(s, chars string) int`

该函数返回 chars 中的 Unicode code point 出现在 s 中的首字符位置，如果不存在则返回 -1。

~~~go
func main() {
	fmt.Println(strings.IndexAny("", ""))        // -1
	fmt.Println(strings.IndexAny("", "gopher"))
	fmt.Println(strings.IndexAny(" ", "gopp"))

	fmt.Println(strings.Index("中国", "国"))     // 3
    
     // 2，只要 s 中存在 chars 的存在任意字符即可，体现在方法名上的 Any 概念
	fmt.Println(strings.IndexAny("chicken", "aeiouy"))  
	fmt.Println(strings.IndexAny("crwth", "aeiouy"))    // -1
}
~~~

`func LastIndexAny(s, chars string) int`

该方法和 `IndexAny` 恰好相反，表示最后的含义，如果不存在则返回 -1。

~~~go
func main() {
	fmt.Println(strings.LastIndexAny("", "")) // -1
	fmt.Println(strings.LastIndexAny("", "gopher"))
	fmt.Println(strings.LastIndexAny(" ", "gopp"))

	fmt.Println(strings.LastIndexAny("中国", "国")) // 3

	fmt.Println(strings.LastIndexAny("chicken", "aeiouy"))   // 5，s 中最后能匹配的字符是 e
	fmt.Println(strings.LastIndexAny("crwth", "aeiouy"))     // -1
	fmt.Println(strings.LastIndexAny("go gopher", "go"))     // 4
	fmt.Println(strings.LastIndexAny("go gopher", "rodent")) // 8
	fmt.Println(strings.LastIndexAny("go gopher", "fail"))   // -1
}
~~~

----

`func IndexByte(s string, c byte) int`

该函数特别注意的是其匹配参数值是 `byte`，而不是 string。该函数返回 s 中首次出现 c 的位置索引，如果不存在，则返回 -1。

~~~go
func main() {
	fmt.Println(strings.IndexByte("", ' '))           // -1
	fmt.Println(strings.IndexByte("", 'g'))           // -1
	fmt.Println(strings.IndexByte(" ", ' '))          // 0

	fmt.Println(strings.IndexByte("chicken", 'a'))    // -1
	fmt.Println(strings.IndexByte("crwth", 'w'))      // 2
}
~~~

`func LastIndexByte(s string, c byte) int`

~~~go
func main() {
	fmt.Println(strings.LastIndexByte("Hello, world", 'l'))
	fmt.Println(strings.LastIndexByte("Hello, world", 'o'))
	fmt.Println(strings.LastIndexByte("Hello, world", 'x'))
}
~~~

----

`func IndexFunc(s string, f func(rune) bool) int`

该函数返回的是 s 中首个能够满足 `f(r)` 的 Unicode code point 位置。

~~~go
func main() {
	f := func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	}
	fmt.Println(strings.IndexFunc("Hello, 中国!", f))
	fmt.Println(strings.IndexFunc("Hello, world!", f))
}
~~~

`func LastIndexFunc(s string, f func(r rune) bool) int`

~~~go
func main() {
	fmt.Println(strings.LastIndexFunc("go 123", unicode.IsNumber))
}

// IsNumber reports whether the rune is a number (category N).
func IsNumber(r rune) bool {
	if uint32(r) <= MaxLatin1 {
		return properties[uint8(r)]&pN != 0
	}
	return isExcludingLatin(Number, r)
}
~~~

从 `unicode.IsNumber` 函数签名来看，是符合 `strings.LastIndexFunc` 函数第二个参数的。

-----

`func IndexRune(s string, r rune) int`

该函数返回的是 s 中首个值为 r 的 Unicode code point 位置，如果不存在则返回 -1；如果 r 是 `utf8.RuneError`，则返回第一个无效的 UTF-8 字节序列。

~~~go
func main() {
	fmt.Println(strings.IndexRune("中国人I am Chinese-!", 'C'))
}
~~~

==这个函数的使用是有疑问的！而且为什么没有 `LastIndexRune`？==

----

`func Join(elems []string, sep string) string`

该函数会级联 `elems` 中的 string 并重新构成新的 string，同时 sep 会被放置在 elems 各个元素之间。

~~~go
func main() {
	s := []string{"foo", "zebr", "zoo"}
	fmt.Println(strings.Join(s, " - "))
}
~~~

----

`func Map(mapping func(rune) rune, s string) string`

该方法依据 `mapping` 函数将 s 字符串序列转化为指定字符串。如果转化得到的 rune 是负数，则舍弃该 rune。

~~~go
func main() {
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'z'+13)%26
		}
		return r
	}
	fmt.Println(strings.Map(rot13, "Twas brillig and the slithy gopher..."))
}
~~~

----

`func Repeat(s string, count int) string`

该方法将对 s 做 count 次的重复拷贝，得到新的字符串序列。如果 count 是负数或者 `len(s) * count` 溢出，将会引发 panic。

~~~go
func main() {
	fmt.Println("ba" + strings.Repeat("na", -2))
}
~~~

----

`func ReplaceAll(s, old, new string) string`

该函数将 s 中的 old 替换为 new。如果 old 为空串，则从 s 的开头开始匹配，匹配每一个 UTF-8 序列，最多产生 k + 1 个替换项。

~~~go
func main() {
	fmt.Println(strings.ReplaceAll("abc", "", "moo")) // mooamoobmoocmoo
    fmt.Println(strings.ReplaceAll("Abcaaadd", "aa", "moo")) // Abcmooadd
}
~~~

`func Replace(s , old, new string, n int) string`

和 `ReplaceAll` 类似，只是增加了匹配数量的限制，n 代表最先匹配的数量。如果 n 小于 0，则相当于没有数量限制。

~~~go
func main() {
	fmt.Println(strings.Replace("Abcaaadd", "aa", "moo", 0))   // Abcaaadd
	fmt.Println(strings.Replace("Abcaaaadd", "aa", "moo", 1))  // Abcmooaadd
	fmt.Println(strings.Replace("Abcaaaadd", "aa", "moo", -1)) // Abcmoomoodd
}
~~~

---

`func Split(s, sep string) []string`

`func SplitN(s, sep string, n int) []string`

该方法将 s 按照 sep 的格式进行切分，并返回切分后的 `[]string`；如果 sep 不为空，且 s 中不存在该 sep 时，结果将返回仅含 s 元素的 slice；如果 sep 为空，则将 s 中的每一个 UTF-8 序列做拆分；如果 s 和 sep 都为空自测，结果返回空的 slice。在 `SplitN` 函数中，如果 n 为 -1，则其结果和 `Split` 一致。

~~~go
func main() {
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.Split(" xyz", ""))
	fmt.Printf("%q\n", strings.Split("", "Benardo O'Higgins"))
	fmt.Printf("%q\n", strings.Split("", ""))
}

/**
["a" "b" "c"]
["" "man " "plan " "canal panama"]
[" " "x" "y" "z"]
[""]   --> 包含 1 个元素
[]     --> 不包含任何元素，slice 为空
**/
~~~

那如果是 `splitN` 则对返回的内容做如下限制：

~~~go
n > 0: at most n substrings; the last substring will be the unsplit remainder.
n == 0: the result is nil (zero substrings)
n < 0: all substrings（相当于退化成 Split()）
~~~

作如下示例代码：

~~~go
func main() {
	fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))

	z := strings.SplitN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}

/**
["a" "b,c"]
[] (nil = true)
**/
~~~

`func SplitAfter(s, sep string) []string`

`func SplitAfterN(s, sep string, n int) []string`

`SplitAfter` 和 `Split` 的区别在于前者切分的结果会包含 sep，比如：

~~~go
func main() {
	fmt.Printf("%q\n", strings.SplitAfter("a,b,c", ","))
}

/**
["a," "b," "c"] --> slice 中的第一个元素是 'a,' 而不是 a
**/
~~~

----

`func Title(s string) string`

title 的含义在于，相当于是一篇文章的 Title 命令格式，也就是每一个单次的==首字母大写==。

~~~go
func main() {
	fmt.Println(strings.Title("her royal highness"))  // Her Royal Highness
	fmt.Println(strings.Title("loud noises"))         // Loud Noises
	fmt.Println(strings.Title(""))
}
~~~

但该函数有一个 bug：该函数包含的单词边界规则，不能正确处理 Unicode 标点符号。

`func ToTitle(s string) string`

该函数将 s 的所有 Unicode letters 转化为大写：

~~~go
func main() {
	fmt.Println(strings.ToTitle("her royal highness"))  // HER ROYAL HIGHNESS
	fmt.Println(strings.ToTitle("loud noises"))         // LOUD NOISES
	fmt.Println(strings.ToTitle(""))
}
~~~

`func ToTitleSpecial(c unicode.SpecialCase, s string) string`

该方法将 s 中的==所有== Unicode letters 优先使用给定的 `unicode.SpecialCase` 转化为标题体格式。

~~~go
func main() {
	fmt.Println(strings.ToTitleSpecial(unicode.TurkishCase,
		"dünyanın ilk borsa yapısı Aizonai kabul edilir"))  // DÜNYANIN İLK BORSA YAPISI AİZONAİ KABUL EDİLİR
}
~~~

---

`func ToLower(s string) string`

`func ToLowerSpecial(c unicode.SpecialCase, s string) string`

`func ToUpper(s string) string`

`func ToUpperSpecial(c unicode.SpecialCase, s string) string`

上述函数将指定的 s 转化对应的小写、大写格式。

-----

`func ToValidUTF8(s, replacement string) string`

该函数将 s 的所有无效的 UTF-8 字节序列替换成 replacement。

----

`func Trim(s, cutset string) string`

该函数将 s 中的前导、后缀中的所有包含在 cutset 的 Unicode code point 截取，并返回新的字符串。

~~~go
func main() {
	fmt.Println(strings.Trim("¡¡¡Hello, Gophers!!!", "!¡")) // Hello, Gophers
}
~~~

`func TrimFunc(s string, f func(rune) bool) string`

`func TrimLeft(s, cutset string) string`，==要注意 `cutset` 集合的命名==

`func TrimLeftFunc(s string, f func(rune) bool) string`

`func TrimPrefix(s, prefix string) string`，==要注意此时使用的是 `prefix` 命名==

~~~go
func main() {
	fmt.Println(strings.TrimLeft("¡¡¡Hello, Gophers!!!", "!¡"))   // Hello, Gophers!!!
	fmt.Println(strings.TrimPrefix("¡¡¡Hello, Gophers!!!", "!¡")) // ¡¡¡Hello, Gophers!!!
}
~~~

区别在于，`strings.TrimPrefix()` 只会考虑 s 中和 prefix 完全匹配的前缀。

`func TrimRight(s, cutset string) string`

`func TrimRightFunc(s string, f func(rune) bool) string`

`func TrimSuffix(s, suffix string) string`

`func TrimSpace(s string) string`

该函数将前导的和后置的所有空格字符去掉。

~~~go
func main() {
	z := " \t\n Hello, Gophers \n\t\r\n"
	fmt.Println(strings.TrimSpace(z)) // Hello, Gophers
}
~~~

# 2 strings.Builder

`strings.Builder` 的使用能够**高效率地创建字符串**，最小化内存拷贝。==不允许对已写入内容的 `Builder` 执行拷贝，但对于空内容的 `Builder` 除外==。对 `strings.Builder` 拷贝的检测是==难点==！

~~~go
func main() {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
        // strings.Builder 实现了 io.Writer 接口
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("ignition")
	fmt.Println(b.String())
}
~~~

对于一个 `strings.Builder` 变量，实现了如下方法：

1. `func (b *Builder) Cap() int`
2. `func (b *Builder) Grow(n int)`：增大 Builder 的容量，方法执行后，至少有 n 个字节可被写入；
3. `func (b *Builder) Len() int`：b.Len() == len(b.String())
4. `func (b *Builder) Reset()`：对 Builder 执行重置操作；
5. `func (b *Builder) String() string`
6. `func (b *Builder) Write(p []byte) string`
7. `func (b *Builder) WriteByte(c byte) error`
8. `func (b *Builder) WriteRune(r rune) (int, error)`：写入的是 rune 类型值，也就是 Unicode code point；
9. `func (b *Builder) WriteString(s string) (int, error)`

下面读 `strings.Builder` 的源代码进行分析：

~~~go
// A Builder is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type Builder struct {
	addr *Builder // of receiver, to detect copies by value
	buf  []byte
}
~~~

结构体的定义中，确实是可以将该类型的指针作为其成员。下面按照难易程度作为标准，分 4 个阶段分析 strings.Builder 的源代码：

~~~go
// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int { return len(b.buf) }

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *Builder) Cap() int { return cap(b.buf) }

// Reset resets the Builder to be empty.
func (b *Builder) Reset() {
	b.addr = nil
	b.buf = nil // 为什么此处是赋值 nil？
}
~~~

从上面的 `Len()`、`Cap()` 和 `Reset()` 方法可看出，`strings.Builder` 中实际上是 `[]byte` 这个切片装载了所有写入到 `Bilder` 的内容。

~~~go
// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *Builder) Write(p []byte) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *Builder) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *Builder) WriteRune(r rune) (int, error) {
	b.copyCheck()
	if r < utf8.RuneSelf {
		b.buf = append(b.buf, byte(r))
		return 1, nil
	}
	l := len(b.buf)
	if cap(b.buf)-l < utf8.UTFMax {
		b.grow(utf8.UTFMax)
	}
	n := utf8.EncodeRune(b.buf[l:l+utf8.UTFMax], r)
	b.buf = b.buf[:l+n]
	return n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *Builder) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, s...)
	return len(s), nil
}
~~~

所有的 `WriteXxx` 系列方法在实现时，都调用了内置的 `append()` 将目标值添加到 `b.buf` 切片中。

此处存在的问题在于 `append()` 在执行时，有可能会重新开辟内存并执行拷贝（==底层自动扩容==）。正因为这个可能的内存操作，是会影响性能的。因此引出 `Grow` 方法：

~~~go
// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *Builder) grow(n int) {
	buf := make([]byte, len(b.buf), 2*cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *Builder) Grow(n int) {
	b.copyCheck()
	if n < 0 {
		panic("strings.Builder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}
~~~

调用 `Grow` 后，在需要的情况下（有时候根本不需要），会重新开辟更大的内存空间，用来装载 n 个字节数据。

`strings.Builder` 之所以性能高，原因在于：

~~~go
// String returns the accumulated string.
func (b *Builder) String() string {
    return *(*string)(unsafe.Pointer(&b.buf)) // &(b.buff) 指针创建 Pointer 变量
    
    // 将 Pointer 变量强制类型转化为 *string 类型，并对其指针取值得到 string
}

// unsafe 包下定义类型
type Pointer *ArbitraryType
~~~

使用了一个 `unsafe.Pointer` 的指针转换操作，实现了直接将 `buf []byte` 转换成 string 类型，同时避免了内存申请、分配和销毁的问题。==**我们也可以进行 string 到 `[]byte` 的零内存拷贝和申请转换**==：==**存疑问**==！

~~~go
func StringToBytes(str string) []byte {
	s := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{s[0], s[1], s[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func main() {
	bytes := StringToBytes("中国")
	fmt.Printf("%v \n", bytes)
}
~~~

最后再来看看 `strings.Builder` 不允许复制的问题：

~~~go

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (b *Builder) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*Builder)(noescape(unsafe.Pointer(b)))  // b.addr 设置为 b 变量的指针，为 *strings.Builder 类型
	} else if b.addr != b {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
}
~~~

在所有的 `WriteXxx` 方法和 `Grow` 方法中，都调用了 `copyCheck()`，这个方法是用来**防止复制**的。

1. 在创建 `strings.Builder` 时，`b.addr` 的初始值是 `nil`；
2. 当第一次调用时，会把当前 `b.addr` 赋值为 `strings.Builder` 实例的指针（变量的内存地址）；
3. 后续每次调用，都会检查当前实例的指针是否和 `b.addr` 相等，如果不相等，则引发 panic。

下面来验证：

~~~go
func main() {
	var b1 strings.Builder
	b1.WriteString("ABC") // b1.addr 已被赋值为 &b1

	b2 := b1 // 结构体变量的赋值，b2.addr 值仍然为 &b1

	b2.WriteString("DEF") // copyCheck()执行时 b2.addr 的值为 &b1，和 b2 不相等
}
~~~

在执行 `b2.WriteString("DEF")` 时，引发了 panic！

# 3 strings.Reader

`strings.Reader` 类型通过**从一个 string 读取数据**，该类型实现了 `io.Reader`、`io.Seeker`、`io.ReaderAt`、`io.WriterTo`、`io.ByteScanner`、`io.RuneScanner` 接口。

~~~go
// A Reader implements the io.Reader, io.ReaderAt, io.Seeker, io.WriterTo,
// io.ByteScanner, and io.RuneScanner interfaces by reading
// from a string.
// The zero value for Reader operates like a Reader of an empty string.
type Reader struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}
~~~

**对于 Reader 的零值，可以认为是空字符串的 Reader！**将成员初始化其类型的零值。

`func NewReader(s string) *Reader`

创建一个从 s 读取数据的 Reader，和 `bytes.NewBufferString` 类似，但是更有效率，且为**==只读的==**。

~~~go
func main() {
	name := "中国"
	reader := strings.NewReader(name)
	fmt.Println(reader.Len())
}
~~~

其他方法依次是：

1. `func (r *Reader) Len() int`：获取还未读取的 string 长度；
2. `func (r *Reader) Read(b []byte) (n int, err error) `：从 Reader 中读取字节数组到 b 中；
3. `func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)`：从 Reader 的 off 开始读取字节数组到 b 中；
4. `func (r *Reader) ReadByte() (byte, error)`：从 Reader 中读取一个 byte 值；
5. `func (r *Reader) ReadRune() (ch rune, size int, err error)`：从 Reader 中读取一个 rune 值；
6. `func (r *Reader) Reset(s string)`：使用 s 重置 Reader；
7. `func (r *Reader) Seek(offset int64, whence int) (int64, error)`：按照 offset 和 whence 修改 Reader 当前读取字节索引值；
8. `func (r *Reader) Size() int64`：获取底层 string 的字节长度；
9. `func (r *Reader) UnreadByte() error`：以 byte 为单位，后退已读取的字节位置索引；
10. `func (r *Reader) UnreadRune() error`：以 rune 为单位，后退已读取的字节位置索引；
11. `func (r *Reader) WriteTo(w io.Writer) (n int64, err error)`：将 Reader 的字节数组内容写入到 `io.Writer` 中。

有一处疑问：==为什么指向地址相同？==

~~~go
func test() {
	dst := make([]byte, 2)

	var src = "abc"
	reader := strings.NewReader(src) // *Reader
	n, err := reader.Read(dst)
	fmt.Printf("%d, %v, %s\n", n, err, dst)  // 2, <nil>, ab

	fmt.Printf("%p\n", reader)  // 0xc000004540

	var srcOther = "fff"
	reader.Reset(srcOther)
	fmt.Printf("%p\n", reader) // 0xc000004540
}
~~~

我们看源代码来分析：

~~~go
// Reset resets the Reader to be reading from s.
func (r *Reader) Reset(s string) { *r = Reader{s, 0, -1} }

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *Reader { return &Reader{s, 0, -1} }
~~~

比较难以理解的是 `*r = Reader{s, 0, -1}` 这个实际上是结构体变量的赋值，也就是说相当是将另一个结构体赋给 `*r` 结构体，实质上是会发生成员的赋值。`reader` 的指向关系自然不会改变！

# 4 strings.Replacer

`strings.Replacer` 类型用于字符串的替换，在多 goroutine 中运行时是安全的。

~~~go
// Replacer replaces a list of strings with replacements.
// It is safe for concurrent use by multiple goroutines.
type Replacer struct {
	once   sync.Once // guards buildOnce method
	r      replacer
	oldnew []string
}
~~~

可导出的方法有：

1. `func NewReplacer(oldnew ...string) *Replacer`：使用提供的多组 old、new 字符串对创建并返回一个 `*Replacer`。其替换是依次进行的，匹配时不会重叠；
2. `func (r *Replacer) Replace(s string) string`：执行替换，并返回替换后的字符串；
3. `func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)`：执行替换后，将结果写入到 w 中。
