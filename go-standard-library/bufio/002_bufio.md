# Reader

bufio 包实现的是**缓冲式 I/O 功能**，其中的类型实体封装 io.Reader 和 io.Writer 并创建了新的类型 Reader 和 Writer。新的类型实现了 io.Reader 和 io.Writer 接口，并提供了缓冲式的 I/O 功能。

下面仅分析 bufio.Reader 的源代码：

~~~go
// Reader implements buffering for an io.Reader object.
type Reader struct {
	buf          []byte
	rd           io.Reader // reader provided by the client
	r, w         int       // buf read and write positions
	err          error
	lastByte     int // last byte read for UnreadByte; -1 means invalid
	lastRuneSize int // size of last rune read for UnreadRune; -1 means invalid
}
~~~

从类型结构体的域值可看出，**关键缓冲功能是在 buf 字段实现的**，rd 字段相当于提供了**数据源**。另外 2 个 int 类型的值 r 和 w 用来标记 buf 中读写位置的索引，其中 [r, w) 是可读的部分，[w, len(buf)-1] 的部分是可写入的部分。上面说到的**读位置索引**，是指从 buf 中取数据；**写位置索引**，是指从 rd 中取数据，并写到 buf 中。读位置索引和写位置索引都是作用在 buf 上的！

~~~go
// NewReaderSize returns a new Reader whose buffer has at least the specified
// size. If the argument io.Reader is already a Reader with large enough
// size, it returns the underlying Reader.
func NewReaderSize(rd io.Reader, size int) *Reader {
	// Is it already a Reader?
	b, ok := rd.(*Reader) // bufio.Reader
	if ok && len(b.buf) >= size {
		return b
	}
	if size < minReadBufferSize {
		size = minReadBufferSize
	}
	r := new(Reader)
	r.reset(make([]byte, size), rd)
	return r
}

// NewReader returns a new Reader whose buffer has the default size.
func NewReader(rd io.Reader) *Reader {
	return NewReaderSize(rd, defaultBufSize)
}

func (b *Reader) reset(buf []byte, r io.Reader) {
	*b = Reader{
		buf:          buf,
		rd:           r,
		lastByte:     -1,
		lastRuneSize: -1,
	}
}

const (
	defaultBufSize = 4096 // 默认缓冲区大小是 4096 个字节，4KB
)

const minReadBufferSize = 16 // 最小缓冲区大小是 16 个字节
~~~

上述使用 io.Reader 创建 bufio.Reader 类型的值，并返回了该值的指针，其中在适当的时候根据不同的 size 大小调整了 buf 的大小。

~~~go
// Size returns the size of the underlying buffer in bytes.
func (b *Reader) Size() int { return len(b.buf) }
~~~

实际是返回缓冲区的 len(b.buf) 值，是一个切片的 length 值。

~~~go
// Reset discards any buffered data, resets all state, and switches
// the buffered reader to read from r.
func (b *Reader) Reset(r io.Reader) {
	b.reset(b.buf, r)
}
~~~

有点==疑惑==是此处 b.buf 的设置，重复使用了原先的 b.buf 值。

~~~go
// Buffered returns the number of bytes that can be read from the current buffer.
func (b *Reader) Buffered() int { return b.w - b.r }
~~~

返回从当前的 buf 中能读取到的字节数量。

从==最简单==的 ReadByte() (byte, error) 开始分析：

~~~go
// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (b *Reader) ReadByte() (byte, error) {
	b.lastRuneSize = -1
	for b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		b.fill() // buffer is empty
	}
	c := b.buf[b.r]
	b.r++
	b.lastByte = int(c)
	return c, nil
}

// fill reads a new chunk into the buffer.
func (b *Reader) fill() {
	// Slide existing data to beginning.
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		panic("bufio: tried to fill full buffer")
	}

	// Read new data: try a limited number of times.
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.rd.Read(b.buf[b.w:])
		if n < 0 {
			panic(errNegativeRead)
		}
		b.w += n
		if err != nil {
			b.err = err
			return
		}
		if n > 0 {
			return
		}
	}
	b.err = io.ErrNoProgress
}

const maxConsecutiveEmptyReads = 100
~~~

