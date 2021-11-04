package builtin

import (
	"reflect"
	"testing"
)

// builtin package
func TestCopySliceUseBuiltin(t *testing.T) {
	type param struct {
		input []int
	}

	params := []param{
		{[]int{1, 2, 3, 4}},
		{[]int{0}},
		{[]int{1, 2}},
	}

	for i := 0; i < len(params); i++ {
		input := params[i].input
		count, dst := copySliceUseBuiltin(input)
		// log.Println(count, tools.SliceInfo("result", dst))
		if count != len(input) || !reflect.DeepEqual(input, dst) {
			t.Fatal("copy slice use Builtin function error")
		}
	}
}

func TestSliceCallFunc(t *testing.T) {
	sliceCallFunc()
}

func TestSliceMake(t *testing.T) {
	createSliceUseMake(3)
}

func TestSliceAddr(t *testing.T) {
	getSliceAddr()
}

func TestSliceSizeof(t *testing.T) {
	getSliceSizeof()
}

func TestSliceAndArrayAddr(t *testing.T) {
	sliceAndArrayAddr()
}

func TestSliceAppend(t *testing.T) {
	sliceAppend()
}

func TestSliceRead(t *testing.T) {
	readEleFromSlice()
}

func TestSliceGrow(t *testing.T) {
	sliceGrow()
}

func TestSliceGrowTest(t *testing.T) {
	sliceGrowTest()
}

func TestSliceAgagin(t *testing.T) {
	sliceAgain()
}

func TestSliceNil(t *testing.T) {
	nilSlice()
}

func TestSliceConcurrent(t *testing.T) {
	sliceConcurrent()
	sliceConcurrentMutex()
}
