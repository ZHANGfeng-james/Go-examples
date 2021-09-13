package gobytes

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
	"unicode"
)

func TestCompare(t *testing.T) {
	// Binary search to find a matching byte slice.
	var needle []byte = []byte("d")

	var haystack [][]byte // Assume sorted
	haystack = append(haystack, []byte("abc"))
	haystack = append(haystack, []byte("abd"))
	haystack = append(haystack, []byte("d"))
	haystack = append(haystack, []byte("efg"))

	i := sort.Search(len(haystack), func(i int) bool {
		// Return haystack[i] >= needle.
		return bytes.Compare(haystack[i], needle) >= 0
	})
	if i < len(haystack) && bytes.Equal(haystack[i], needle) {
		fmt.Println("Found it! index:", i)
	}

	var a []byte = []byte("0xFE")
	var b []byte = []byte("0xff")
	fmt.Println(bytes.Compare(a, b))
}

func TestCount(t *testing.T) {
	value := "中国\xE4\xB8\xAD"
	fmt.Println(bytes.Count([]byte(value), nil))
	for index, ele := range []rune(value) {
		fmt.Printf("rune %d#, %q\n", index, ele)
	}

	Count("中国人来自中国", "中国")
	Count("中国人来自中国", "中国人")
	Count("众里寻他千百度，慕容回首，那人却在灯火阑珊处！", "百度")

	Count("five", "ve")
	Count("five", "v")
	Count("five", "")
}

func Count(s, sep string) {
	fmt.Println(bytes.Count([]byte(s), []byte(sep)))
}

func TestFields(t *testing.T) {
	Fields("  foo bar  baz   ")
	Fields("  foo bar baz")
	Fields("\tfoo \rbar    \nbaz")
	Fields("\t\r\n")
}

func Fields(s string) {
	result := bytes.Fields([]byte(s))
	fmt.Printf("%T, Fields are %q\n", result, result)
}

func TestFieldsFunc(t *testing.T) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Printf("Fields are %q\n", bytes.FieldsFunc([]byte("  foo1;bar2,baz3..."), f))
}

func TestHasPrefix(t *testing.T) {
	fmt.Println(bytes.HasPrefix([]byte("Michoi"), []byte("mi")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), []byte("中")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), []byte("")))
	fmt.Println(bytes.HasPrefix([]byte("中国人"), nil))

	fmt.Println(bytes.HasPrefix([]byte(""), nil))
	fmt.Println(bytes.HasPrefix(nil, nil))
}

func TestHasSuffix(t *testing.T) {
	fmt.Println(bytes.HasSuffix([]byte("Michoi"), []byte("oi")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), []byte("中")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), []byte("")))
	fmt.Println(bytes.HasSuffix([]byte("中国人"), nil))

	fmt.Println(bytes.HasSuffix([]byte(""), nil))
	fmt.Println(bytes.HasSuffix(nil, nil))
}

func TestJoin(t *testing.T) {
	var s [][]byte
	s = append(s, []byte("M"))
	s = append(s, []byte("i"))
	s = append(s, []byte("c"))
	s = append(s, []byte("h"))
	s = append(s, []byte("o"))
	s = append(s, []byte("i"))
	fmt.Printf("%q\n", bytes.Join(s, []byte("-")))
}

func TestMap(t *testing.T) {
	rot13 := func(r rune) rune {
		// F 之后的字符输出正值，F 之前的字符输出负值
		switch {
		case r >= 'F' && r <= 'Z':
			return 'A'
		case r >= 'f' && r <= 'z':
			return 'a'
		case r < 'F':
			return -1
		case r < 'f':
			return -1
		}

		return r
	}
	fmt.Printf("%s", bytes.Map(rot13, []byte("'agAG...")))
}

func TestRepeat(t *testing.T) {
	fmt.Printf("ba%s", bytes.Repeat([]byte("na"), 2))
}

func TestReplaceAll(t *testing.T) {
	s := "oink oink oink"
	old := "oink"
	new := "moo"
	ReplaceAll(s, old, new)

	s = "HiLink"
	old = ""
	new = "-"
	ReplaceAll(s, old, new)
}

func ReplaceAll(s, old, new string) {
	fmt.Printf("%s\n", bytes.ReplaceAll([]byte(s), []byte(old), []byte(new)))
}

func TestReplaceCount(t *testing.T) {
	s := "oink oink oink"
	old := "k"
	new := "ky"
	count := 0
	Replace(s, old, new, count)

	old = "oink"
	new = "moo"
	count = -1
	Replace(s, old, new, count)
}

func Replace(s, old, new string, count int) {
	fmt.Printf("%s\n", bytes.Replace([]byte(s), []byte(old), []byte(new), count))
}

func TestRunes(t *testing.T) {
	s := "中Michoi国\xcc, \u4E2D"

	runes := bytes.Runes([]byte(s))
	for i, ele := range runes {
		fmt.Printf("%d# %q; ", i, ele)
	}
	fmt.Println()
}

func TestTitle(t *testing.T) {
	fmt.Printf("%s", bytes.Title([]byte("her royal highness")))
}
