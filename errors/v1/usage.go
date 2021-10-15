package v1

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func newUsage() {
	err := errors.New("a new error")
	log.Println(err.Error())

	if unwrap := errors.Unwrap(err); unwrap == nil {
		log.Println("error is not Wrapped error")
	}
}

func wrapErr() {
	err := errors.New("a new error")
	log.Println(err.Error())

	wrapErr := fmt.Errorf("%w, [%s]", err, "<other>") // 如果包含了 %w 占位符，表示对应是 error 实例
	if _, ok := wrapErr.(interface {
		Unwrap() error
	}); ok {
		log.Printf("[%v] is a Wrapped error, wrapErr.s:%s", wrapErr, wrapErr.Error())
	}
}

type errorIs struct {
	msg string
	err error
}

func (e errorIs) Error() string {
	return e.msg
}

func (e errorIs) Is(target error) bool {
	return e.err == io.ErrClosedPipe
}

func (e errorIs) Unwrap() error {
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

type errorChain struct {
	msg string
	err error
}

func (e errorChain) Error() string {
	return e.msg
}

func (e errorChain) Unwrap() error {
	return e.err
}

func (e errorChain) Wrap(msg string) errorChain {
	err := errorChain{
		msg: msg,
	}
	e.err = err
	//FIXME how to create error chain?
	return err
}

//FIXME 如何组建 error Chain？
func createErrorChain() {
	// os.PathError

	_ = errorChain{
		msg: "1",
	}

	// 构造一个 error Chain 状结构
	// http.ErrBodyNotAllowed: a normal error instance

}
