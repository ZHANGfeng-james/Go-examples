# Buffer

用一句话说清楚 bytes.Buffer 模型的功能：用于调度数据的简单字节缓冲器。



~~~go
// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int    // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.
}
~~~

对 Buffer 的理解，重在弄清楚：

* off 域可称之为“已读标记位置”，`read at &buf[off], write at &buf[len(buf)]`
* buf 域，可读取的范围是 buf[off:len(buf)]，`contents are the bytes buf[off : len(buf)]`

可以看出 bytes.Buffer 底层是用 []byte 封装的，另外模型还封装了一个 readOp 类型的域 lastRead 表示的是最后读取内容的操作类型。readOp 的操作类型包括：

~~~go
// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)
~~~

readOp 类型对 int8 类型进行类型重定义，同时定义了一系列的 const 值，用于表示特定的操作类型。

~~~go
// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")
var errNegativeRead = errors.New("bytes.Buffer: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

// MinRead is the minimum slice size passed to a Read call by
// Buffer.ReadFrom. As long as the Buffer has at least MinRead bytes beyond
// what is required to hold the contents of r, ReadFrom will not grow the
// underlying buffer.
const MinRead = 512

var errUnreadByte = errors.New("bytes.Buffer: UnreadByte: previous operation was not a successful read")
~~~

对于上述的 maxInt 可做这样的测试程序：

~~~go
// 64 位 整数值
const maxInt = ^uint(0)

func main() {
	// 1111111111
	// 1111111111
	// 1111111111
	// 1111111111
	// 1111111111
	// 1111111111
	// 1111       uint
	fmt.Printf("%b， %T.\n", maxInt, maxInt)

	value := ^uint8(0)
	// 11111111， uint8.
	fmt.Printf("%b， %T.\n", value, value)
    
    value = ^uint8(1)
	// 11111110， uint8.
	fmt.Printf("%b， %T.\n", value, value)
}
~~~

从上述结果可以看出，^uint8(x) 表示的含义是使用当前 uint8 类型的最大值（二进制位全为 1 的值）进行异或运算。

~~~go
// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }
~~~

注意此处返回的是 b.buf[b.off:] 也就是还未读部分的字节切片内容，此外，还需要考虑在 Read、Write、Reset 和 Truncate 方法调用后对 Bytes() 的影响。什么时候会改变 len(b.buf)？

~~~go
func main() {
	buf := bytes.NewBufferString("a")
	fmt.Printf("%d, %d.\n", buf.Cap(), buf.Len())

	other := buf.Bytes()
	fmt.Printf("%d, %d, %s.\n", cap(other), len(other), other)

	other[0] = 0x62
	fmt.Printf("%s.\n", buf.String())
}

PS C:\Users\Developer\sample> go run main.go
8, 1.
8, 1, a.
b.
~~~

通过 Bytes() 获取到底层字节缓冲后，可直接对底层数据进行修改。这是一个很重要的特征！

~~~go
// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (b *Buffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}
~~~

上述 String() 获取的是还未读的字符串内容，相应的是从 b.buf[b.off:] 的字节切片中进行转化。

~~~go
// NewBuffer creates and initializes a new Buffer using buf as its
// initial contents. The new Buffer takes ownership of buf, and the
// caller should not use buf after this call. NewBuffer is intended to
// prepare a Buffer to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }

// NewBufferString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBufferString(s string) *Buffer {
	return &Buffer{buf: []byte(s)}
}
~~~

可以分别使用 string 和 []byte 获取到 *Buffer 类型值。

~~~go
// empty reports whether the unread portion of the buffer is empty.
func (b *Buffer) empty() bool { return len(b.buf) <= b.off }

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *Buffer) Len() int { return len(b.buf) - b.off }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (b *Buffer) Cap() int { return cap(b.buf) }
~~~

empty() 和 Len() 都针对的是 b.buf 中的 the unread portion of the buffer。对于上述方法，一个特别有意思的调用是：

