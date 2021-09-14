package gobytes

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	value := "中"
	fmt.Printf("%T, len(value)=%d\n", value, len(value))
	bytesSlice := []byte("中")
	fmt.Printf("%T, len(bytesSlice)=%d\n", bytesSlice, len(bytesSlice))

	values := "中\xcc"
	fmt.Printf("%T, len(value)=%d, %q\n", values, len(values), values)

	a := '中'
	b := 'b'
	c := byte('c')
	d := '\x01'
	fmt.Printf("%T, %q\n", a, a)
	fmt.Printf("%T, %q\n", b, b)
	fmt.Printf("%T, %q\n", c, c)
	fmt.Printf("%T, %q\n", d, d)
}

func TestRuneSlice(t *testing.T) {
	value := "中国加油！\xc3\x8c"
	runes := []rune(value)
	for index, ele := range runes {
		fmt.Printf("%d --> %q; ", index, ele)
	}
	fmt.Println()

	bytesEle := []byte(value)
	for i, element := range bytesEle {
		fmt.Printf("%d --> %q; ", i, element)
	}
	fmt.Println()

	fmt.Printf("%q\n", 0x4E2D)
	fmt.Printf("%q\n", '\u4E2D')
	fmt.Printf("%q\n", '\u00cc')
	fmt.Printf("%q\n", '\xcc')
	// fmt.Printf("%q\n", '\x4E2D')
}

func TestRuneNormal(t *testing.T) {
	fmt.Printf("%q\n", 0x4E2D)
	fmt.Printf("%q\n", '\u4E2D')
	fmt.Printf("%q\n", '\U00004E2D')

	fmt.Printf("%q\n", '\u00cc')
	fmt.Printf("%q\n", '\xcc')
	fmt.Printf("%q\n", '\314')
	// fmt.Printf("%q\n", '\x4E2D') illegal
}

func TestStringNormal(t *testing.T) {
	value := "中Michoi国"
	fmt.Printf("%q\n", value[0])

	for index, ele := range value {
		fmt.Printf("%d, %q\t", index, ele)
	}
	fmt.Println()

}

func TestStringToSlice(t *testing.T) {
	value := "中Michoi国"
	slice := []byte(value)

	slice[3] = 'm'
	fmt.Println(value)

	origin := string(slice)
	fmt.Println(origin)

	slice[4] = 'I'
	fmt.Println(origin)
}
