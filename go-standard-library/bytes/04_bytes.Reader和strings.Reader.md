> 每个人都需要给自己的人生设立目标。没有目标的生活，就像没有目的地的航行，找不到前进的方向。只有清晰地知道想要什么，才能主动突破周围环境的限制，找到行动的动力，塑造更好的自己。

本篇我们把 bytes.Reader 和 strings.Reader 结合在一起比较，看看有什么不一样，使用场景分别是什么？

# 1 NewReader

bytes.Reader 的定义：

~~~go
type Reader struct {
	s        []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}
~~~

strings.Reader 的定义：

~~~go
type Reader struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}
~~~

特别需要指出的是：

1. s：相当于是缓存的内容，区别在于 bytes 包中**缓存的是 []byte**，而 strings 包中**缓存的是 string**；
2. i：当前正在读取的位置索引；
3. prevRune：前一个 rune 值的位置索引。

**如下的所有实现的方法，都是通过上述 s、i、prevRune 实现读取的**。对于 bytes.Reader 而言，**封装的用意**在于：读取 []byte 的 byte 内容；对于 strings.Reader 而言，**封装的用意**在于：读取 string 中的 byte 内容。

bytes 和 strings 包都封装了 Reader 结构体类型，从类型注释上来看，它们**都实现了接口**：

* io.Reader
* io.ReaderAt
* io.Seeker
* io.WriterTo
* io.ByteScanner
* io.RuneScanner

也就是具备上述**接口定义的功能**。

NewReader 用于创建 Reader 实例：

~~~go
// NewReader returns a new Reader reading from b.
func NewReader(b []byte) *Reader {
    return &Reader{b, 0, -1} 
}
~~~

对应的 strings 包，有相同的函数：

~~~go
// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *Reader { 
    return &Reader{s, 0, -1} 
}
~~~

和包名的区分一样，strings 包对应的是参数是 string，而对于 byts 包对应的参数是 []byte。

# 2 实现 io.Reader 接口

对于 io.Reader 接口：

~~~go
type Reader interface {
    Read(p []byte) (n int, err error)
}
~~~

我这样理解 io.Reader 接口定义的功能：**可被读取的能力**！在实现 `Read(p []byte) (n int, err error)` 时，实际上是把该接口的实现类型（比如 bytes.Reader 或者是 strings.Reader）中缓存的内容读取到 p 中。

比如：**bytes 实现 io.Reader 接口**：

~~~go
// Read implements the io.Reader interface.
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
    // 尽可能地读取 r.s[r.i:] 剩下的内容，要么填满 b，要么 r.s 已读完
    n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}
~~~

对应到 strings 包下：

~~~go
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
    // 尽可能地读取 r.s[r.i:] 剩下的内容，要么填满 b，要么 r.s 已读完
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}
~~~

**代码是完全一样的！**

特别的，我们新接触了一个函数：`copy(dst, src []Type) int`

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

copy **内建函数**是用来实现 `[]Type` 的拷贝，其中 Type 代表（stand-in）的是 Go 中的任意一种数据类型。在 copy 的含义，表示 dst 和 src 是相同的类型（Type）的切片类型。既然 byes.Reader.Read 和 strings.Reader.Read 方法都是基于 copy 的，那我们来测试一下 copy 的功能：

~~~go
func TestBuiltinCopy(t *testing.T) {
	dst := make([]byte, 0)
	n := copy(dst, "Michoi") // 确实能够将 string 实例的字节内容通过 copy 拷贝到 []byte 中！
	fmt.Println("Success copy byte size:", n)

	fmt.Printf("%q\n", dst)
}

""
// dst := make([]byte, 0) --> dst := make([]byte, 10)
Michoi
"Michoi\x00\x00\x00\x00"
~~~

也就是说，如果 len(dst) 的长度为 0，则根本不会执行任何拷贝，因为 dst 没有开辟任何内存空间，用于存放拷贝的字节内容。其拷贝的结果——n 值——是 len(dst) 和 len(src) 的较小的值。**copy 内建函数不会对 dst 做内存的重新分配，仅仅依靠 dst 本身的 len 和 cap 执行拷贝动作**。另外，copy 函数是不会返回 error 值的（但 Reader 方法是会返回 error 值）。在上述测试示例中，即便 dst 的值是 nil，也不会有任何 panic 发生！

现在我们测试 bytes.Reader.Read 和 strings.Reader.Read 方法：

~~~go
func TestReader(t *testing.T) {
	reader := strings.NewReader("Michoi")
	p := make([]byte, 0)
	n, err := reader.Read(p)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}