~~~go
func main() {
	slice := []byte("a")
	fmt.Printf("%d, %d.\n", cap(slice), len(slice))

	other := slice[:0]
	fmt.Printf("%d, %d.\n", cap(other), len(other))

	buf := bytes.NewBufferString("a")
	fmt.Printf("%d, %d.\n", buf.Cap(), buf.Len())

	unRead := buf.Bytes()
	// 61.
	fmt.Printf("%0 x.\n", unRead)
	// 8, 1.
	fmt.Printf("%d, %d.\n", cap(unRead), len(unRead))
}

PS C:\Users\Developer\sample> go run main.go
1, 1.
1, 0.
8, 1.
61.
8, 1.
~~~

特别的地方是：buf.Cap() 和 buf.Len() 值分别是 8 和 1，但是问题在于为什么底层自动分配了 8 个字节的数组作为整个容器的初始容量？如果是运行下面的内容，却得到了不同的结果：

~~~go
func main() {
	buf := bytes.NewBufferString("a")
	fmt.Printf("%d, %d.\n", buf.Cap(), buf.Len())
}

PS C:\Users\Developer\sample> go run main.go
32, 1.
~~~

![](./pics/Snipaste_2021-04-14_11-44-08.png)

从源代码来看，确实是没有关于底层 slice 的分配相关的内容，这有可能是关于结果会变化的原因。

> 这部分内容可以参考：runtime 包中叫做 stringtoslicebyte 的函数！

~~~go
// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
	b.lastRead = opInvalid
}
~~~

b.buf[:0] 对 b.buf 做重新切片，切片后 len 值为 0，cap 值不变。也就是说，底层的字节数组是不改变的！

~~~go
// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (b *Buffer) Truncate(n int) {
	if n == 0 {
		b.Reset()
		return
	}
	b.lastRead = opInvalid
	if n < 0 || n > b.Len() {
		panic("bytes.Buffer: truncation out of range")
	}
	b.buf = b.buf[:b.off+n]
}
~~~

相当于是从 [0, b.off + n] 进行截取，抛弃掉其他的 slice 内容。

为了能够理解 Grow() 方法，以及 Write 方法，先来理解 Read 方法：

~~~go
// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (b *Buffer) Read(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
    // 实际读取 b.buf[b.off:] 切片数据，返回实际拷贝的字节数
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead // Any other read operation.
	}
	return n, nil
}
~~~

如果此时 b.empty() 为 true，即表示 buffer 中已经没有数据可读了！p 切片量，实际相当于是 output。

~~~go
// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (b *Buffer) ReadByte() (byte, error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		return 0, io.EOF
	}
	c := b.buf[b.off]
	b.off++
	b.lastRead = opRead
	return c, nil
}
~~~

读取一个字节的值，并返回。同样的，需要判断是否已经无可读数据。

~~~go
// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := b.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return line, err
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (b *Buffer) readSlice(delim byte) (line []byte, err error) {
    // 在 b.buf[b.off:] 切片中找到 delimiter 界限符，返回其索引
	i := IndexByte(b.buf[b.off:], delim)
	end := b.off + i + 1
	if i < 0 {
		end = len(b.buf)
		err = io.EOF
	}
	line = b.buf[b.off:end]
	b.off = end
	b.lastRead = opRead
	return line, err
}
~~~

读取 buffer 中首次遇到 delimiter 之前的内容，包括 delimiter 字节内容。

~~~go
// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (b *Buffer) ReadString(delim byte) (line string, err error) {
	slice, err := b.readSlice(delim)
	return string(slice), err
}
~~~

将读取到的 line 转化为 string 并返回。

