# Timer





~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	time.AfterFunc(time.Second*5, func() {
		fmt.Println("Time is up!")
		ch <- struct{}{}
	})

	go func() {
		tick := time.Tick(time.Second)
		for {
			select {
			case <-tick:
				fmt.Println("Tick!")
			}
		}
	}()
	<-ch
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
Tick!
Tick!
Tick!
Tick!
Time is up!
~~~

这只“猫”只会 Tick! 5 次在第五次时，会立即停止。作为结果，其最后一次并不会输出 Tick!

~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	timer := time.AfterFunc(time.Second*5, func() {
		fmt.Println("Time is up!")
		ch <- struct{}{}
	})

	go func() {
		tick := time.Tick(time.Second)
		index := 0
		for {
			select {
			case <-tick:
				fmt.Println("Tick!")
				index++
				if index == 3 {
					timer.Stop()
				}
			}
		}
	}()
	<-ch
}
~~~

而这个示例程序，因为在超时时间到来之前已经执行了 timer.Stop()，因此永远不会在执行 AfterFunc 中的匿名函数，也就永远不会停止！









~~~go
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	ch := make(chan struct{})

	var callStatus uint32
	timer := time.AfterFunc(time.Second*2, func() {
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
			// it is not request timeout! atomic.CompareAndSwapUint32(&callStatus, 0, 1) success!
			fmt.Println("it is not request timeout!")
			return
		}
		fmt.Println("Time is up!")

		time.Sleep(time.Second * 3)
		sendMessage(ch)
	})

	// simulat call http request
	time.Sleep(time.Second * 3)
	if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
		// request timeout! atomic.CompareAndSwapUint32(&callStatus, 0, 2) success!
		fmt.Println("request timeout!")
		return
	}
	timer.Stop()
	fmt.Println("Stopp latter!")
	sendMessage(ch)

	<-ch
}

func sendMessage(ch chan<- struct{}) {
	go func() {
		ch <- struct{}{}
	}()
}
~~~

`time.Sleep(time.Second * 3)` 用来模拟网络请求的延时，该延时是不确定的，有可能会导致 time.Timer 超时。

问题在于：AfterFunc 中的匿名函数是在是哪个 goroutine 执行的？

作为校验，修改上述代码，删除 return：

~~~go
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	ch := make(chan struct{})

	var callStatus uint32
	timer := time.AfterFunc(time.Second*2, func() {
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
			// it is not request timeout! atomic.CompareAndSwapUint32(&callStatus, 0, 1) success!
			fmt.Println("it is not request timeout!")
			return
		}
		fmt.Println("Time is up!")

		time.Sleep(time.Second * 3)
		sendMessage(ch)
	})

	// simulat call http request
	time.Sleep(time.Second * 3)
	if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
		// request timeout! atomic.CompareAndSwapUint32(&callStatus, 0, 2) success!
		fmt.Println("request timeout!")
		// return
	}
	if !timer.Stop() {
		size := len(timer.C)
		fmt.Println(size)
		if (size) != 0 {
			<-timer.C
		}
	}
	fmt.Println("Stopp latter!")
	sendMessage(ch)

	<-ch
}

func sendMessage(ch chan<- struct{}) {
	go func() {
		ch <- struct{}{}
	}()
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
Time is up!
request timeout!
0
Stopp latter!
~~~



