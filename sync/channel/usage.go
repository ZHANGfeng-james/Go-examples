package channel

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func unBufferChannelUsage() {
	ch := make(chan struct{})
	go func() {
		log.Println("Go Go Go...")
		<-ch
	}()

	ch <- struct{}{} // write or read get the same result
}

func unBufferChannelOther() {
	ch := make(chan struct{})
	go func() {
		log.Println("Go Go Go...")
		ch <- struct{}{}
	}()

	<-ch // write or read get the same result
}

func bufferChannelUsage() {
	ch := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		log.Println("Go Go Go...")
		ch <- struct{}{}
	}(ch)
	<-ch
}

func bufferChannelOther() {
	ch := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		log.Println("Go Go Go...")
		<-ch
	}(ch)
	ch <- struct{}{}
}

func bufferChannelSameGoroutine() {
	ch := make(chan struct{}, 1)

	log.Println("same goroutine, write:")
	ch <- struct{}{}

	log.Println("same goroutine, read:")
	<-ch

	log.Println("over")
}

func selectChannelUsage() {
	ch := make(chan struct{})

	// select 的语法特性，选择未被阻塞的路径执行下去
	select {
	case <-ch:
		log.Println("read ele from ch")
	default:
		log.Println("contine, no block")
	}
}
