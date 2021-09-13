package gostrings

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestBuiltinCopy(t *testing.T) {
	var dst []byte
	if dst == nil {
		fmt.Println("dst is nil!")
	}
	n := copy(dst, "Michoi")
	fmt.Println("Success copy byte size:", n)

	fmt.Println(string(dst))
	fmt.Printf("%q\n", dst)
}

func TestReaderInterface(t *testing.T) {
	reader := strings.NewReader("")
	p := make([]byte, 0)
	n, err := reader.Read(p)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}

func TestReaderAtInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")
	p := make([]byte, 10)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		if err == io.EOF {
			fmt.Println("r.s --> io.EOF")
		} else {
			t.Fatal(err.Error())
		}
	}
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)

	reader.Read(p)
	fmt.Println(string(p))
}

func TestSeekerInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")
	p := make([]byte, 10)

	abs, err := reader.Seek(-4, io.SeekEnd)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("abs:", abs)

	n, _ := reader.Read(p)
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}

func TestWriterToInterface(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	reader := strings.NewReader("Michoi")

	n, err := reader.WriteTo(buffer)
	if err != nil {
		if err == io.ErrShortWrite {
			fmt.Println("r.s --> io.ErrShortWrite")
		} else {
			t.Fatal(err.Error())
		}
	}
	fmt.Println("n:", n)
	fmt.Println(buffer.String())
}

func TestByteScannerInterface(t *testing.T) {
	reader := strings.NewReader("Michoi")

	err := reader.UnreadByte()
	if err != nil {
		fmt.Println(err.Error())
	}

	p := make([]byte, 10)
	n, err := reader.Read(p)
	fmt.Printf("Success Read size:%d; p:%q\n", n, p)
}

func TestRuneScanner(t *testing.T) {
	reader := strings.NewReader("中\xc5国")

	for {
		ch, size, err := reader.ReadRune()
		if err != nil {
			t.Fatal(err.Error())
		}

		if !utf8.ValidRune(ch) {
			fmt.Println("InvalidRune")
		}

		fmt.Printf("%q, %d\n", ch, size)
	}
}

func TestByteSlice(t *testing.T) {
	buffer := []byte("Michoi")
	fmt.Printf("len(buffer)=%d, %q\n", len(buffer), buffer)

	start := -1
	slice := buffer[start:1]
	fmt.Printf("%q\n", slice)
}

func TestString(t *testing.T) {
	value := "Michoi"
	slice := value[:2]
	fmt.Printf("%T, %v\n", slice, slice)
}
