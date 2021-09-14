package gobytes

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"unicode"
)

func TestTrim(t *testing.T) {
	fmt.Printf("%q\n", bytes.Trim([]byte("michoim"), "mio"))
	fmt.Printf("[%q]\n", bytes.Trim([]byte(" !!! Achtung! Achtung! !!! "), "! "))
	fmt.Printf("%q\n", bytes.Trim([]byte("▓▽○中国国中○◇□▷★■●▽▒░▓"), "▽中▓"))
}

func TestTrimFunc(t *testing.T) {
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("\"go-gopher!\""), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}

func TestTrimLeft(t *testing.T) {
	fmt.Printf("%q\n", bytes.TrimLeft([]byte("michoim"), "mio"))
	fmt.Printf("[%q]\n", bytes.TrimLeft([]byte(" !!! Achtung! Achtung! !!! "), "! "))
	fmt.Printf("%q\n", bytes.TrimLeft([]byte("▓▽○中国国中○◇□▷★■●▽▒░▓"), "▽中▓"))
}

func TestTrimLeftFunc(t *testing.T) {
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher!"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("\"go-gopher!\""), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}

func TestTrimPrefixAndSuffix(t *testing.T) {
	var b = []byte("Hello, goodbye, etc!")
	b = bytes.TrimSuffix(b, []byte("goodbye, etc!"))
	b = bytes.TrimSuffix(b, []byte("gopher"))
	b = append(b, bytes.TrimSuffix([]byte("world!"), []byte("x!"))...)
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))

	var value = []byte("MichoiMi")
	os.Stdout.Write(bytes.TrimPrefix(value, []byte("Mio")))
}

func TestTrimSpace(t *testing.T) {
	values := " \t\n a lone gopher \n\t\r\n"
	fmt.Printf("%s.\n", bytes.TrimSpace([]byte(values)))

	fmt.Printf("%s.\n", bytes.TrimFunc([]byte(values), unicode.IsSpace))
}
