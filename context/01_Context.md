和 Context 相关的**疑惑**：

1. context 是什么？

2. context 在 Go 中有什么用途？为什么要引入 context？

3. WithCancel、WithDeadline、WithTimeout 的应用场景分别是什么？

4. 在使用时，有哪些特别需要注意的坑？

5. Context：上下文定义、可以 cancel 的上下文、带时间的上下文、带 KV 的上下文、多线程上下文控制、上下文树、关闭上下文

6. context 如何被取消？

7. cancelFunc 和 Done 之间是什么关系？

   cancelFunc 的调用或执行，会导致 Done channel 关闭。但这——cancelFunc 只是 Done channel 关闭的引发因素。针对不同的情况，还有其他因素。WithCancel：

8. context.Value 的查找过程是怎样的



Go 标准库的 context 包定义了 Context 类型，这个类型的变量会“携带”最后期限、取消 signal 和一些穿梭在 API 边界（似乎和 API 调用时的栈帧相关）和调用处理中的与“请求区间”（似乎和 request 网络请求相关）有关的值。

Server 接收到一个到来的 request 时（似乎和 HTTP 服务器相关），会创建一个 Context 类型值，紧接着的一系列调用都应该接受一个 Context 类型值。于此相关的**方法调用链**必需要能传递这个 Context 类型值，另外还可以**使用 WithCancel、WithDeadline、WithTimeout 或者是 WithValue 创建一个衍生出的 Context 作为可选项**替换原先的 Context 类型值。当一个 Context 值被取消，所有由该 Context 值衍生出的 Context 值都会被取消。

使用 WithCancel、WithDeadline 和 WithTimeout 函数能够在 Context（父辈）的基础上衍生出一个 Context 值（孩子），并携带一个 CancelFunc 类型值。<u>调用 CancelFunc 时，会**取消**其孩子以及子孙的**执行**，**移除**孩子和父辈之间的**引用关系**，并停止与之关联的任何**计时器**。如果没有调用 CancelFunc 会导致**泄漏**其孩子及其子孙，直到父辈的 Context 类型值被取消，或者**计时时间到**为止。</u>go vet 工具会检测 CancelFunc 类型值是否用在**控制流路径**上。

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

### Variables

Canceled 是一个 error 类型的值。当 Context 类型值被取消时，在 Context.Err 方法中返回该值。

~~~go
var Canceled = errors.New("context canceled")
~~~

与之对应，DeadlineExceeded 也是一个 error 类型的值，当 Context 类型值抵达期限时 deadline，可通过 Context.Err 方法获得：

~~~go
// DeadlineExceeded is the error returned by Context.Err when the context's
// deadline passes.
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }
~~~

### func WithCancel

~~~go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// A CancelFunc tells an operation to abandon its work.
// A CancelFunc does not wait for the work to stop.
// A CancelFunc may be called by multiple goroutines simultaneously.
// After the first call, subsequent calls to a CancelFunc do nothing.
type CancelFunc func()

// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     chan struct{}         // created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}
~~~

函数 WithCancel 的返回值有 2 个，分别是一份对 parent 的拷贝，同时附带了一个新建的 Done Channel。这个返回的 Done Channle 会在如下情形会**关闭 Close**：

1. 调用了 cancel
2. parent 的 Context 类型值的 Done 关闭

上述无论哪个情况发生先发生，都会让这个 Done Channel 关闭。

**取消**（动词，表示调用这个 Context 的 cancel 函数）这个 Context 会**释放与之相关的系统资源**（否则会导致泄漏），因此，在程序中应该尽可能要在 Context 结束时调用 cancel 函数（避免资源泄漏）。

比如下述示例程序：

~~~go
package main

import (
	"context"
	"fmt"
)

