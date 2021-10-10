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