如果当前的 buf 中没有任何数据，用 b.r == b.w 进行判断，也就是读标记位和写标记位值相等，则需要向 buf 中填充数据。调用 fill 私有方法填充数据，从 rd 中读取数据写到 buf 中，同时更新 w 值（标记了写 buf 的位置索引）。填充数据后，b.buf[b.r] 值就是待读取的当前 byte 值，并将该值保存到 lastByte 中。

~~~go
// ReadRune reads a single UTF-8 encoded Unicode character and returns the
// rune and its size in bytes. If the encoded rune is invalid, it consumes one byte
// and returns unicode.ReplacementChar (U+FFFD) with a size of 1.
func (b *Reader) ReadRune() (r rune, size int, err error) {
	for b.r+utf8.UTFMax > b.w 
		&& !utf8.FullRune(b.buf[b.r:b.w]) 
		&& b.err == nil && b.w-b.r < len(b.buf) {
		b.fill() // b.w-b.r < len(buf) => buffer is not full
	}
	b.lastRuneSize = -1
	if b.r == b.w {
		return 0, 0, b.readErr()
	}
	r, size = rune(b.buf[b.r]), 1
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRune(b.buf[b.r:b.w])
	}
	b.r += size
	b.lastByte = int(b.buf[b.r-1])
	b.lastRuneSize = size
	return r, size, nil
}
~~~

比较复杂是需要判断何时需要重新填充 buf，填充完后，首先读 1 个字节，并判断是否表示的是多个字节的 Unicode code point。如果是，则从 b.r 到 b.w 之间读取一个 rune，并更新 b.r 值。

~~~go
// ReadSlice reads until the first occurrence of delim in the input,
// returning a slice pointing at the bytes in the buffer.
// The bytes stop being valid at the next read.
// If ReadSlice encounters an error before finding a delimiter,
// it returns all the data in the buffer and the error itself (often io.EOF).
// ReadSlice fails with error ErrBufferFull if the buffer fills without a delim.
// Because the data returned from ReadSlice will be overwritten
// by the next I/O operation, most clients should use
// ReadBytes or ReadString instead.
// ReadSlice returns err != nil if and only if line does not end in delim.
func (b *Reader) ReadSlice(delim byte) (line []byte, err error) {
	s := 0 // search start index
	for {
		// Search buffer.
		if i := bytes.IndexByte(b.buf[b.r+s:b.w], delim); i >= 0 {
			i += s
			line = b.buf[b.r : b.r+i+1]
			b.r += i + 1
			break
		}

		// Pending error?
		if b.err != nil {
			line = b.buf[b.r:b.w]
			b.r = b.w
			err = b.readErr()
			break
		}

		// Buffer full?
		if b.Buffered() >= len(b.buf) {
			b.r = b.w
			line = b.buf
			err = ErrBufferFull
			break
		}

		s = b.w - b.r // do not rescan area we scanned before

		b.fill() // buffer is not full
	}

	// Handle last byte, if any.
	if i := len(line) - 1; i >= 0 {
		b.lastByte = int(line[i])
		b.lastRuneSize = -1
	}

	return
}
~~~

ReadSlice 在 buf 中查找指定 delimiter 的切片。以当前 buf 内容为开始查找，若未找到，则会持续填充并查找，==直到达到 Buffer full 状态==。此时如果仍在 full buf 中未找到，则会返回 ErrBufferFull 的错误值，同时返回的是当前 buf 中的所有内容，而且此时 b.r 被赋值为 b.w，意味着这部分内容已被读取。

~~~go
// collectFragments reads until the first occurrence of delim in the input. It
// returns (slice of full buffers, remaining bytes before delim, total number
// of bytes in the combined first two elements, error).
// The complete result is equal to
// `bytes.Join(append(fullBuffers, finalFragment), nil)`, which has a
// length of `totalLen`. The result is strucured in this way to allow callers
// to minimize allocations and copies.
func (b *Reader) collectFragments(delim byte) (fullBuffers [][]byte, finalFragment []byte, totalLen int, err error) {
	var frag []byte
	// Use ReadSlice to look for delim, accumulating full buffers.
	for {
		var e error
		frag, e = b.ReadSlice(delim)
		if e == nil { // got final fragment
			break
		}
		if e != ErrBufferFull { // unexpected error 比如 io.EOF 等
			err = e
			break
		}

		// Make a copy of the buffer.
		buf := make([]byte, len(frag))
		copy(buf, frag)
		fullBuffers = append(fullBuffers, buf) // fullBuffers 保存已排查的所有内容
		totalLen += len(buf)
	}

	totalLen += len(frag)
	return fullBuffers, frag, totalLen, err // frag 找到的内容
}
~~~

