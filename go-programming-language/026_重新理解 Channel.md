

~~~go
package main

import (
	"fmt"
)

func main() {
    // 无缓存的 Channel
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
}
~~~

上述代码会直接阻塞，导致 `fatal error: all goroutines are asleep - deadlock! `，而下述代码却不会：

~~~go
package main

import (
	"fmt"
)

func main() {
    // 带有一个缓存位置的 Channel
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
}
~~~

上述两者唯一的区别是：`ch := make(chan int)` 和 `ch := make(chan int)`，即 Channel 是否有缓存！