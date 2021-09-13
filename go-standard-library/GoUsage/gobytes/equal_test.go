package gobytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	fmt.Println(bytes.Equal([]byte(""), []byte("")))
	fmt.Println(bytes.Equal(nil, []byte("")))
	fmt.Println(bytes.Equal([]byte(""), nil))

	fmt.Println(bytes.Equal([]byte("Michoi"), []byte("michoi")))
	fmt.Println(bytes.Equal([]byte("Go"), []byte("Go")))
	fmt.Println(bytes.Equal([]byte("Go"), []byte("C++")))
}

func TestEqualFold(t *testing.T) {
	s := []byte{0xFF}
	tt := []byte{0xff}

	fmt.Println(bytes.EqualFold(s, tt))

	fmt.Println(bytes.EqualFold([]byte("Go"), []byte("go")))
}