collectFragments 会不断调用 RaedSlice 知道找到指定 delimiter 或者出现 error 不为 ErrBufferFull 的情况；又因为 ReadSlice 会不断调用 fill() 直到本次 full buffer 状态才返回，因此 fullBuffers 会返回所有已排查的内容。

~~~go
// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
// For simple uses, a Scanner may be more convenient.
func (b *Reader) ReadBytes(delim byte) ([]byte, error) {
	full, frag, n, err := b.collectFragments(delim)
	// Allocate new buffer to hold the full pieces and the fragment.
	buf := make([]byte, n)
	n = 0
	// Copy full pieces and fragment in.
	for i := range full {
		n += copy(buf[n:], full[i])
	}
	copy(buf[n:], frag)
	return buf, err
}
~~~

底层调用了 collectFragments(delim) 方法，并在 ReadBytes 返回了所有已排查的内容，其中将 `fullBuffers [][]byte` 转化成了 []byte 类型。

~~~go
// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end in
// delim.
// For simple uses, a Scanner may be more convenient.
func (b *Reader) ReadString(delim byte) (string, error) {
	full, frag, n, err := b.collectFragments(delim)
	// Allocate new buffer to hold the full pieces and the fragment.
	var buf strings.Builder
	buf.Grow(n)
	// Copy full pieces and fragment in.
	for _, fb := range full {
		buf.Write(fb)
	}
	buf.Write(frag)
	return buf.String(), err
}
~~~

实际上是和 ReadByte(delim byte) 相同的方法，不同之处仅在于返回值是 string 类型。

~~~go
func (b *Reader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// To read exactly len(p) bytes, use io.ReadFull(b, p).
// At EOF, the count will be zero and err will be io.EOF.
func (b *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, b.err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			if n > 0 {
				b.lastByte = int(p[n-1])
				b.lastRuneSize = -1
			}
			return n, b.readErr()
		}
		// One read.
		// Do not use b.fill, which will loop.
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	b.lastByte = int(b.buf[b.r-1])
	b.lastRuneSize = -1
	return n, nil
}
~~~

Read(p []byte) 中有一个段注释：Do not use b.fill, which will loop. 也就是说，上述方法仅会触发一次 b.rd.Read(b.buf)。



# Writer

bufio.Writer 封装了 io.Writer 实例，表示的是一种带有==缓冲==功能的 I/O ==写入器==。

带有缓冲功能的目的是为了提升写入的性能，如何提升的？通过内置 []byte 切片，将写入的内容先缓存到 []byte 中，==当时机成熟时==，==一次性==写入底层的 io.Writer 值中。

~~~go
// Writer implements buffering for an io.Writer object.
// If an error occurs writing to a Writer, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
type Writer struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}
~~~

bufio.Writer 结构体类型中封装的 []byte 就是用于==缓冲功能的切片，提升 Output 的性能==。

~~~go
// NewWriterSize returns a new Writer whose buffer has at least the specified
// size. If the argument io.Writer is already a Writer with large enough
// size, it returns the underlying Writer.
func NewWriterSize(w io.Writer, size int) *Writer {
	// Is it already a Writer?
	b, ok := w.(*Writer) // bufio.Writer
	if ok && len(b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}
	return &Writer{
		buf: make([]byte, size),
		wr:  w,
	}
}

// NewWriter returns a new Writer whose buffer has the default size.
func NewWriter(w io.Writer) *Writer {
	return NewWriterSize(w, defaultBufSize)
}
~~~

和 bufio.Reader 类似，默认的底层缓冲区大小是 4KB。

~~~go
// Available returns how many bytes are unused in the buffer.
func (b *Writer) Available() int { return len(b.buf) - b.n }

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *Writer) Buffered() int { return b.n }

// Size returns the size of the underlying buffer in bytes.
func (b *Writer) Size() int { return len(b.buf) }
~~~

b.n 值代表的含义是：已经写入到 buf 中的字节数（已缓存的字节数）；Available() 返回的是还能够填充到 buf 中的字节数，也就是还未使用的空间数量。

