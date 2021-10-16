

# 1 Channel 的使用场景

Channel 的使用时机、特征，也就是：

1. 在何种应用场景中，可以使用 Channel 实现需求，这些场景有什么特征、共同之处？
2. Buffered Channel 和 Unbuffered Channel 有什么区别？怎么使用？

## 1.1 Unbuffered Channel VS Buffered Channel

**｜Unbuffered Channel**

~~~go
func unBufferChannelUsage() {
	ch := make(chan struct{})
	go func() {
		log.Println("Go Go Go...")
		<-ch
	}()

	ch <- struct{}{} // write or read get the same result
}
~~~

通过 Channel，我们解决了**同步问题**。此次需求就是要等到输出了 `Go Go Go`，之后程序才能结束，和**等待**是一个意思。也就是说 main goroutine 需要等待打印 goroutine 运行结束后才能继续执行。另外对于 Unbuffered Channel 来说，下面的程序能**达到相同的效果**：

~~~go
func unBufferChannelOther() {
	ch := make(chan struct{})
	go func() {
		log.Println("Go Go Go...")
		ch <- struct{}{}
	}()

	<-ch // write or read get the same result
}
~~~

总结出 Unbuffered Channel 的**特征**：不管是在 goroutine（需要是不同的 goroutine）内读或者是写，都只能是在一个读一个写完成后，两者才能继续进行下去。这就是使用 Unbuffered Channel 能够达到的**同步效果**。

**｜Buffered Channel**

同样的程序，但是**换成了有缓冲区**的 Channel：

~~~go
func bufferChannelUsage() {
	ch := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		log.Println("Go Go Go...")
		ch <- struct{}{}
	}(ch)
    <-ch // 初始状态，Buffer区没有 stuct{}{}，阻塞直到 goroutine 写入 struct{}{}
}
~~~

与之形成对比的是:

~~~go
func bufferChannelOther() {
	ch := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		log.Println("Go Go Go...")
		<-ch
	}(ch)
    ch <- struct{}{} // 带有缓冲区（可容纳1个struct{}{}实例，必会导致阻塞）
}
~~~

只要 Buffered Channel 的缓冲区足够大，写进去的实例不需要被读出来，goroutine 也会继续执行。

因此，这一节内容，只要记住这个结论：

Unbuffer Channel 只有等到读和写都结束后，不同的 goroutine 才会继续执行，否则一方总是会被阻塞！这就是能够达成**同步效果**的原因。

## 1.2 