package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	set := make(map[int]bool)

	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		go func() {
			lock.Lock()
			defer lock.Unlock()
			if _, ok := set[100]; !ok {
				fmt.Println("100")
				set[100] = true
			}
		}()
	}

	time.Sleep(1 * time.Second)
}
