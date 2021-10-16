os/signal 标准库用于访问到来的**信号（操作系统在检测系统状态、输入后发出的系统信号）**。这些系统信号大部分被使用在**类 Unix 操作系统**上，当然也会有一些在 Windows 和 Plan9 上。



# 1 介绍







# 2 API 的使用



~~~go
func Notify(c chan<- os.Signal, sig ...os.Signal)
~~~

Notify 函数会让进程接下来**接收到的 os.Signal 会写入到 c Channel 中**。如果入参中没有给定任何 os.Signal，那么**所有的系统信号**都会 Write 到 c Channel 中；反之，**只有指定的 os.Signal 会**。

`Package signal` 发送到 Channel 中时，**不会引起阻塞**，但是调用者需要确保 Channel 有足够的 Buffer 空间容纳可能会到来的 os.Signal。**对于接受单一 os.Signal 来说，容量为 1 的缓冲区是足够的**。

~~~go
func process(sig os.Signal) {
	n := signum(sig)
	if n < 0 {
		return
	}

	handlers.Lock()
	defer handlers.Unlock()

	for c, h := range handlers.m {
		if h.want(n) {
			// send but do not block for it
			select {
			case c <- sig:
			default: // 即便 c <- sig 阻塞了，也因为 default 不会导致阻塞！
			}
		}
	}

	// Avoid the race mentioned in Stop.
	for _, d := range handlers.stopping {
		if d.h.want(n) {
			select {
			case d.c <- sig:
			default:
			}
		}
	}
}
~~~

标准库允许在同一个 Channel 上多次调用 Notify，每次调用相当于扩大了发送给 Channel 的 os.Signal 集合。当然，调用 Stop 能够清空。标准库允许在不同的 Channel 和相同的 os.Signal 上多次调用 Notify，每个 Channel 都会收到这些 os.Signal 的副本。

比如下面尝试写这样的需求：

~~~go
package v1

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func signalUsage() {
	log.Println("Process running...")

    // Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt) // wait for get os.Interrupt signal

	log.Println("goroutine start to sleep")
	time.Sleep(5 * time.Second)
	log.Println("goroutine sleep over, and weak up...")

	s := <-ch
	log.Println(s)
}
~~~

上述程序中的 `time.Sleep(5 * time.Second)` 就是模拟在启动监听后，还有很复杂的工作要做。此时进入 goroutine Sleep 状态后，**连续多次按下 Ctrl + C 都不能让程序停止运行**，直到退出 goroutine Sleep 状态后，再次按下才会收到 os.Signal。

更换程序：使用 Buffered Channel，则不会出现这种情况。在退出 goroutine Sleep 状态后，程序自动停止运行，相当于是此时 ch 已经接受到了 os.Signal。

解答另一个疑惑：signal_unix.go 中有如下程序，应该是有 goroutine 在持续侦测 os.Signal，并不会因为某个 goroutine 进入到了 time.Sleep 状态而停止

~~~go
func loop() {
	for {
		process(syscall.Signal(signal_recv()))
	}
}
~~~







