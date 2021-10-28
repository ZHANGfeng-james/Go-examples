package builtin

import "log"

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

type noEmptyStruct struct {
	name string
}

func compareNoEmptyStruct() {
	a := &noEmptyStruct{}
	b := &noEmptyStruct{}

	log.Printf("a'addr:%p, b'addr:%p", a, b)

	log.Printf("a.name:%s, b.name:%s", a.name, b.name)

	log.Printf("a == b? --> %v", a == b)
}
