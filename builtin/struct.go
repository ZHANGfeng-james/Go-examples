package builtin

import (
	"fmt"
	"log"
	"unsafe"
)

type Person struct{}

func createEmptyStructPointer() {
	a := &Person{}
	b := &Person{}

	// 逃逸
	// log.Printf("a'addr:%p, b'addr:%p", a, b)
	log.Printf("a == b? --> %v", a == b)

	// 出现逃逸时，true
	// 未出现逃逸时，false，但是指向的是同一个底层对象
}

func createEmptyStructVariable() {
	a := Person{}
	b := Person{}

	log.Printf("a == b? --> %v", a == b)

	// 涉及到 struct 变量的比较，引出的问题是：struct 变量能否做 map 的 key？
}

type NoEmptyStruct struct {
	name string
	age  int
}

func NewNoEmptyStruct(name string, age int) *NoEmptyStruct {
	return &NoEmptyStruct{
		name: name,
		age:  age,
	}
}

func (s *NoEmptyStruct) String() string {
	return "name:" + s.name + "; age:" + fmt.Sprint(s.age)
}

func compareNoEmptyStruct() {
	a := &NoEmptyStruct{}
	b := &NoEmptyStruct{}

	log.Printf("a'addr:%p, b'addr:%p", a, b)

	log.Printf("a.name:%s, b.name:%s", a.name, b.name)

	log.Printf("a == b? --> %v", a == b)
}

type t1 struct {
	a int8

	b int64

	c int16
}

type t2 struct {
	a int8

	c int16

	b int64
}

func structMemoryAllocation() {
	v1Struct := t1{}
	v2Struct := t2{}
	log.Printf("Memory Allocation: %d, %d", unsafe.Sizeof(v1Struct), unsafe.Sizeof(v2Struct))
}
