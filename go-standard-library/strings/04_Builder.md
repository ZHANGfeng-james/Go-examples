Go 中的标准库 strings，用于处理 UTF-8 编码的字符串，可将 strings 包当作是一个（UTF-8编码格式的）字符串处理工具箱。

用一句话来展示 Go 项目封装 `strings.Builder` 的用意：能够**高效率地创建字符串**，**最小化内存拷贝**。

`*strings.Builder` 类型的方法：

* `func (b *Builder) Cap() int`：获取底层 []byte 的 cap；
* `func (b *Builder) Grow(n int)`：扩大底层 []byte 的容量；
* `func (b *Builder) Len() int`：获取底层 []byte 的字节长度；
* `func (b *Builder) Reset()`：让底层 []byte 清空；
* `func (b *Builder) String() string`：将底层的 []byte 转为 string；
* `func (b *Builder) Write(p []byte) (int, error)`：向 `*Builder` 中写入 []byte；
* `func (b *Builder) WriteByte(c byte) error`：向 `*Builder` 中写入 byte；
* `func (b *Builder) WriteRune(r rune) (int, error)`：向 `*Builder` 中写入 rune；
* `func (b *Builder) WriteString(s string) (int, error)`：向 `*Builder` 中写入 string。

示例程序如下：

~~~go
func TestBuilder(t *testing.T) {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("ignition")
	b.WriteRune('中')
	b.WriteByte('\xcc')
	b.Write([]byte("Michoi"))
	fmt.Println(b.String(), "cap:", b.Cap(), "; len:", b.Len())

	b.Reset()
	fmt.Println(b.String())
}
3...2...1...ignition中�Michoi cap: 32 ; len: 30

~~~

下面，我们用一种特别的方式来分析 strings.Builder 的源代码：

~~~go
// A Builder is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type Builder struct {
	addr *Builder // of receiver, to detect copies by value
	buf  []byte
}
~~~

结构体的定义中，佐证了将自身类型的指针作为其 Field 的形式。**对于一个命名的结构体类型来说，不能再包含该类型的成员，即一个聚合的值不能包含它自身；但是该结构体却可以包含该命名类型的指针类型的成员**。（用途）在创建递归数据结构时，比如创建链表和树结构时就使用了这种方式。

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
		b.addr = (*Builder)(noescape(unsafe.Pointer(b)))  
        // b.addr 设置为 b 变量的指针，为 *strings.Builder 类型
	} else if b.addr != b {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
}
~~~

源代码中调用 copyCheck() 的地方包括：Grow(n int)、Write(p []byte)、WriteByte(c byte)、WriteRune(r rune)、WriteString(s string)。我们简单看一个上述代码的实现：

~~~go
// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *Builder) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}
~~~

对于一个新建的变量 `var sb strings.Builder` 来说，结构体成员 addr 值为 nil；此时，unsafe.Pointer(b) 将 *strings.Builder 转化为 unsafe.Pointer 类型；调用 noescap 函数时，uintptr(p) 将 unsafe.Pointer 转化为 uintptr；但是在返回时，再次将其转化为了 unsafe.Pointer，仅仅做了一次异或操作。

> unsafe.Pointer 类型是这样一类值，指向的是任何可寻址的值的指针，是指针值和 uintptr 之间的桥梁；uintptr 是 Go 中的内建数据类型，是一种数值类型，其代表任何指针的位模式，也就是原始的内存地址值。

~~~go
func main() {
	var value = 4
	fmt.Printf("%b.\n", value)

	fmt.Printf("%b.\n", value^0)

	// ^
	// 1 ^ 0 = 1
	// 0 ^ 0 = 0
}
~~~

noescape 函数返回后，又一次做了强制类型转换，将返回值转化成了 *Builder 类型，同时赋值给了 addr 字段。经过这个转换，实际上就是将当前声明变量的内地地址值赋值给了自身结构体的 add 域。

copyCheck 方法的目的是来判断此时操作的变量是否已经发生了改变，如若发现已改变，则直接抛出 panic。如果不在 runtime 状态下 recover，则会出现应用进程崩溃的现象，导致进程终止。

`addr` 字段存在的意义在于，保存其所属值所在的内存地址。如此一来，一旦这个值被拷贝了，使用内存地址比较的方式会可以检测出来！

~~~go
// String returns the accumulated string.
func (b *Builder) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}
~~~

上述方法的实现中，&b.buf 实际上是分为 2 个部分：

1. b.buf：是一种 selector **表达式**（并不是一种 operator），是指结构体中的 buf 域（此处做了一次 *b 取指针指向的变量值）；
2. &(b.buf)：取 buf 域的地址，也就是 *[]byte 类型。

unsafe.Pointer(&b.buf) 将 *[]byte 类型转化为 unsafe.Pointer，再次类型转化为 *string，并返回取指针值后的结果作为其结果——字符串表示形式。

~~~go
func main() {
	var str = "爱"
	// string --> []byte
	bytes := []byte(str)
	fmt.Printf("% x.\n", bytes)

	// []byte --> string
	fmt.Println(getString(bytes))
    
    // []byte --> string
    fmt.Println(string(bytes))
}

func getString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

PS C:\Users\Developer\sample> go run main.go
e7 88 b1.
爱
爱
~~~

上述实际上就是 []byte 与 string 相互转换的方法。

注意，上述方法存在 2 种方式：[]byte --> string，但是前一种方式省去了类型转换的开销，效率会高很多。

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

因为 b.buf 的类型是 []byte，是一种指针类型，其默认的初始值是 nil。因此对于 Reset() 而言，需要赋值为 nil。

~~~go
// The copy built-in function copies elements from a source slice into a
// destination slice. (As a special case, it also will copy bytes from a
// string to a slice of bytes.) The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
func copy(dst, src []Type) int

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

Grow(n int) 方法并不是总是需要调用 grow(n)，如果此时已存在能容量 n 字节的空余空间，则不会调用 grow(n)。另外，grow(n int) 中关于 b.buf 的扩容值是有规律的，新切片的 cap 值是 2 * cap(b.buf) + n。

~~~go
// The append built-in function appends elements to the end of a slice. If
// it has sufficient capacity, the destination is resliced to accommodate the
// new elements. If it does not, a new underlying array will be allocated.
// Append returns the updated slice. It is therefore necessary to store the
// result of append, often in the variable holding the slice itself:
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
// As a special case, it is legal to append a string to a byte slice, like this:
//	slice = append([]byte("hello "), "world"...)
func append(slice []Type, elems ...Type) []Type

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
    // characters below RuneSelf are represented as themselves in a single byte.
	if r < utf8.RuneSelf {
		b.buf = append(b.buf, byte(r))
		return 1, nil
	}
	l := len(b.buf)
    // maximum number of bytes of a UTF-8 encoded Unicode character.
	if cap(b.buf)-l < utf8.UTFMax {
		b.grow(utf8.UTFMax)
	}
    // 方法返回实际写入的字节数
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

最后再来看看 `strings.Builder` 不允许复制的问题：在所有的 `WriteXxx` 方法和 `Grow` 方法中，都调用了 `copyCheck()`，这个方法是用来**防止复制**的。

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