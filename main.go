package main

import (
	"log"
	"sync/atomic"
	"unsafe"
)

type T3 struct {
	b int64
	c int32
	d int64
}

func main() {
	t3 := &T3{}
	atomic.AddInt64(&t3.d, 1)
	log.Println(unsafe.Sizeof(*t3))
}
