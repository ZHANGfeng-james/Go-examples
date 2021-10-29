package unsafe

import (
	"log"
	"testing"
)

func TestUnsafePointer(t *testing.T) {
	usageUnsafePointer()
}

func TestUnsafeSizeof(t *testing.T) {
	usageUnsafeSizeof()
}

func TestUnsafeGetAddr(t *testing.T) {
	var i int64 = 100
	// *int64 --> uintptr, uintptr(&i) error!
	log.Printf("%p, %#x", &i, getVarAddrUsingPointer(&i))
}

func TestUnsafeChangeString(t *testing.T) {
	changeStringContent()
}

func TestUnsafeChangeStruct(t *testing.T) {
	changeStructField()
}

func TestUnsafeChangeArray(t *testing.T) {
	changeArrayEle()
}

func TestUnsafeGetSliceInfo(t *testing.T) {
	getSliceHeaderInfo()
}

func TestUnsafeBytesToString(t *testing.T) {
	byteSliceToStringNoCopy()
}
