package builtin

import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

func channelUsage() {
	var ch chan int
	ch = make(chan int, 10)
	log.Printf("ch sizeof:%d", unsafe.Sizeof(ch))
}

func getEleFromClosedChannel() {
	ch := make(chan int, 10)

	go func() {
		ch <- 100
		close(ch)
	}()

	time.Sleep(1 * time.Second)

	log.Printf("cap:%d, len:%d", cap(ch), len(ch))

	select {
	case ele, ok := <-ch:
		if ok {
			log.Printf("channel is not closed! value:%d", ele)
		} else {
			log.Printf("channel is closed! value:%d", ele)
		}
	}
}

func chanForRange() {
	var ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		select {
		case ch <- i:
			log.Println("input:", i)
		case v := <-ch:
			log.Println("output:", v)
		}
	}
	log.Printf("cap:%d, len:%d", cap(ch), len(ch))
}

func bufferChanRead() {
	ch := make(chan int, 3)

	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	for i := 0; i < 10; i++ {
		v := <-ch
		log.Printf("read index:%d, value:%d", i, v)
	}
}

func channelType() {
	// chan T 是否可以给 <- chan T 和 chan<- T 类型的变量赋值？反过来呢？

	_ = make(chan<- chan int, 1)

}

func channelTaskLoop() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})
	ch4 := make(chan struct{})
	go func() {
		goroutineID := "goroutine 1#"
		for {
			select {
			case <-ch1:
				log.Println(goroutineID)
				time.Sleep(1 * time.Second)
				ch2 <- struct{}{}
			}
		}
	}()

	go func() {
		goroutineID := "goroutine 2#"
		for {
			select {
			case <-ch2:
				log.Println(goroutineID)
				time.Sleep(1 * time.Second)
				ch3 <- struct{}{}
			}
		}
	}()

	go func() {
		goroutineID := "goroutine 3#"
		for {
			select {
			case <-ch3:
				log.Println(goroutineID)
				time.Sleep(1 * time.Second)
				ch4 <- struct{}{}
			}
		}
	}()

	go func() {
		goroutineID := "goroutine 4#"
		for {
			select {
			case <-ch4:
				log.Println(goroutineID)
				time.Sleep(1 * time.Second)
				ch1 <- struct{}{}
			}
		}
	}()

	ch1 <- struct{}{}

	select {}
}

func channelTaskLoop2() {
	ch := make(chan struct{})
	for i := 1; i <= 4; i++ {
		go func(index int) {
			time.Sleep(time.Duration(index*10) * time.Millisecond)
			for {
				<-ch
				fmt.Printf("I am No %d Goroutine\n", index)
				time.Sleep(time.Second)
				ch <- struct{}{}
			}
		}(i)
	}
	ch <- struct{}{}
	time.Sleep(time.Minute)
}
