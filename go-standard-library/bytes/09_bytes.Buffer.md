> 一年，每天1个小时，看完 2 本书：TCP/IP 详解，Windows 高级编程。基础知识，是你最重要的部分。如果基础知识越补越快的话，你的学习能力就会越来越快。因此，**只要你愿意，有成长的渴望，你就一定能够挤出时间**。主动驱动自己，掌控自己的时间。弄清楚为什么要做这个事情，才会主动掌控时间，才能优化、管理好自己的时间。

Go 的 bytes 包，其功能是处理 `[]byte` 类型实例的相关函数；其功能和 strings 包类似，可以做类比（analogous）。

bytes.Buffer 封装的类型，相当于是**可变长度的字节缓冲区**。该类型的零值，相当于缓冲区为空，可以“拿来即用”。

~~~go
func TestBufferInit(t *testing.T) {
	buffer := &bytes.Buffer{}
	fmt.Printf("%d, %d\n", buffer.Len(), buffer.Cap())

	buffer.Grow(10)
	fmt.Printf("%d, %d\n", buffer.Len(), buffer.Cap())
}
0, 0
0, 64
~~~

另外，bytes.Buffer 实现了底层 []byte 的自动扩容，也就是当其容量无法满足要求时，会自动扩容。

bytes.Buffer 类型实现了如下接口：

1. io.Reader：**可被读取**的能力；
2. io.Writer：**可被写入**的能力；
3. io.ReadWriter：集成了 io.Reader 和 io.Writer，内嵌入了上述 2 个接口；
4. io.ReaderFrom：从指定的 io.Reader 中**读取数据的能力**；
5. io.WriterTo：**向**指定的 io.Writer **写入数据的能力**；
6. io.ByteReader：从缓冲区中读取 byte 的能力；
7. io.ByteScanner：内嵌入了 io.ByteReader 接口，新增了 UnreadByte() 方法，相当于是一个 byte 的 Scanner；
8. io.ByteWriter：向缓冲区写入指定的 byte 值；
9. io.RuneReader：从缓冲区中读取单个 `UTF-8 encoded Unicode character`；
10. io.RuneScanner：内嵌入了 io.RuneReader 接口，新增了 UnreadRune() 方法，相当于是一个 Rune 的 Scanner；
11. io.StringWriter：向缓冲区中写入 string 值。

下面针对 bytes.Buffer 的能力，做详细的论述。

# 1 创建 bytes.Buffer 实例

bytes.Buffer 实际是一个结构体类型：

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

bytes.Buffer 提供了 2 种创建实例的函数：

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

可以分别使用 string 和 []byte 创建 *Buffer 类型值。

~~~go
func TestBuffer(t *testing.T) {
	origin := []byte("foo")
	var buffer *bytes.Buffer = bytes.NewBuffer(origin)
	fmt.Printf(" %s\n", buffer.String())

	origin[0] = 'o'
	fmt.Println(string(origin))
	fmt.Printf(" %s\n", buffer.String())
}
 foo
ooo
 ooo
~~~

正如注释中所述，NewBuffer 函数的入参 buf，是为 `*Buffer` 所有的，在调用该函数之后，**调用者就不能再使用（修改）buf 的内容**，否则会导致 `*Buffer` 的改动。

# 2 io.Reader

io.Reader 接口封装的是：**可被读取**的能力。也就是从 bytes.Buffer 中读取字节序列到指定的 p 中

~~~go
// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}
~~~

注释中指出，如果 n > 0，表示已经从 bytes.Buffer 中读取到了字节，需要在考虑处理 err 之前（不管 err 是否是 nil），先行对这 n 个字节进行处理。关于 io.Reader 的注释，是对所有实现了 io.Reader 接口的类型说明的。

~~~go
package gobytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {
	origin := []byte("foo")
	var buffer *bytes.Buffer = bytes.NewBuffer(origin)
	fmt.Printf(" %s\n", buffer.String())

	value := make([]byte, 2)
	n, err := buffer.Read(value)
	fmt.Printf("%d, %s, %s\n", n, value, buffer.String())
	check(err, t)

	origin[2] = 'T'
	fmt.Println(string(origin))
	fmt.Printf(" %s\n", buffer.String())
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}
 foo
2, fo, o
foT
 T
~~~

# 3 io.ByteReader

~~~go
// ByteReader is the interface that wraps the ReadByte method.
//
// ReadByte reads and returns the next byte from the input or
// any error encountered. If ReadByte returns an error, no input
// byte was consumed, and the returned byte value is undefined.
//
// ReadByte provides an efficient interface for byte-at-time
// processing. A Reader that does not implement  ByteReader
// can be wrapped using bufio.NewReader to add this method.
type ByteReader interface {
	ReadByte() (byte, error)
}
~~~

很有意思的结论：io.ByteReader 接口中，实际上封装的是 ReadByte()，恰好和其名称相反。ReadByte() 表达了“读取一个 byte 的能力”。

~~~go
func TestBufferByteReader(t *testing.T) {
	buffer := bytes.NewBufferString("Michoi")
	char, _ := buffer.ReadByte()
	fmt.Printf("%q\n", char)
}
'M'
~~~

# 4 io.ByteScanner

~~~go
// ByteScanner is the interface that adds the UnreadByte method to the
// basic ReadByte method.
//
// UnreadByte causes the next call to ReadByte to return the same byte
// as the previous call to ReadByte.
// It may be an error to call UnreadByte twice without an intervening
// call to ReadByte.
type ByteScanner interface {
	ByteReader
	UnreadByte() error
}
~~~

