package gobytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	value := "Michoi"
	contains := bytes.Contains([]byte(value), []byte("Mi"))
	fmt.Printf("%v\n", contains)

	value = "seefood"
	contains = bytes.Contains([]byte(value), []byte("bar"))
	fmt.Printf("%v\n", contains)

	value = ""
	contains = bytes.Contains([]byte(value), []byte(""))
	fmt.Printf("%v\n", contains)

	value = "Michoi"
	contains = bytes.Contains(nil, []byte(""))
	fmt.Printf("%v\n", contains)

	contains = bytes.Contains(nil, nil)
	fmt.Printf("%v\n", contains)
}

func TestContainsAny(t *testing.T) {
	value := "Michoi"
	contains := bytes.ContainsAny([]byte(value), "abo")
	fmt.Printf("%v\n", contains)

	fmt.Println(bytes.ContainsAny([]byte("I like seafood."), "fÄo!"))
	fmt.Println(bytes.ContainsAny([]byte("I like seafood."), "去是伟大的."))
	fmt.Println(bytes.ContainsAny([]byte("I like seafood."), ""))
	fmt.Println(bytes.ContainsAny([]byte(""), ""))

	fmt.Println(bytes.ContainsAny(nil, ""))
}

func TestContainsRune(t *testing.T) {
	slice := []byte{0x00, 0xE4, 0xB8, 0xAD, 0xff}
	fmt.Println(bytes.ContainsRune(slice, '中'))

	fmt.Println(bytes.ContainsRune([]byte("I like seafood."), 'f'))
	fmt.Println(bytes.ContainsRune([]byte("I like seafood."), 'ö'))
	fmt.Println(bytes.ContainsRune([]byte("去是伟大的!"), '大'))
	fmt.Println(bytes.ContainsRune([]byte("去是伟大的!"), '!'))
	fmt.Println(bytes.ContainsRune([]byte(""), '@'))
}
