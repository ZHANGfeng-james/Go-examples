package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func withCancel() {
	gen := func(ctx context.Context) <-chan int {
		origin := make(chan int)
		n := 0

		go func() {
			for {
				select {
				case <-ctx.Done(): // ctx 和原先 gen 函数的调用者中 context 不在同一个 goroutine
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

const shortDuration = 1 * time.Second

func withDeadline() {
	d := time.Now().Add(shortDuration)

	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done(): // context deadline exceeded 截止时间到时，channel 关闭
		fmt.Println(ctx.Err())
	}
}

func withDeadlineProducer() {
	gen := func(ctx context.Context) chan int {
		origin := make(chan int)
		n := 0

		go func() {
			for {
				select {
				case <-ctx.Done(): // context deadline exceeded 时，channel 关闭
					close(origin)
					return
				case origin <- n:
					n++
				}
			}
		}()

		return origin
	}

	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	for ele := range gen(ctx) { // 持续阻塞，直到 channel close
		time.Sleep(200 * time.Millisecond)
		fmt.Println(ele)
	}
}

func withTimeout() error {
	delay, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 必须在 delay 之前完成，否则返回超时 error
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

func doneUsage() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	channel := make(chan int)
	go stream(ctx, channel)

	for ele := range channel {
		fmt.Println(ele)
	}
}

func stream(ctx context.Context, out chan<- int) error {
	for {
		time.Sleep(500 * time.Millisecond)
		v, err := doSomething(ctx)
		if err != nil {
			return err // reason 2:
		}
		select {
		case out <- v:
			// nornal case
		case <-ctx.Done():
			close(out)
			return ctx.Err() // reason 1:
		}
	}
}

func doSomething(ctx context.Context) (int, error) {
	return 1, nil
}

func withValueUsage() {
	type favContextKey string
	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))
}

func withValueUseNormalType() {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := "language"
	ctx := context.WithValue(context.Background(), k, "Go")
	f(ctx, k)
	f(ctx, "color")
}

func cancelContextPropagate() {
	parent, pCancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		pCancel()
	}()

	val1 := context.WithValue(parent, "1", "1")
	son, sCancel := context.WithCancel(parent)
	defer sCancel()

	select {
	case <-son.Done():
		log.Println("son cancel() run...")
	}

	log.Println("Value Context:", val1.Value("1"))
}
