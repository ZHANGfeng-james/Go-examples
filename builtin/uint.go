package builtin

import (
	"log"
	"sync/atomic"
)

func uint64Value() {
	// uint64
	var val uint64
	val += 2 << 32
	val += 100
	log.Printf("val:%d", val)

	var delta int
	delta = -1
	atomic.AddUint64(&val, uint64(delta<<32))

	high32 := int32(val >> 32)
	low32 := uint32(val)
	log.Printf("hight32:%d, low32:%d", high32, low32)
}
