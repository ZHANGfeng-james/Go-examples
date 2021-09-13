SectionReader 实现了 Read、Seek 和 ReadAt 方法，其底层都是依赖于其封装的 ReaderAt 字段（本身是一个接口）。

~~~go
// SectionReader implements Read, Seek, and ReadAt on a section
// of an underlying ReaderAt.
type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
}
~~~

其中 ReaderAt 接口定义：

~~~go
// ReaderAt is the interface that wraps the basic ReadAt method.
//
// ReadAt reads len(p) bytes into p starting at offset off in the
// underlying input source. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered.
//
// When ReadAt returns n < len(p), it returns a non-nil error
// explaining why more bytes were not returned. In this respect,
// ReadAt is stricter than Read.
//
// Even if ReadAt returns n < len(p), it may use all of p as scratch
// space during the call. If some data is available but not len(p) bytes,
// ReadAt blocks until either all the data is available or an error occurs.
// In this respect ReadAt is different from Read.
//
// If the n = len(p) bytes returned by ReadAt are at the end of the
// input source, ReadAt may return either err == EOF or err == nil.
//
// If ReadAt is reading from an input source with a seek offset,
// ReadAt should not affect nor be affected by the underlying
// seek offset.
//
// Clients of ReadAt can execute parallel ReadAt calls on the
// same input source.
//
// Implementations must not retain p.
type ReaderAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
}
~~~

ReaderAt 接口和 Reader 接口的不同是：其 ReadAt 和 Read 方法有区别，前者包含了一个 off 入参，表示起始读取点偏移量。

**创建 SectionReader 实例**的函数：

~~~go
// NewSectionReader returns a SectionReader that reads from r
// starting at offset off and stops with EOF after n bytes.
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {
	return &SectionReader{r, off, off, off + n}
}
~~~

其首个入参是一个实现了 ReaderAt 接口的实例。这个 *SectionReader 实例，表示从 ReaderAt 实例距离首字节位置 off 的地方开始读取，依次读取 n 个字节，此时意味着输出流的结束，也就是返回的是 io.EOF。

这里有个应用实例：

~~~go
package main

import (
	"io"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("./A20_Apex.apk")
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := io.NewSectionReader(file, 0, fileinfo.Size())
	reader := sectionReadCloser{r}
	defer reader.Close()

	out, err := os.Create("./copy.apk")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()
	// 将 src 拷贝到 out 中
	io.Copy(out, reader)
}

type sectionReadCloser struct {
	*io.SectionReader
}

func (rc sectionReadCloser) Close() error {
	return nil
}
~~~

