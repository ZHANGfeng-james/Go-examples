package builtin

import (
	"log"
	"sync"
	"unsafe"
)

func printMap(input map[string]int) {
	// map 是一个集合（key-value的容器），如果是要查询 key 是否在集合中，是不知道的
	key := "3"
	if value, ok := input[key]; ok {
		log.Printf("get key-value:[%s]=%d", key, value)
	}
	// 遍历 map 集合
	for key, value := range input {
		log.Printf("for range map: map[%s]=%d", key, value)
	}
}

func createMap() {
	kVal := make(map[string]int)

	kVal["1"] = 1
	kVal["2"] = 2
	printMap(kVal)

	val := map[string]int{}
	val["1"] = 10
	val["2"] = 20
	printMap(val)
}

func createMapNil() {
	var kVal map[string]int // kVal 是 nil 值，不能向其中 assign entry
	log.Printf("kVal Type:%T", kVal)

	if kVal == nil {
		log.Println("kVal is nil!")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	kVal["1"] = 1
}

func getMapSizeof() {
	var kVal map[string]int
	log.Printf("sizeof(map):%d, type: %T", unsafe.Sizeof(kVal), kVal)
}

func retainEle() {
	kVal := make(map[string]int)

	kVal["1"] = 1
	kVal["2"] = 2

	// retain 查询 map 是否存在某个值
	val := kVal["1"]
	log.Printf("value is %d", val)

	val = kVal["3"]
	log.Printf("key:3 --> val:%d", val)
}

func callFuncUseMap() {
	kVal := make(map[string]int)
	kVal["1"] = 1
	kVal["2"] = 2
	kVal["3"] = 3

	called(kVal)

	if _, ok := kVal["1"]; ok {
		log.Println("called delete failed!")
	} else {
		log.Println("called delete success!")
	}
}

func called(kVal map[string]int) {
	delete(kVal, "1")
}

func forRangeMap() {
	kVal := make(map[string]int)
	kVal["1"] = 1
	kVal["2"] = 2
	kVal["3"] = 3

	for range []int{1, 2, 3} {
		for k, v := range kVal {
			log.Printf("key:%s, value:%d", k, v)
		}
		log.Println()
	}
}

func mapConcurrent() {
	kVal := make(map[string]int)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kVal["1"] = 10
		}()
	}

	wg.Wait()
	log.Printf("key:\"1\" --> valu:%d", kVal["1"])
}

func mapGoAction() {
	counter := struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Lock()
			counter.m["1"] = 100
			counter.Unlock()
		}()
	}

	wg.Wait()
	log.Printf("key:\"1\" --> valu:%d", counter.m["1"])
}

type FooMap struct {
	sync.Mutex
	kVal map[int]int
}

var kv FooMap = FooMap{
	kVal: make(map[int]int),
}

func storeUseBuiltinMap(key, value int) {
	kv.Lock()
	kv.kVal[key] = value
	kv.Unlock()
}

func mapKey() {
	kVal := make(map[float32]int)

	var key float32
	key = 1.2
	value := 1

	kVal[key] = value

	log.Printf("key:%f, value:%d", key, kVal[key])
}

func forRangeAndDelete() {
	// 边遍历边删除
	kVal := make(map[string]int)

	kVal["1"] = 1
	kVal["2"] = 2
	kVal["3"] = 3

	for k, v := range kVal {
		log.Printf("key:%s, value:%d", k, v)
		if k == "1" {
			delete(kVal, k)
		}
	}

	if _, ok := kVal["1"]; ok {
		log.Println("for range and delete failed!")
	}
	log.Printf("size of map:%d", len(kVal))

	for k, v := range kVal {
		log.Printf("key:%s, value:%d", k, v)
	}
}

func forRangeAndAdd() {
	// 边遍历边添加
	kVal := make(map[string]int)

	kVal["1"] = 1
	kVal["2"] = 2
	kVal["3"] = 3

	for k, v := range kVal {
		log.Printf("key:%s, value:%d", k, v)
		if k == "1" {
			kVal["10"] = 10
		}
	}
	log.Printf("size of map:%d", len(kVal))

	for k, v := range kVal {
		log.Printf("key:%s, value:%d", k, v)
	}
}
