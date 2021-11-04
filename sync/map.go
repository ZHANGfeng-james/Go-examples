package sync

import (
	"log"
	"sync"
)

// sync.Map 的基本使用
func syncMap() {
	var kVal sync.Map

	values := [4]string{
		"我不行",
		"我一定能行",
		"我必须能行",
		"我确定我能行",
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			kVal.LoadOrStore(i, values[i])
		}(i)
	}

	wg.Wait()

	kVal.Store("string", "int")

	if v, ok := kVal.LoadAndDelete(0); ok {
		log.Printf("delete success! original value:%s", v)
	}

	kVal.Range(func(key, value interface{}) bool {
		log.Printf("key:%v, value:%v", key, value)
		return true
	})
}

var kVal sync.Map

func storeUseSyncMap(key, value int) {
	kVal.Store(key, value)
}
