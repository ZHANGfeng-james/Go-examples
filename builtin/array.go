package builtin

import (
	"log"
	"unsafe"
)

func printArray(arr [3]int) {
	for i, v := range arr { // for...range...遍历数组集合
		log.Printf("Array[%d]=%d", i, v)
	}
}

func arrayCreate() {
	var val [3]int // 每个元素，赋值为 int 类型令旨
	log.Printf("len:%d, cap:%d", len(val), cap(val))
	printArray(val)
}

func arrayCallFunc() {
	// [3]int 仅仅只是一种类型，和 []int 一样，若用于定义变量，则必须初始化
	val := [3]int{1, 2, 3}
	log.Printf("val[%d]=%d", len(val), val[len(val)-1])
	// [3]int{} 定义了一个变量，该变量的类型是 [3]int，同时赋予初始值
	callfuncWithArray(val)
	log.Printf("val[%d]=%d", len(val), val[len(val)-1])
}

// callfunc call function with a array, and return array's length
func callfuncWithArray(arr [3]int) int {
	arr[len(arr)-1] = 100
	printArray(arr)
	return len(arr)
}

func getArrayDst() {
	val := [10]int{1, 2, 3}
	log.Printf("array original address:%p", &val)
	for i, v := range val {
		log.Printf("val[%d]=%d, addr:%p, %p", i, v, &v, &val[i])
	}
	printAddr(val)

	log.Printf("unsafe.Sizeof(val)=%d", unsafe.Sizeof(val))
}

func printAddr(value [10]int) {
	for index := 0; index < len(value); index++ {
		log.Printf("index= %d; addr= %p.\n", index, &value[index])
	}
}
