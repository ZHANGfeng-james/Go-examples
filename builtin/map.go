package builtin

import (
	"log"
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
	log.Printf("sizeof(map):%d", unsafe.Sizeof(kVal))
}
