package builtin

import "testing"

func TestEmptyStructPointerEqual(t *testing.T) {
	createEmptyStructPointer()
}

func TestEmptyStructEqual(t *testing.T) {
	createEmptyStructVariable()
}

func TestNoEmptyStructEqual(t *testing.T) {
	compareNoEmptyStruct()
}

func TestStructMemoryAllocation(t *testing.T) {
	structMemoryAllocation()
}
