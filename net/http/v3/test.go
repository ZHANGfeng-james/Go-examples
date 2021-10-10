package main

import (
	"fmt"
	"time"
)

func test() {

	// 测试 panic 的 recover 机制，是对 goroutine 有效的
	// 也就是说，如果一个 goroutine 中出现 panic，不会导致系统整体崩溃

	ch := make(chan int)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		fmt.Println("panic goroutine start!")
		// a panic goroutine
		time.Sleep(3 * time.Second)
		panic("a panic")
	}()

	go func() {
		fmt.Println("normal goroutine start!")
		time.Sleep(5 * time.Second)
		ch <- 0
	}()

	<-ch
	fmt.Println("over")
}

func testRecover() {
	defer func() {
		fmt.Println("a panic!")
	}()

	fmt.Println("tesetRecover")
	panic("panic") // 模拟出现 panic
}
