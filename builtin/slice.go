package builtin

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"unsafe"

	"github.com/go-examples-with-tests/tools"
)

func printSlice(slice []int) {
	log.Printf("cap(slice):%d, len(slice):%d", cap(slice), len(slice))
	for i, v := range slice { // for...range...遍历数组集合
		log.Printf("Array[%d]=%d", i, v)
	}
}

func readEleFromSlice() {
	slice := []int{1, 2, 3}
	ele := slice[0]
	log.Printf("%d", ele)

	nums := slice[0:2]
	log.Printf("%d", nums)

	arr := make([]int, 0, 3)

	arr = append(arr, 1)
	arr = append(arr, 2)
	arr = append(arr, 3)
	fmt.Println(arr) // [1 2 3]

	slice = make([]int, 1, 3)
	fmt.Printf("%d, %d\n", len(slice), cap(slice)) // 1, 3

	slice = append(slice, 1)
	fmt.Println(slice)                             // [0 1]
	fmt.Printf("%d, %d\n", len(slice), cap(slice)) // 2, 3
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

func makeSliceUsingArray() {
	arr := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	slice := arr[2:4:7]
	// 2, 5, slice:[3 4]
	fmt.Printf("%d, %d, slice:%v\n", len(slice), cap(slice), slice)

	slice = append(slice, 1) // 此时 len 值为 2，也就是 index 处理下标是 2
	slice = append(slice, 1)
	slice = append(slice, 1)
	slice = append(slice, 1)
	/*
		1, 2, 3, 4, 5, 6, 7, 8, 9
		      ^  ^        ^
	*/
	// 6, 10, slice:[3 4 1 1 1 1]
	fmt.Printf("%d, %d, slice:%v\n", len(slice), cap(slice), slice)
	fmt.Println(arr) // [1 2 3 4 1 1 1 8 9]

	fmt.Printf("%p", &slice[0])
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

	log.Printf("sizeof(slice):%d, sizeof(tmp):%d", unsafe.Sizeof(slice), unsafe.Sizeof(tmp))

	var values []int64
	log.Printf("values:%p, variable values's address:%p", values, &values)
	values = make([]int64, 0)
	log.Printf("values:%p, variable values's address:%p", values, &values)
}

func getSliceSizeof() {
	var slice []int
	// 此处得到的是 slice 这个 slice 变量内存占有情况
	log.Printf("sizeof(slice):%d", unsafe.Sizeof(slice))
}

func sliceAndArrayAddr() {
	// 比较 slice 和 array 的地址
	slice := []int{1, 2, 3}
	array := [3]int{20, 3, 40}
	log.Printf("variable addr(defalut order: slice, array):%p, %p", &slice, &array)

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	log.Printf("slice element addr:%#x", hdr.Data)

	// 证明 &array 获取到的就是数组地址
	// *[3]int --> unsafe.Pointer --> uintptr
	addr := (*int)(unsafe.Pointer((uintptr)(unsafe.Pointer(&array)) + uintptr(len(array)-1)*unsafe.Sizeof(&array[0])))
	log.Printf("array first element:%d", *addr)
	log.Printf("%#x,%#x", &array[len(array)-1], addr)
}

func sliceAppend() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i) // cap: 1-->2-->4, len: 1-->2-->3
		log.Printf("cap:%d, len:%d", cap(s), len(s))
	}

	// modifySlice(s)
	modifSliceMore(s)
	log.Println(s)
}

func modifySlice(s []int) {
	s = append(s, 2048)
	s[0] = 1024
}

func modifSliceMore(s []int) {
	s = append(s, 2048)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	log.Printf("%#x", hdr.Data)

	s = append(s, 4096)
	log.Printf("%#x", hdr.Data)

	log.Printf("cap:%d, len:%d", cap(s), len(s))

	s[0] = 1024
	log.Printf("%v", s) // [1024 1 2 2048 4096]
}

type node struct {
	insert   int
	capacity int
}

func sliceGrow() {
	var tmp = make([]int, 5, 7)
	add := make([]int, 3, 3)
	tmp = append(tmp, add...)
	log.Printf("len:%d, cap:%d", len(tmp), cap(tmp))

	list := make([]node, 0)

	slice := []int{}
	capacity := cap(slice)
	for i := 1; i < 2500; i++ {
		slice = append(slice, i)
		if capacity != cap(slice) {
			list = append(list, node{
				insert:   i,
				capacity: cap(slice),
			})
			capacity = cap(slice)
		}
	}

	for _, v := range list {
		log.Printf("insert:%d, cap:%d", v.insert, v.capacity)
	}
}

func sliceGrowTest() {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	log.Printf("s.cap=%d, s.len=%d", cap(s), len(s))

	x := append(s, 11) // slice grow
	log.Printf("x.cap=%d, x.len=%d", cap(x), len(x))
	y := append(s, 12)

	log.Printf("s:%v, x:%v, y:%v", s, x, y)
}

func sliceAgain() {
	slice := make([]int, 10) // len:10  cap:10
	for i := 0; i < 10; i++ {
		slice[i] = i
	}
	log.Println(tools.SliceInfo("original", slice))

	before := slice[:5]
	log.Println(tools.SliceInfo("before", before))

	after := slice[7:]
	log.Println(tools.SliceInfo("after", after))

	// 删除 5 ~ 6 位置索引的元素
	slice = append(slice[:5], slice[7:]...)
	log.Println(tools.SliceInfo("original", slice))
}

func nilSlice() {
	var slice []int
	log.Println(tools.SliceInfo("origin", slice))
	if slice == nil {
		log.Println("slice is nil")
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
		log.Printf("element data addr:%#x", hdr.Data)
	}
	slice = make([]int, 0)
	if slice == nil {
		log.Println("slice is nil")
	} else {
		log.Println("slice is not nil")
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
		log.Printf("element data addr:%#x", hdr.Data)
		slice = append(slice, 1)
		log.Printf("element data addr:%#x", hdr.Data)
	}
	log.Println(tools.SliceInfo("origin", slice))

	tmp := make([]int, 0)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&tmp))
	log.Printf("element data addr:%#x", hdr.Data)
}

func sliceConcurrent() {
	var slice []int
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			slice = append(slice, i)
		}(i)
	}

	wg.Wait()
	log.Printf("slize len:%d", len(slice))
}

func sliceConcurrentMutex() {
	var slice []int
	var wg sync.WaitGroup

	var mutex sync.Mutex

	for i := 0; i < 10000; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			mutex.Lock()
			slice = append(slice, i)
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	log.Printf("With sync.Mutex, slize len:%d", len(slice))
}