~~~go
// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (b *Buffer) ReadRune() (r rune, size int, err error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		return 0, 0, io.EOF
	}
	c := b.buf[b.off]
	if c < utf8.RuneSelf {
		b.off++
		b.lastRead = opReadRune1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.buf[b.off:])
	b.off += n
	b.lastRead = readOp(n) // 设置为实际读取的字节数（rune 最大是4个字节）
	return r, n, nil
}
~~~

和读取 byte 值类似，不同之处在于此处读取到的是 rune 值。

~~~go
// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer) Next(n int) []byte {
	b.lastRead = opInvalid
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}
~~~

读取 b.buf[b.off : b.off + n] 的字节内容，并返回。

~~~go
// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a successful ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (b *Buffer) UnreadRune() error {
	if b.lastRead <= opInvalid {
		return errors.New("bytes.Buffer: UnreadRune: previous operation was not a successful ReadRune")
	}
	if b.off >= int(b.lastRead) {
		b.off -= int(b.lastRead)
	}
	b.lastRead = opInvalid
	return nil
}

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
func (b *Buffer) UnreadByte() error {
	if b.lastRead == opInvalid {
		return errUnreadByte
	}
	b.lastRead = opInvalid
	if b.off > 0 {
		b.off--
	}
	return nil
}
~~~

上面的 UnreadRune 和 UnreadByte 方法都是可 b.lastRead 的类型是相关的。如果 b.lastRead 是 opInvalid 或者是 opRead，在执行 UnreadRune 时会报错。

~~~go
// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
    m := b.grow(n)
	b.buf = b.buf[:m] // reslice
}
~~~

先从宏观上理解 Grow(n int) 带来的实际效果：在需要的时候，扩大底层 []byte 的容量。此时，m 值下一次写入的位置索引值。

~~~go
// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
    // cap(b.buf) 能够容纳 n 个字节，这种情况下只需要 reslice
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 { // 相当于没有任何数据可读
        b.Reset() // 重置，len(b.buf) 为0，但底层字节数组仍可用
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
        // 64 个字节
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m { // n + m <= c/2 
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
    return make([]byte, n) // len(slice) 结果为 n
}
~~~

扩容规则比较难以理解，其原则在于：能够用当前的 cap(b.buf) 容纳的，绝不重新创建底层字节数组。其效果也仅仅只是用来扩大底层 []byte 的容量，而不会扩大“窗口”的长度。

~~~go
// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lastRead = opInvalid
	if nBytes := b.Len(); nBytes > 0 {
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	b.Reset()
	return n, nil
}
~~~

WriteTo 表示将 bytes.Buffer 的**可读内容**写入到 io.Writer 中，直到 buffer 无输出的内容（已枯竭）。

~~~go
// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.lastRead = opInvalid
	for {
		i := b.grow(MinRead) // 循环从 io.Reader 中读取，每次增长 512 字节
		b.buf = b.buf[:i]
		m, e := r.Read(b.buf[i:cap(b.buf)]) // reslice
		if m < 0 {
			panic(errNegativeRead)
		}

		b.buf = b.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}
~~~

ReadFrom 从 io.Reader 中读取字节内容到 bytes.Buffer 中，直到 io.Reader 枯竭。此时，ReadFrom 方法内部，并不会修改 b.off 值。

~~~go
// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (b *Buffer) WriteByte(c byte) error {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(1)
	if !ok {
		m = b.grow(1)
	}
	b.buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (b *Buffer) WriteRune(r rune) (n int, err error) {
	if r < utf8.RuneSelf {
		b.WriteByte(byte(r))
		return 1, nil
	}
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = b.grow(utf8.UTFMax)
	}
	n = utf8.EncodeRune(b.buf[m:m+utf8.UTFMax], r)
	b.buf = b.buf[:m+n]
	return n, nil
}
~~~

WriteByte 和 WriteRune 分别向 bytes.Buffer 中写入 byte 和 rune 内容。

~~~go
// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(s))
	if !ok {
		m = b.grow(len(s))
	}
	return copy(b.buf[m:], s), nil
}
~~~

