

bytes 包，我们回到**最初定义**的地方：

Package bytes implements functions for the manipulation of byte slices. It is analogous to the facilities of the strings package.

bytes 包具备用来**操作 []byte 类型值的方法**，和 strings 包的功能类似，可以当做是一个**工具包**使用。

bytes 包中有**两颗“闪亮的星”**：Reader 和 Buffer

1. bytes.Reader：在构造了 Reader 后，可通过各种类似 Reader 的方式读取其中的字节或字节序列；
2. bytes.Buffer：是一种 byte 的缓冲区。











* `func copy(dst, src []Type) int`：copy **内建函数**是用来实现 `[]Type` 的拷贝，其中 Type 代表（stand-in）的是 Go 中的任意一种数据类型。
* `func WriteString(w Writer, s string) (n int, err error)`：是 io 包下的函数，用于向 w 写入 s，其中 s 也可以是 []byte。
* `fmt.Println("%q\n", []byte("Michoi"))`：将 []byte 内容以字符串的形式输出。

