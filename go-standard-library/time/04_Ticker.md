# Ticker









看看 time.Ticker 在出现延迟获取 chan 中 Time 时的结果：

~~~go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	tick := time.NewTicker(time.Second)
	ch := tick.C

	index := 0
	for {
		select {
            case cur := <-ch:
			fmt.Println(cur)

			index++
			if index == 10 {
				return
			}
            // 随机延迟一段时间
			num := rand.Intn(10)
			fmt.Println(num)
			time.Sleep(time.Duration(num) * time.Second)
		}
	}
}
~~~

第 1 次获取 chan 值：

~~~go
2021-05-13 13:45:49.851535 +0800 CST m=+1.002252701
1
~~~

执行 `time.NewTicker(time.Second)` 时，计时已开始，首次获取时间，并**休眠 1s**。

第 2 次获取 chan 值：

~~~go
2021-05-13 13:45:50.8519443 +0800 CST m=+2.002662001
7
~~~

实际频率和 time.Ticker 一致，本次获取时间是正常的。此时，goroutine 需要**休眠 7s**。goroutine 休眠，但是 time.Ticker 仍然能够在 1s 后将 Now() 值写到 chan 中，实际写入的是：`2021-05-13 13:45:51.8518651 +0800 CST m=+3.002582801`。在写入该值后，chan 只能容纳一个元素，之后试图写入 Now() 将会直接阻塞（不再接收 Now() 值），直到下次 <-chan！

第 3 次获取 chan 值：

~~~go
2021-05-13 13:45:51.8518651 +0800 CST m=+3.002582801
7
~~~

此时读取的值已在第 2 次获取 chan 值时，goroutine 延迟 7s 的过程中已经分析。在第 2 次获取 chan 值后需要休眠 7s 后获取了 chan 值，在延迟到 7s 结束后，向 ch<-57 阻塞，执行 default。紧接着 51<-ch ，输出该值：`2021-05-13 13:45:51.8518651 +0800 CST m=+3.002582801`，继续休眠 7s。

第 4 次获取 chan 值：

~~~go
2021-05-13 13:45:58.8519567 +0800 CST m=+10.002674401
9
~~~

ch<-64 阻塞，执行 default；58<-ch

第 5 次获取 chan 值：

~~~go
2021-05-13 13:46:05.8532722 +0800 CST m=+17.003989901
1
~~~

ch<-73 阻塞，执行 default；65<-ch

第 6 次获取 chan 值：

~~~go
2021-05-13 13:46:14.8527267 +0800 CST m=+26.003444401
8
~~~

从上述结果列表来看，关键的一点是：time.Ticker 的 Tick 声到时，是先从 chan 中取值？还是先向 chan 写值？

1. 如果是先从 chan 中取值，则值取出后 chan 为空，可向 chan 中写入值；
2. 如果是先向 chan 写值，因 chan 容量已满而会直接阻塞，time.Ticker 会丢弃当前值；紧接着从 chan 中读取值，并返回。

从结果上来看，time.Ticker 使用的是第 2 种方式，问题是如何去验证？

~~~go
func sendTime(c interface{}, seq uintptr) {
	// Non-blocking send of time on c.
	// Used in NewTimer, it cannot block anyway (buffer).
	// Used in NewTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan Time) <- Now():
	default:
	}
}
~~~

如果从 channel 的角度来分析，当 runtime 中的 time.go 定时时间到时，是先向 chan 中写入 Now()？还是先从 chan 中读取？

因此，问题转移到：如果 goroutine 同时发生了读写操作，究竟顺序是怎样的？为什么在 time.Ticker 中，向 chan 写值为先？