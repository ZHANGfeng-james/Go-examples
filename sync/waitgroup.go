package sync

import (
	"log"
	"sync"
	"time"
	"unsafe"
)

func waitGroup() {
	var wg sync.WaitGroup

	log.Printf("sizeof:%d", unsafe.Sizeof(wg)) // 12个字节，正好是 state1 [3]uint32 内存占用

	wg.Add(2)

	go func() {
		time.Sleep(3 * time.Second)
		wg.Done()
		log.Println("3 second waiting, over")
	}()

	go func() {
		time.Sleep(5 * time.Second)
		wg.Done()
		log.Println("5 second waiting, over")
	}()

	go func() {
		wg.Wait()
		log.Println("goroutine wg.Wait()...")
	}()

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func waitGroupReuse() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done() // 计数器减1
		time.Sleep(10 * time.Millisecond)
		wg.Add(1) // 计数值加1
	}()
	wg.Wait() // 主goroutine等待，有可能和第7行并发执行
}

func calcWaitGroupCount() {
	// 获取 WaitGroup 的计数值
	var wg sync.WaitGroup

	wg.Add(20)
	// state1 [3]uint32，当前是 64bit 的
	// state1[0]：Waiter数目，也就是调用了 Wait() 的 goroutine 的数量
	// state1[1]：计数值

	for i := 10; i > 0; i-- {
		go func(i int) {
			wg.Wait()
		}(i)
	}

	time.Sleep(1 * time.Second)
	ptr := (*uint64)(unsafe.Pointer((uintptr(unsafe.Pointer(&wg)))))
	counter := int32(*ptr >> 32)
	waiters := uint32(*ptr)
	log.Printf("waiters:%d, counter:%d", waiters, counter)

	wg.Add(-20)
	wg.Wait()
}
