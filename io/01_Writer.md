# 1 Writer

Writer 接口，包装的是基本的 Writer 方法。实现这个接口的类型，其含义就是该类型具备有**被写的能力**。

Writer 接口，对应的是 **Input 属性**，也就是**将外部的数据写入到实现了 Writer 接口的实例中**。

Writer 接口相关的 Write 方法，从 p 中读取 len(p) 的字节序列，写入到底层的数据流（实现 Writer 接口的实例）中。该方法返回 `n <= n <= len(p)`，以及与之相关的 error 实例（其表现是：导致 write 动作提前结束）。

write 方法在 `n < len(p)` 时，必然会返回一个 non-nil 的 error 实例。

~~~go
type Writer interface {
    Write(p []byte) (n int, err error)
}
~~~

在 io 包中有一个特殊的 Writer 接口实现：

~~~go
type discard struct{}

// Discard is an Writer on which all Write calls succeed
// without doing anything.
var Discard Writer = discard{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}
~~~

可以看到，Discard 实例实现了 Writer 接口，但是其对应方法没有做任何改变。

# 2 MultiWriter

io.MultiWriter 是一个函数：

~~~go
func MultiWriter(writers ...Writer) Writer
~~~

该函数由多个 Writer 创建单一的 Writer，此时写入到该返回的 Writer 的内容会同时写入到这多个 Writer 中，类似 Unix 中的 tee(1) 指令。如果中间某处 Write 操作出现错误，则 Write 动作**不会继续下去**：

~~~go
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

func main() {
	r := strings.NewReader("michoi")

	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)

	if _, err := io.Copy(w, r); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf1.String(), buf2.String())
}
michoi michoi
~~~

