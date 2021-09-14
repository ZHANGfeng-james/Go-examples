

1. context 是什么
2. context 在 Go 中有什么用途？为什么要引入 context？
3. 应用场景是什么？
4. 在使用时，有哪些特别需要注意的坑？
5. Context：上下文定义、可以 cancel 的上下文、带时间的上下文、带 KV 的上下文、多线程上下文控制、上下文树、关闭上下文



Go 标准库的 context 包定义了 Context 类型，这个类型的变量会“携带”最后期限、取消 signal 和一些穿梭在 API 边界（似乎和 API 调用时的栈帧相关）和调用处理中的与“请求区间”（似乎和 request 网络请求相关）有关的值。

Server 接收到一个到来的 request 时（似乎和 HTTP 服务器相关），会创建一个 Context 类型值，紧接着的一系列调用都应该接受一个 Context 类型值。于此相关的**方法调用链**必需要能传递这个 Context 类型值，另外还可以**使用 WithCancel、WithDeadline、WithTimeout 或者是 WithValue 创建一个衍生出的 Context 作为可选项**替换原先的 Context 类型值。当一个 Context 值被取消，所有由该 Context 值衍生出的 Context 值都会被取消。

使用 WithCancel、WithDeadline 和 WithTimeout 函数能够在 Context（父辈）的基础上衍生出一个 Context 值（孩子），并携带一个 CancelFunc 类型值。调用 CancelFunc 时，会**取消**其孩子以及子孙的**执行**，**移除**孩子和父辈之间的**引用关系**，并停止与之关联的任何**计时器**。如果没有调用 CancelFunc 会导致**泄漏**其孩子及其子孙，直到父辈的 Context 类型值被取消，或者**计时时间到**为止。go vet 工具会检测 CancelFunc 类型值是否用在**控制流路径**上。

使用 Context 的程序需要遵循如下的规则在包之间保持**接口兼容**，同时让 go **静态分析工具**能有效检查 Context 的衍生关系：

* 不要在一个结构体类型中保有 Context 变量值，取而代之的是：显式地为每个需要 Context 类型值的函数传递一个参数，且应该是首个参数值。

~~~go
func DoSometing(ctx context.Context, arg Arg) error {
    // ... use ctx ...
}
~~~

* 不要传递一个 nil 值给 Context 类型变量，即使函数允许这样做。当不确定是否需要使用时，可传递 context.TODO 值。
* 仅仅是在和 HTTP 请求区域（范围）相关的处理和 API 调用中使用 Context Values，而不是在其他函数中为了传递额外的参数使用 Context 类型值。

相同的 Context 类型值可能会被传递到运行在不同 goroutine 的函数中，这种场景中，Context 的使用是并发安全的。





context 如何被取消



context.Value 的查找过程是怎样的