~~~go
// Flush writes any buffered data to the underlying io.Writer.
func (b *Writer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
            // 存在 b.n - n 个字节的内容还未写入到 io.Writer 中
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}
~~~

Flush() 是一次性将 b.n 个字节（已缓存的）全部写入到底层的 io.Writer 中，不过，当然这个操作可能是存在异常的。上述的代码兼容了异常的情况。

~~~go
// WriteByte writes a single byte.
func (b *Writer) WriteByte(c byte) error {
	if b.err != nil {
		return b.err
	}
	if b.Available() <= 0 && b.Flush() != nil {
		return b.err
	}
	b.buf[b.n] = c
	b.n++
	return nil
}
~~~

如果 b.Available() <=0，表示 buf 中已不存在缓冲空间，会启动 b.Flush() 将缓冲区内容真实地写入到 io.Writer 中。

~~~go
// WriteString writes a string.
// It returns the number of bytes written.
// If the count is less than len(s), it also returns an error explaining
// why the write is short.
func (b *Writer) WriteString(s string) (int, error) {
	nn := 0
	for len(s) > b.Available() && b.err == nil {
		n := copy(b.buf[b.n:], s)
		b.n += n
		nn += n
		s = s[n:]
		b.Flush()
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], s)
	b.n += n
	nn += n
	return nn, nil
}
~~~

for 循环中分片段将 s 各个字节写入到 buf 中，在适当的时候触发 Flush()。

~~~go
// WriteRune writes a single Unicode code point, returning
// the number of bytes written and any error.
func (b *Writer) WriteRune(r rune) (size int, err error) {
	if r < utf8.RuneSelf {
		err = b.WriteByte(byte(r))
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	if b.err != nil {
		return 0, b.err
	}
	n := b.Available()
	if n < utf8.UTFMax {
		if b.Flush(); b.err != nil {
			return 0, b.err
		}
		n = b.Available()
		if n < utf8.UTFMax {
			// Can only happen if buffer is silly small.
			return b.WriteString(string(r))
		}
	}
	size = utf8.EncodeRune(b.buf[b.n:], r)
	b.n += size
	return size, nil
}
~~~

如果待写入的 rune 是单字节的，则直接调用 WriteByte 方法，否则可能会调用 Flush()。

~~~go
// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *Writer) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}
~~~

将 []byte 依次写入到 buf 中，并可能触发 Flush()。

~~~go
// ReadFrom implements io.ReaderFrom. If the underlying writer
// supports the ReadFrom method, and b has no buffered data yet,
// this calls the underlying ReadFrom without buffering.
func (b *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	if b.err != nil {
		return 0, b.err
	}
	if b.Buffered() == 0 {
		if w, ok := b.wr.(io.ReaderFrom); ok {
			n, err = w.ReadFrom(r)
			b.err = err
			return n, err
		}
	}
	var m int
	for {
		if b.Available() == 0 {
			if err1 := b.Flush(); err1 != nil {
				return n, err1
			}
		}
		nr := 0
		for nr < maxConsecutiveEmptyReads {
			m, err = r.Read(b.buf[b.n:])
			if m != 0 || err != nil {
				break
			}
			nr++
		}
		if nr == maxConsecutiveEmptyReads {
			return n, io.ErrNoProgress
		}
		b.n += m
		n += int64(m)
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		// If we filled the buffer exactly, flush preemptively.
		if b.Available() == 0 {
			err = b.Flush()
		} else {
			err = nil
		}
	}
	return n, err
}
~~~



比如，通过 HTTP 请求下载文件：

~~~go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var url string

func init() {
	flag.StringVar(&url, "url", "", "Download Resource URI")
	fmt.Println(url)
}

func main() {
	flag.Parse()

	if url == "" {
		log.Fatal("Download Resource URI error!")
	}

	req, err := http.NewRequest("GET", url, nil)
	checkError(err)
	res, err := http.DefaultClient.Do(req)
	checkError(err)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	formatTime := time.Now().Format("2006-01-02_15-04-05")
	file, err := os.Create("Snipaste_" + formatTime + ".png")
	checkError(err)
	defer file.Close()

	buf := bufio.NewWriter(file) // 通过 Writer 写入到文件中
	buf.Write(body)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
~~~









# Scanner





https://segmentfault.com/a/1190000013493942