Success Read size:0; p:""
// p := make([]byte, 0) --> p := make([]byte, 10)
Success Read size:6; p:"Michoi\x00\x00\x00\x00"
~~~

从 Read 方法的实现来看，有 2 种特殊的返回：

1. n 值为 0，error 值为 nil：可能的原因是 len(p) 值为 0；
2. n 值为 0，error 值为 io.EOF：s 已经读取完毕，没有其他可读取的内容。

# 3 实现 io.ReaderAt 接口

对于 io.ReaderAt 接口：

~~~go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
~~~

和 io.Reader 接口定义的功能不同的是：从 s 的 off 位置索引**开始读取**字节内容，并存储到 p 实例中。

比如：**bytes 实现 io.Reader 接口**：

~~~go
// ReadAt implements the io.ReaderAt interface.
func (r *Reader) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("bytes.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}
~~~

对应到 strings 包下：

~~~go
func (r *Reader) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("strings.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF // 没有把 b 读满，意味着 r.s 已经到达 io.EOF 了
	}
	return
}
~~~

**代码是完全一样的！**

测试代码如下：

~~~go
func TestReaderAtInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")
	p := make([]byte, 10)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		if err == io.EOF {
			fmt.Println("r.s --> io.EOF")
		} else {
			t.Fatal(err.Error())
		}
	}
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)

	reader.Read(p)
	fmt.Println(string(p))
}
r.s --> io.EOF
Success Read size:4; p:"choi\x00\x00\x00\x00\x00\x00"
Michoi
~~~

ReaderAt 方法的另一个特点在于：对 r.s 的读取后，并不会修改 r.i 的值。也就是说，紧接着的一次读取仍然从 r.i 开始。

# 4 实现 io.Seeker 接口

对于 io.Seeker 接口：

~~~go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
~~~

该接口表示的是这样的功能：能够修改当前缓存区的下次读取（或者写入）时的起始位置索引值—— r.i 值。其中（whence: from where）：

* SeekStart：以最开始的位置作为参照，设置 r.i 的值；
* SeekCurrent：以当前 r.i 的值作为参照，设置 r.i 的值；
* SeekEnd：以可读内容末尾位置作为参照，设置 r.i 的值。

Seek 方法设置后的返回值表示的是：当前可被读取（或写入）的 r.i 值（相对于可被读取的起始位置来说的相对值）。在上述 Seek 方法中，offset 作为一个偏移量，相当于是根据 whence 策略的不同，**由不同的值和 offset 的和计算得到了 abs**。

比如：**bytes 实现 io.Reader 接口**：

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
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}
~~~

对应到 strings 包下：

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

**代码是完全一样的！**

测试代码如下：

~~~go
func TestSeekerInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")
	p := make([]byte, 10)

	abs, err := reader.Seek(-4, io.SeekEnd)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("abs:", abs)

	n, _ := reader.Read(p)
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}
abs: 2
Success Read size:4; p:"choi\x00\x00\x00\x00\x00\x00"
~~~

# 5 实现 io.WriterTo 接口

对于 io.WriterTo 接口：

~~~go
type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}
~~~

io.Writer 接口实际上表示的是：**能够将自身的内容写入到入参 io.Writer 实例的能力**。

比如：**bytes 实现 io.Reader 接口**：

~~~go
// WriteTo implements the io.WriterTo interface.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	b := r.s[r.i:]
	m, err := w.Write(b) // 重新构建了 []byte
	if m > len(b) {
		panic("bytes.Reader.WriteTo: invalid Write count")
	}
	r.i += int64(m)
	n = int64(m)
	if m != len(b) && err == nil {
		err = io.ErrShortWrite
	}
	return
}
~~~

对应到 strings 包下：

