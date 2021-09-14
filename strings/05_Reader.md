Go 中的标准库 strings，用于处理 UTF-8 编码的字符串，可将 strings 包当作是一个（UTF-8编码格式的）字符串处理工具箱。

`strings.Reader` 类型通过**从一个 string 读取数据**，该类型实现了如下接口：

* io.Reader
* io.ReaderAt
* io.ByteReader
* io.ByteScanner
* io.RuneReader
* io.RuneScanner
* io.Seeker
* io.WriterTo

其结构体构成如下：

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

~~~go
// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *Reader {
	 return &Reader{s, 0, -1} 
}
~~~

创建一个从 s 读取数据的 Reader，和 `bytes.NewBufferString` 类似，但是更有效率，且为**只读的**。

~~~go
func main() {
	name := "中国"
	reader := strings.NewReader(name)
	fmt.Println(reader.Len())
}

PS C:\Users\Developer\sample> go run main.go
6
~~~

下面我们用最直接的方式分析源代码：

~~~go
// Len returns the number of bytes of the unread portion of the
// string.
func (r *Reader) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

// Size returns the original length of the underlying string.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (r *Reader) Size() int64 { return int64(len(r.s)) }
~~~

Len() 获取到的是 Reader 还未读取的字节数。

~~~go
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
    // 读取 r.s[r.i:] 剩余的内容
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

func (r *Reader) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("strings.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
    // offset 表示偏移量
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

func (r *Reader) ReadByte() (byte, error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if r.i >= int64(len(r.s)) {
		r.prevRune = -1
		return 0, 0, io.EOF
	}
	r.prevRune = int(r.i)
	if c := r.s[r.i]; c < utf8.RuneSelf {
		r.i++
		return rune(c), 1, nil
	}
	ch, size = utf8.DecodeRuneInString(r.s[r.i:])
	r.i += int64(size)
	return
}

// The copy built-in function copies elements from a source slice into a
// destination slice. (As a special case, it also will copy bytes from a
// string to a slice of bytes.) The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
func copy(dst, src []Type) int
~~~

上述方法大多使用到了内置的 copy 函数，用于 src 的切片拷贝到 dst 中。此时，实际发生的拷贝内容由 dst 的 len 决定：

~~~go
func main() {
	name := "abcdefg"
	reader := strings.NewReader(name)
	fmt.Println(reader.Len())

    // slice 的 len 值是 4，cap 值是 7
	slice := make([]byte, reader.Len()/2+1, reader.Len())
	n, err := reader.Read(slice)
	fmt.Println(n, err)
}

PS C:\Users\Developer\sample> go run main.go
7
4 <nil>
~~~

上述代码说明了，实际发生的拷贝字节个数是 reader.Len()/2 + 1，也就是 dst 的 len 值。

~~~go
func (r *Reader) UnreadByte() error {
	if r.i <= 0 {
		return errors.New("strings.Reader.UnreadByte: at beginning of string")
	}
	r.prevRune = -1
	r.i--
	return nil
}

func (r *Reader) UnreadRune() error {
	if r.i <= 0 {
		return errors.New("strings.Reader.UnreadRune: at beginning of string")
	}
	if r.prevRune < 0 {
		return errors.New("strings.Reader.UnreadRune: previous operation was not ReadRune")
	}
	r.i = int64(r.prevRune)
	r.prevRune = -1
	return nil
}
~~~

UnreadXxx 方法，实际的作用是将 Reader 中的标记回滚。UnreadByte 方法，是回滚一个字节标记；UnreadRune 方法，是回滚上一个读取到的 rune 值的标记。

~~~go
// Seek implements the io.Seeker interface.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	r.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("strings.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("strings.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}
~~~

Seek 方法用于读取从 whence 开始的指定 offset 位置的内容，其中 whence 的只有 3 种，分别是：io.SeekStart、io.SeekCurrent 和 io.SeekEnd。

~~~go
// WriteTo implements the io.WriterTo interface.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	s := r.s[r.i:]
	m, err := io.WriteString(w, s)
	if m > len(s) {
		panic("strings.Reader.WriteTo: invalid WriteString count")
	}
	r.i += int64(m)
	n = int64(m)
	if m != len(s) && err == nil {
		err = io.ErrShortWrite
	}
	return
}
~~~

Reader 将自身的 string 内容，读取到 io.Writer 接口的值中。

~~~go
// Reset resets the Reader to be reading from s.
func (r *Reader) Reset(s string) { *r = Reader{s, 0, -1} }
~~~

类似于重置，将 Reader 中的 s 重置为指定内容。

示例程序如下：

~~~go
func main() {
	name := "abcdefg"
	reader := strings.NewReader(name)
	fmt.Println(reader.Len())

	slice := make([]byte, reader.Len()/2+1, reader.Len())
	n, err := reader.Read(slice)
	fmt.Println(n, err)

	reader.Reset(name)

	reader.Seek(int64(reader.Len()), io.SeekEnd)
	n, err = reader.Read(slice)
	fmt.Println(n, err)
}

PS C:\Users\Developer\sample> go run main.go
7
4 <nil>
0 EOF
~~~

另外，有一个疑问在于：为什么在 Reader 源代码中多次使用到了 `len(r.s)`，而不把长度信息存储在某个地方？这样做其实根本不会浪费性能。因为字符串值中本身就存着字符串的长度，参考 stringStruct 类型定义。况且，把这种信息记录在 reader 内会造成额外的维护成本。



有一处疑问：为什么指向地址相同？

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
