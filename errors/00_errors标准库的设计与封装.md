Go 标准库 errors 封装了**对程序中 error 类型实例的处理和使用**，error 在 Go 中是一个**接口类型**：

~~~go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
~~~

既然是**接口类型**，如果这个 error 实例的值是 nil，表示当前没有 Error 发生。也就是说，如果某个结构体实现了 error 接口，那么在这个类型上创建的实例就是一个 error。比如标准库 errors 中：

~~~go
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
~~~

创建了 errorString 结构体类型，实现了 `Error() string` 方法（注意此处时在 `*errorString` 类型上实现的）。

**New 函数**原型：`func New(text string) error`，用于创建一个 error 实例，其内容仅包含了 Error 的文本信息。

Unwrap、Is、As 函数应用在**可以 Wrap（动词：包装、包裹） 其他 error 实例的场景**中。如果一个 error 实例**具有 `Unwrap() error` 方法（实际上是实现了接口）**，那么这个 error 实例就 Wrap 了其他的 error 实例。如果 `e.Unwrap()` 返回了一个 no-nil 的 error 实例 w，那么我们就说 e Wrap w。从 Unwrap 函数的实现中找到理解的线索：

~~~go
// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface { // 此处就是必须实现的接口
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
~~~

**Unwrap 函数**原型：`func Unwrap(err error) error`，会对一个 Wrapped 的 error 实例**拆包**，如果函数的参数 err 具有 Unwrap 方法（实现了上述所说的接口），那么会在其中**调用一次该方法**。否则，直接返回 nil 值。**一个创建 Wrapped error 的简单方法**是调用 `fmt.Errorf`，使用示例：

~~~go
package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	err := errors.New("this is error")
	log.Println(err.Error())

	wrapErr := fmt.Errorf("%w", err) // 如果包含了 %w 占位符，表示对应是 error 实例
	if _, ok := wrapErr.(interface {
		Unwrap() error
	}); ok {
		log.Printf("[%v] is a Wrapped error", wrapErr)
	}
}
~~~

`fmt.Errorf` 返回的 error 实例就是符合要求的 Wrapped error 实例。从 `fmt.Errorf` 的实现中可以找到线索：

~~~go
...
func Errorf(format string, a ...interface{}) error {
	p := newPrinter()
	p.wrapErrs = true
	p.doPrintf(format, a)
	s := string(p.buf)
	var err error
	if p.wrappedErr == nil {
		err = errors.New(s)
	} else {
		err = &wrapError{s, p.wrappedErr} // s --> msg, p.wrappedErr --> err
	}
	p.free()
	return err
}

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string { // 返回的必须是 error，也就要求实现 error 接口
	return e.msg
}

func (e *wrapError) Unwrap() error { // *wrapError 类型有 Unwrap() error 方法
	return e.err
}
~~~

从 `fmt.Errorf` 函数的实现上来看，相当于是构成了**一条 error 的 Chain**，这个链条上的**各个节点就是 error 实例**。大白话来讲：就是实现了**对 error 实例的一层包装（封装），在其外面包裹了一层**。

**Is 函数**的原型是：`func Is(err, target error) bool`，该函数会在 err 的 Error Chains 中依次调用 Unwrap （获得一个 Unwrap 之后的 error 实例）以此判断是否和 target 匹配。更加优雅的写法是：

~~~go
if errors.Is(err, fs.ErrExist)
~~~

而不是：

~~~go
if err == fs.ErrExist
~~~

因为前者如果包装了 `fs.ErrExist` 实例，依然会返回 true，而后者不会。使用示例：

~~~go
type errorIs struct {
	msg string
	err error
}

func (e errorIs) Error() string {
	return e.msg
}

func (e errorIs) Is(target error) bool { // 实现指定接口的方法
	return e.err == io.ErrClosedPipe
}

func (e errorIs) Unwrap() error { // 实现指定接口的方法
	return io.ErrClosedPipe
}

func isUsage() {
	wrapedErr := fmt.Errorf("%w", io.ErrClosedPipe)
	if is := errors.Is(wrapedErr, io.ErrClosedPipe); is {
		log.Println("wrappedErr wrap [io.ErrClosedPipe]")
	}

	// errors.Is 还有更加丰富的特征
	msg := "a errorIs instance"
	myErr := errorIs{
		msg: msg,
		err: errors.New(msg),
	}
	if is := errors.Is(myErr, io.ErrClosedPipe); is {
		log.Println("errorIs is a io.ErrClosedPipe")
	}
}
~~~

Is 函数本身在调用时不仅会将 err 和 target 比较，还会**在调用了 UnWrap 之后继续比较**。具体实现：

~~~go
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflectlite.TypeOf(target).Comparable()
	for {
		if isComparable && err == target { // 第一次比较
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) { // 第二次比较
			return true
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}
~~~

**As 函数**的原型：`func As(err error, target interface{}) bool`，该函数会依次调用 Unwrap(err)，以此从其中找到能够赋值给 target 的 error 实例，此处的 target 必须是一个指针类型实例。如果执行成功，target 会得到赋予的值，函数返回 true，反之返回 false。最佳实践：

~~~go
var perr *fs.PathError
if erros.As(err, &perr) {
    fmt.Println(perr.Path)
}
~~~

而不是：

~~~go
if perr, ok := err.(*fs.PathError); ok {
    fmt.Println(perr.Path)
}
~~~

因为前者在 err 包装了 `*fs.PathError` 时，依然会返回 true。使用实例如下：

~~~go
func asUsage() {
	wrappedErr := fmt.Errorf("%w", http.ErrNotSupported) // *ProtocolError 类型实例
	if is := errors.Is(wrappedErr, http.ErrNotSupported); is {
		log.Println("wrappedErr wrap [http.ErrNotSupported]")
	}

	var err *http.ProtocolError // wrappedErr中包装的 error 是 *ProtocolError 类型实例
	if result := errors.As(wrappedErr, &err); result {
		log.Printf("%v", err) // err 就是 wrappedErr 中包装的 error
	}
}
~~~

### 总结

特别需要注意的是：

UnWrap 函数中包含一个接口：

~~~go
interface{
    Unwrap() error
}
~~~

返回 error 中包装的 error 实例。

Is 函数中也包含一个接口：

~~~go
interface {
    Is(error) bool
}
~~~

其中入参一般就是 target，用于判断 error 和 target 的关系。

### 疑惑

由 UnWrap 函数，引出疑惑：**如何组建一个 error 的 Chain**？

