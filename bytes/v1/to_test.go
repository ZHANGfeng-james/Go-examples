package gobytes

import (
	"bytes"
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestToLower(t *testing.T) {
	value := "MiChoi"
	fmt.Printf("%q.\n", bytes.ToLower([]byte(value)))
	fmt.Printf("%q.\n", bytes.ToLower([]byte("没有大小写的 Unicode 码")))
}

func TestToLowerSpecial(t *testing.T) {
	str := []byte("AHOJ VÝVOJÁRİ GOLANG")
	// unicode.AzeriCase / unicode.TurkishCase
	totitle := bytes.ToLowerSpecial(unicode.TurkishCase, str)
	fmt.Println("Original : " + string(str))
	fmt.Println("ToLower : " + string(totitle))
}

func TestToTitle(t *testing.T) {
	fmt.Printf("%s\n", bytes.ToTitle([]byte("loud noises")))
	fmt.Printf("%s\n", bytes.ToTitle([]byte("хлеб")))
}

func TestToTitleSpecial(t *testing.T) {
	str := []byte("ahoj vývojári golang")
	totitle := bytes.ToTitleSpecial(unicode.AzeriCase, str)
	fmt.Println("Original : " + string(str))
	fmt.Println("ToTitle : " + string(totitle))
}

func TestToUpper(t *testing.T) {
	fmt.Printf("%s\n", bytes.ToUpper([]byte("Gopher")))

	str := []byte("ahoj vývojári golang")
	totitle := bytes.ToUpperSpecial(unicode.AzeriCase, str)
	fmt.Println("Original : " + string(str))
	fmt.Println("ToUpper : " + string(totitle))
}

func TestDiffUpperAndTitle(t *testing.T) {
	str := "ǳ"
	fmt.Printf("%q\n", bytes.ToTitle([]byte(str)))
	fmt.Printf("%q\n", bytes.ToUpper([]byte(str)))

	fmt.Printf("%s\n", bytes.ToUpper([]byte("Gopher")))
	fmt.Printf("%s\n", bytes.ToTitle([]byte("Gopher")))
}

func TestToValidUTF8(t *testing.T) {
	value := "\xc5Geeks\xc5Geeks\n"
	fmt.Println(value)
	origin := []byte(value)
	fmt.Printf("%0x\n", origin)

	// Invalid UTF-8 '\xc5' replaced by 'For'
	dst := bytes.ToValidUTF8(origin, []byte("TT"))
	fmt.Printf("%q\n", dst)
}

func TestUnicodeUTF(t *testing.T) {
	code := '\u00c5' // rune 类型
	fmt.Printf("%T, %q, %#U\n", code, code, code)
	code = '\xc5' // rune 类型
	fmt.Printf("%T, %q, %#U\n", code, code, code)

	// rune --> []byte
	fmt.Printf("%0x\n", code)

	fmt.Println(utf8.ValidRune(code))

	fmt.Println(utf8.Valid([]byte("\xc5")))
}

func TestUTF8Valide(t *testing.T) {
	fmt.Println(utf8.Valid([]byte("\xc5")))
}
