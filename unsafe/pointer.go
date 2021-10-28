package unsafe

import (
	"log"
	"reflect"
	"unsafe"
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

func changeSliceLength() {
	var origin = "abc"
	oData := (*reflect.StringHeader)(unsafe.Pointer(&origin))

	var s string = "123"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // case 1
	hdr.Data = oData.Data                              // case 6 (this case)
	hdr.Len = len(origin)

	log.Printf("%s", s)
}
