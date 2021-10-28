package builtin

import "log"

type Bean struct {
	value []int
}

func (b *Bean) callMethod() {
	log.Printf("callMethod is called, %T, %v", b, b)

	if b == nil {
		log.Println("b is nil")
	}

	log.Println("b.value:", b.value) // 此处如果是 pointerTest 会 Panic
}