~~~go
// WriteTo implements the io.WriterTo interface.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	s := r.s[r.i:]
	m, err := io.WriteString(w, s) // 重新构建了 []byte
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

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
// If w implements StringWriter, its WriteString method is invoked directly.
// Otherwise, w.Write is called exactly once.
func WriteString(w Writer, s string) (n int, err error) {
	if sw, ok := w.(StringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}
~~~

上述两者的区别是：bytes.Reader.WriterTo 使用的方法是 w.Write(b)，而 strings.Reader.WriterTo 使用的方式是 io.WriteString(w, s)。

测试程序：

~~~go
func TestWriterToInterface(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	reader := strings.NewReader("Michoi")

	n, err := reader.WriteTo(buffer)
	if err != nil {
		if err == io.ErrShortWrite {
			fmt.Println("r.s --> io.ErrShortWrite")
		} else {
			t.Fatal(err.Error())
		}
	}
	fmt.Println("n:", n)
	fmt.Println(buffer.String())
}
n: 6
Michoi
~~~

# 6 实现 io.ByteScanner 接口

对于 io.ByteScanner 接口（**含有一个内嵌接口**）：

~~~go
type ByteScanner interface {
    ByteReader
    UnreadByte() error
}
~~~

接口封装用于表示：具备有获取上一次 ReadByte() 方法的结果的能力，以及包含有获取下一次读取字节的能力。

比如：**bytes 实现 io.Reader 接口**：

~~~go
// ReadByte implements the io.ByteReader interface.
func (r *Reader) ReadByte() (byte, error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

// UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
func (r *Reader) UnreadByte() error {
	if r.i <= 0 {
		return errors.New("bytes.Reader.UnreadByte: at beginning of slice")
	}
	r.prevRune = -1
	r.i--
	return nil
}
~~~

strings 包下的相关方法，是和 bytes 包中的**代码是完全一样的！**

测试程序：

~~~go
func TestByteScannerInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")

	err := reader.UnreadByte()
	if err != nil {
		fmt.Println(err.Error()) // 抛出了一个 error
	}

	p := make([]byte, 10)
	n, err := reader.Read(p)
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}
strings.Reader.UnreadByte: at beginning of string
Success Read size:6; p:"Michoi\x00\x00\x00\x00"
~~~

和 []byte 相关的测试程序：

~~~go
func TestSlice(t *testing.T) {
	buffer := []byte("Michoi")
	fmt.Printf("len(buffer)=%d, %q\n", len(buffer), buffer)

	start := -1
	slice := buffer[start:1]
	fmt.Printf("%q\n", slice)
}
panic: runtime error: slice bounds out of range [-1:]
~~~

在第 6 行，抛出了一个 panic，也就是说在引用切片时，不能超出范围。

# 7 实现 io.RuneScanner 接口

对于 io.RuneScanner 接口（**含有一个内嵌接口**）：

~~~go
type RuneScanner interface {
    RuneReader
    UnreadRune() error
}
~~~

io.RuneScanner 接口封装的功能和 io.ByteScanner 类似，不同之处是：后者是遍历的是 byte（1个字节），前者遍历的元素是 rune 类型（4个字节）的。io.RuneScanner 接口的功能：**具有以 Rune 的格式扫描的能力**。

比如：**bytes 实现 io.Reader 接口**：

~~~go
// ReadRune implements the io.RuneReader interface.
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
	ch, size = utf8.DecodeRune(r.s[r.i:]) // ch, size = utf8.DecodeRuneInString(r.s[r.i:])
	r.i += int64(size)
	return
}

// UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
func (r *Reader) UnreadRune() error {
	if r.i <= 0 {
		return errors.New("bytes.Reader.UnreadRune: at beginning of slice")
	}
	if r.prevRune < 0 {
		return errors.New("bytes.Reader.UnreadRune: previous operation was not ReadRune")
	}
	r.i = int64(r.prevRune)
	r.prevRune = -1
	return nil
}
~~~

测试程序：

~~~go
func TestRuneScanner(t *testing.T) {
	reader := strings.NewReader("中\xc5国")

	for {
		ch, size, err := reader.ReadRune()
		if err != nil {
			t.Fatal(err.Error())
		}

		if !utf8.ValidRune(ch) {
			fmt.Println("InvalidRune")
		}

		fmt.Printf("%q, %d\n", ch, size)
	}
}
'中', 3
'�', 1
'国', 3
    reader_test.go:100: EOF
~~~

# 8 Size 和 Reset 方法

Size 和 Reset 方法很简单：

~~~go
// Size returns the original length of the underlying byte slice.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (r *Reader) Size() int64 { 
    return int64(len(r.s)) 
}

// Reset resets the Reader to be reading from b.
func (r *Reader) Reset(b []byte) { 
    *r = Reader{b, 0, -1} 
}
~~~

# 9 bytes.Reader 和 strings.Reader 的比较

bytes.Reader 和 strings.Reader 的区别要从底层结构体的区别说起：

string：字符串类型

~~~go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string

type stringStruct struct {
	str unsafe.Pointer
	len int
}
~~~

[]byte：byte 的切片类型

~~~go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
~~~

从结构体的定义来看，string 和 []byte 实际上相差并不是很大。

我们把这个问题的剖析留到下一篇！

