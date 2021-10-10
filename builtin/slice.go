package builtin

import (
	"log"

	"github.com/go-examples-with-tests/tools"
)

func getSliceInfo() {
	origin := make([]int, 0)
	origin = append(origin, 1, 2, 3, 4, 5, 6)
	log.Print(tools.SliceInfo("origin", origin))

	dst := make([]int, len(origin), len(origin)) // len(dst) should not be zero!
	log.Print(tools.SliceInfo("dst", dst))

	num := copy(dst, origin) // copy() return the minmum of len(dst) and len(src)
	log.Printf("copy num:%d", num)
	log.Print(tools.SliceInfo("dst", dst))
}

// builtin function:[func copy(dst, src []Type) int]
func copySliceUseBuiltin(src []int) (int, []int) {
	result := make([]int, len(src), len(src))
	num := copy(result, src)
	return num, result
}
