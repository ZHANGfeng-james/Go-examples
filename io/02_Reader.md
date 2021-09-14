Reader 接口封装了基本的 Read 方法：

~~~go
type Reader interface {
    Read(p []byte) (n int, err error)
}
~~~

Reader **具备 Output 属性**，也就是可以将自身底层流**输出**到对应的**数据载体**上去。

Read 方法会读取 len(p) 的字节内容，并填充到 p 的字节切片中。该方法返回的是读取到的字节长度，以及与之相关的 error 实例。特别的 n 的取值范围是：`0 <= n <= len(p)`。即便 `n < len(p)`，也会在调用 Read 方法时使用整个 p 字节切片。如果 Reader 实例在调用 Read 方法时，仅有部分数据可用（可读），该方法会立即返回，而不会等待。

当 Read 方法在读取了 n （n >0）个字节时，遇到 error 或者 io.EOF（Reader 具备被读的属性，因此会遇到 io.EOF）时，会返回读取到的字节数。**本次调用**会返回一个 non-nil 的 error，而**紧接着的一次 Read 方法调用**后，会返回一个 error 且 n = 0。比如：在一个 Reader 实例读取到输入流的末尾时，会返回读取到的字节数，以及 io.EOF 或者 nil 的错误值，而**下一次 Read** 会返回 0 和 io.EOF。

Read 方法的调用者需要在考虑 error 之前，处理已经读取到的 n 个字节。这是一种正确处理 IO 错误的方式，另外还需要特别处理可被允许的 io.EOF 错误。

实现 Read 方式时，最好不要在返回 0 的同时返回一个 nil 的 error 值，除非 len(p) 的值是 0。对于调用者来说，如果出现返回值是 0 和 nil，则表示什么都没有发生，而且不是 io.EOF（也就是没有到达 input 的末尾）。

Read 方法的实现者，不应该持有 p 字节切片。

# 1 Reader