和 io.ByteReader 不同的是，io.ByteScanner 还包含了 UnreadByte()，表示让下一次调用 readByte() 返回结果和上一次相同。

~~~go
func TestBufferByteScanner(t *testing.T) {
	buffer := bytes.NewBufferString("Michoi")
	char, _ := buffer.ReadByte()
	fmt.Printf("%q\n", char)

	buffer.UnreadByte()

	char, _ = buffer.ReadByte()
	fmt.Printf("%q\n", char)
}
'M'
'M'
~~~

# 5 io.RuneReader

~~~go
// RuneReader is the interface that wraps the ReadRune method.
//
// ReadRune reads a single UTF-8 encoded Unicode character
// and returns the rune and its size in bytes. If no character is
// available, err will be set.
type RuneReader interface {
	ReadRune() (r rune, size int, err error)
}
~~~

和 io.ByteReader 类似，不同之处在于将读取的单位从 byte 更换成了 rune。

~~~go
func TestBufferRuneReader(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	r, n, _ := buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)
}
'中', 3
~~~

# 6 io.RuneScanner

~~~go
// RuneScanner is the interface that adds the UnreadRune method to the
// basic ReadRune method.
//
// UnreadRune causes the next call to ReadRune to return the same rune
// as the previous call to ReadRune.
// It may be an error to call UnreadRune twice without an intervening
// call to ReadRune.
type RuneScanner interface {
	RuneReader
	UnreadRune() error
}
~~~

和 io.ByteScanner 类似，不同之处在于将读取的单位从 byte 更换成了 rune。

~~~go
func TestBufferRuneScanner(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	r, n, _ := buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)

	buffer.UnreadRune()

	r, n, _ = buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)
}
'中', 3
'中', 3
~~~

# 7 io.Writer

~~~go
// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
type Writer interface {
	Write(p []byte) (n int, err error)
}
~~~

io.Writer 封装的是**一种可被写入的能力**，其中入参 p 相当于是数据源，是一个 []byte 类型。bytes.Buffer 调用该方法时，会将 p 写入到 Buffer 中。

~~~go
func TestBufferWriter(t *testing.T) {
	var buffer *bytes.Buffer = &bytes.Buffer{}

	content := "Michoi"
	n, err := buffer.Write([]byte(content))
	if n < len(content) {
		t.Fatal(err.Error())
	}
	fmt.Printf("%d, %s, %s\n", n, content, buffer.String())
}
6, Michoi, Michoi
~~~

# 8 io.ByteWriter

~~~go
// ByteWriter is the interface that wraps the WriteByte method.
type ByteWriter interface {
	WriteByte(c byte) error
}
~~~

io.ByteWriter 接口，能够像 Buffer 中写入一个 byte 的值。很有意思的是，其中包含的就是 WriteByte()：

~~~go
func TestBufferByteWriter(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	err := buffer.WriteByte('\x41')
	if err == nil {
		fmt.Println(buffer.String())
	}
}
中Michoi国A
~~~

# 9 io.StringWriter

~~~go
// StringWriter is the interface that wraps the WriteString method.
type StringWriter interface {
	WriteString(s string) (n int, err error)
}
~~~

io.StringWriter 和 io.ByteWriter 类似，不同之处是将 byte 替换成了 string 值：

~~~go
func TestBufferStringWriter(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	buffer.WriteString("人，在中国")
	fmt.Println(buffer.String())
}
中Michoi国人，在中国
~~~

# 10 io.ReadWriter

~~~go
// ReadWriter is the interface that groups the basic Read and Write methods.
type ReadWriter interface {
	Reader
	Writer
}
~~~

集成了 io.Reader 和 io.Writer，内嵌入了 io.Reader 和 io.Writer 接口。

# 11 io.ReaderFrom

~~~go
// ReaderFrom is the interface that wraps the ReadFrom method.
//
// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except io.EOF encountered during the read is also returned.
//
// The Copy function uses ReaderFrom if available.
type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}
~~~

io.ReaderFrom 从指定的 io.Reader 中读取数据的能力，实际上和 io.Writer 接口有类似的地方，只是 io.ReaderFrom 指定输入源是 io.Reader。

~~~go
func TestBufferReaderFrom(t *testing.T) {
	buffer := &bytes.Buffer{}
	reader := bytes.NewReader([]byte("Michoi"))

	n, _ := buffer.ReadFrom(reader)
	fmt.Printf("%d, %s\n", n, buffer.String())
}
6, Michoi
~~~

# 12 io.WriterTo

~~~go
// WriterTo is the interface that wraps the WriteTo method.
//
// WriteTo writes data to w until there's no more data to write or
// when an error occurs. The return value n is the number of bytes
// written. Any error encountered during the write is also returned.
//
// The Copy function uses WriterTo if available.
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}
~~~

io.WriterTo：向指定的 io.Writer 写入数据的能力。这个接口的功能和 io.Reader 类似，都是将 Buffer 中的内容写入目的地，不同之处是，此处指定的是一个 io.Writer 实例。

~~~go
func TestBufferWriterTo(t *testing.T) {
	buffer := bytes.NewBufferString("Michio")
	writer := &bytes.Buffer{}
	n, _ := buffer.WriteTo(writer)
	fmt.Printf("%d, %s\n", n, writer.String())
}
6, Michio
~~~