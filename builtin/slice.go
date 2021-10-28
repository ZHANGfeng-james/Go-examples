package builtin

import (
	"log"
	"unsafe"

	"github.com/go-examples-with-tests/tools"
)

func printSlice(slice []int) {
	log.Printf("cap(slice):%d, len(slice):%d", cap(slice), len(slice))
	for i, v := range slice { // for...range...遍历数组集合
		log.Printf("Array[%d]=%d", i, v)
	}
}

func createSliceUseMake(len int) {
	slice := make([]int, len, 2*len) // 使用 make 创建 slice，而不能用于 Array 中
	printSlice(slice)
}

func getSliceInfo() {
	origin := make([]int, 0)
	origin = append(origin, 1, 2, 3, 4, 5, 6)
	log.Print(tools.SliceInfo("origin", origin))

	dst := make([]int, len(origin), len(origin)) // len(dst) should not be zero!
	log.Print(tools.SliceInfo("dst", dst))

	num := copy(dst, origin) // copy() return the minmum of len(dst) and len(src)
	log.Printf("copy num:%d", num)
	log.Print(tools.SliceInfo("dst", dst))
}

// builtin function:[func copy(dst, src []Type) int]
func copySliceUseBuiltin(src []int) (int, []int) {
	result := make([]int, len(src), len(src))
	num := copy(result, src)
	return num, result
}

func sliceCallFunc() {
	// [3]int 仅仅只是一种类型，和 []int 一样，若用于定义变量，则必须初始化
	val := []int{1, 2, 3}
	log.Printf("val[%d]=%d", len(val), val[len(val)-1])
	// [3]int{} 定义了一个变量，该变量的类型是 [3]int，同时赋予初始值
	callfuncWithSlice(val)
	log.Printf("val[%d]=%d", len(val), val[len(val)-1])
}

// callfunc call function with a array, and return array's length
func callfuncWithSlice(arr []int) int {
	arr[len(arr)-1] = 100
	printSlice(arr)
	return len(arr)
}

func getSliceAddr() {
	var tmp []int64
	log.Printf("tmp:%p, variable tmp's address:%p", tmp, &tmp)

	var slice = []int64{1, 2, 3, 4}
	log.Printf("slice:%p, variable slice's address:%p", slice, &slice)
	for i, v := range slice {
		log.Printf("slice[%d]=%d, addr:%p", i, v, &slice[i])
	}

	tmp = slice
	log.Printf("tmp:%p, variable tmp's address:%p", tmp, &tmp)
	slice = nil
	log.Printf("tmp:%p, variable tmp's address:%p", tmp, &tmp)

	log.Printf("slice:%p, variable slice's address:%p", slice, &slice)

	var values []int64
	log.Printf("values:%p, variable values's address:%p", values, &values)
	values = make([]int64, 0)
	log.Printf("values:%p, variable values's address:%p", values, &values)
}

func getSliceSizeof() {
	var slice []int
	log.Printf("sizeof(slice):%d", unsafe.Sizeof(slice))
}