func main() {
	gen := func(ctx context.Context) <-chan int {
		origin := make(chan int)
		n := 0

		go func() {
			for {
				select {
				case <-ctx.Done(): // 注意 context 不在同一个 goroutine
					return
				case origin <- n:
					n++
				}
			}
		}()

		return origin
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for val := range gen(ctx) { // 消费者：消费 gen(ctx) 中获取到的 <-chan int
		fmt.Println("val:", val)
		if val == 5 {
			break // trigger to call cancel()
		}
	}
}
~~~

gen 函数会返回一个 Channel，在另外的 goroutine 中产生数据，并将产生的数据发送到 Channel 中。该 Channel 作为数据源，相当于是**生产者**。

原先的 goroutine 中是一个数据的消费者，执行 cancel 表示的是不需要数据，相当于通知生产者退出。此时退出后，相当于是**不会让 goroutine 一直运行，从而导致资源泄漏**。

**关键内容**：gen 函数中不会让 goroutine 泄漏资源

### func WithDeadline

~~~go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)

// A timerCtx carries a timer and a deadline. It embeds a cancelCtx to
// implement Done and Err. It implements cancel by stopping its timer then
// delegating to cancelCtx.cancel.
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}
~~~

关键点：timerCtx 结构体类型**内嵌**了 cancelCtx 类型

WithDeadline 函数会返回返回一个 parent 的拷贝，同时会附带有一个已调整过的 deadline。如果 parent Context 的 deadline 早于参数 d 表示的时间，返回的结果就和 parent 相同。也就是说，d 必须早于 parent 对应的 deadline，才是会创建一个新的 timerCtx。

![](./Snipaste_2021-09-15_12-00-02.png)

有 3 种情况能让 WithDeadline 的 timerCtx 的 Channel 关闭：

1. 截止时间已到
2. cancel 被调用
3. parent 的 Channel 被关闭

和 cancelCtx 类似，都是需要调用 cancel 的，否则会导致资源泄漏。

~~~go
const shortDuration = 1 * time.Second

func withDeadline() {
	d := time.Now().Add(shortDuration)

	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
~~~

在需要具备有倒计时功能（deadline）中使用 timeCtx，起到作用的是结构体中的 timer。

### func WithTimeout

~~~go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
~~~

比较 WithTimeout 和 WithDeadline 的区别，前者的入参类型是 time.Duration（**时间段**），后者是 time.Time（确定的**时间点**）。

因此，本质上，WithTimeout 相当于是：`WithDeadline(parent, time.Now().Add(timeout))` 可以这样理解，在当前时间的基础上增加**时间段**。因此在使用上，其实没什么区别的。

~~~go
func withTimeout() error {
	delay, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // if slowOperation 在 timeout 之前完成，则释放资源
	return slowOperationWithTimeout(delay)
}

const duration = 2 * time.Second

func slowOperationWithTimeout(ctx context.Context) error {
	channel := make(chan int)
	go func() {
		// mock for slow operation
		time.Sleep(duration)
		channel <- 2
	}()

	select {
	case <-ctx.Done():
		fmt.Println("times up!")
		return errors.New(ctx.Err().Error())
	case <-channel:
		fmt.Println("nornal result return")
		return nil
	}
}
~~~

上面是一个**超时控制模型**，也就是说：必须在 delay 之前完成，否则返回超时 error。

### type CancelFunc

CancelFunc 类型的值表示的是一个舍弃功能，即让执行流程结束。

~~~go
// A CancelFunc tells an operation to abandon its work.
// A CancelFunc does not wait for the work to stop.
// A CancelFunc may be called by multiple goroutines simultaneously.
// After the first call, subsequent calls to a CancelFunc do nothing.
type CancelFunc func()
~~~

一个 CancelFunc 可能会并发执行，但这是安全的。首次执行 CancelFun 后，后续再次调用，不会做任何事。

比如：

~~~go
// cancel closes c.done, cancels each of c's children, and, if
// removeFromParent is true, removes c from its parent's children.
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
~~~

并发安全，是因为在其中加了 `c.mu.Lock()`。

### type Context

