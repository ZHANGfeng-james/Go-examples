package gobytes

import (
	"bytes"
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestIndexNormal(t *testing.T) {
	index := bytes.Index([]byte("Michoi"), []byte("hoi"))
	fmt.Printf("index = %d\n", index)

	index = bytes.Index([]byte("chicken"), []byte("Michoi"))
	fmt.Printf("index = %d\n", index)
}

func TestIndexAny(t *testing.T) {
	index := bytes.IndexAny([]byte("Mi国cho中i"), "中国")
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexAny([]byte("chicken"), "aeiouy")
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexAny([]byte("crwth"), "aeiouy")
	fmt.Printf("index = %d\n", index)
}

func TestIndexByte(t *testing.T) {
	index := bytes.IndexByte([]byte("Michoi"), byte('i'))
	fmt.Printf("index = %d\n", index)

	index = bytes.IndexByte([]byte("Michoi"), byte('a'))
	fmt.Printf("index = %d\n", index)

	character := 'i'
	fmt.Printf("%T, %q\n", character, character)
}

func TestIndexFunc(t *testing.T) {
	satisfyFunc := func(r rune) bool {
		fmt.Printf("%q\n", r)
		return unicode.Is(unicode.Han, r)
	}

	index := bytes.IndexFunc([]byte("MichoiЛ中国"), satisfyFunc)
	fmt.Printf("index = %d\n", index)
}

func TestIndexRune(t *testing.T) {
	index := bytes.IndexRune([]byte{0x01, 0x02, 0xff, 0xee, 0x2f, 0x33}, utf8.RuneError)
	fmt.Printf("index = %d\n", index)
}

func TestIndexLast(t *testing.T) {
	index := bytes.LastIndex([]byte("Michoic"), []byte("ic"))
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndex([]byte("chicken"), []byte("Michoi"))
	fmt.Printf("index = %d\n", index)

	index = bytes.LastIndexAny([]byte("Mi国cho中i"), "中国")
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexAny([]byte("chicken"), "aeiouy")
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexAny([]byte("crwth"), "aeiouy")
	fmt.Printf("index = %d\n", index)

	index = bytes.LastIndexByte([]byte("Michoi"), byte('i'))
	fmt.Printf("index = %d\n", index)
	index = bytes.LastIndexByte([]byte("Michoi"), byte('a'))
	fmt.Printf("index = %d\n", index)

	satisfyFunc := func(r rune) bool {
		fmt.Printf("%q\n", r)
		return unicode.Is(unicode.Han, r)
	}
	index = bytes.LastIndexFunc([]byte("MichoiЛ中国"), satisfyFunc)
	fmt.Printf("index = %d\n", index)
}
