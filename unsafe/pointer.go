package unsafe

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/go-examples-with-tests/builtin"
)

func usageUnsafePointer() {
	var ptr unsafe.Pointer
	log.Println(ptr)

	var f float64
	f = 3
	log.Printf("float64bits:%d", float64bits(f))
}

func float64bits(f float64) uint64 {
	var ptr *float64
	ptr = &f
	return *((*uint64)(unsafe.Pointer(ptr)))
}

func getVarAddrUsingPointer(ptr *int64) uintptr {
	add := (uintptr)(unsafe.Pointer(ptr))
	return add
}

type Person struct {
	name string
}

func usageUnsafeSizeof() {
	person := Person{}
	log.Printf("%d", unsafe.Sizeof(person))
	var s string
	log.Printf("%d", unsafe.Sizeof(s))
}

func changeStructField() {
	// builtin.NoEmptyStruct 中的 name 是不可导出的结构体类型
	s := builtin.NewNoEmptyStruct("Katyusha", 18)
	log.Println(s)

	// s --> *NoEmptyStruct + 16 --> age
	age := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(16)))
	*age = 28
	log.Println(s)
}

func changeArrayEle() {
	// 此处必须是数组类型，[4]int 如果是切片则会发生改变
	var slice = [4]int{1, 2, 3, 4}
	log.Printf("origin data:%v", slice)

	changeID := len(slice) - 1
	// *[]int --> unsafe.Pointer --> *int
	ptr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice)) + unsafe.Sizeof(&slice[0])*uintptr(changeID)))
	*ptr = 100

	log.Printf("changed data:%v, changeID:%d", slice, changeID)
}

func getSliceHeaderInfo() {
	slice := []int{1, 2, 3, 4, 5, 6}

	slice = append(slice, 100)

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	log.Printf("slice cap:%d, len:%d", hdr.Cap, hdr.Len)
}

func changeStringContent() {
	var origin = "abc"
	oData := (*reflect.StringHeader)(unsafe.Pointer(&origin))

	var s string = "123"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // case 1
	hdr.Data = oData.Data                              // case 6 (this case)
	hdr.Len = len(origin)

	log.Printf("%s", s)
}

func byteSliceToStringNoCopy() {
	bytes := []byte{0x41, 0x42}
	log.Printf("origin data:%s", bytes)
	bytesHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))

	var s string
	log.Printf("dst s:%s", s)

	// []byte --> string 实现 zero-copy
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	log.Printf("string length:%d", hdr.Len)

	hdr.Data = bytesHdr.Data
	hdr.Len = bytesHdr.Len

	bytes[0] = 0x42
	log.Printf("changed:%s", s)
}
