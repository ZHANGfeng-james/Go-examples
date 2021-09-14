package godefault

import (
	"fmt"
	"testing"
	"time"
)

type StrTo string

func TestPanic(t *testing.T) {

	fmt.Println("logger")

	next()

	fmt.Println("main")

}

func next() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error is not nil")
			fmt.Printf("异常抛出，发生时间: %d\n", time.Now().Unix())
			fmt.Printf("错误信息: %v\n", err)
		}
	}()

	handlerFunc()
}

func handlerFunc() {
	fmt.Println("---")
	panic("模拟panic")
}
