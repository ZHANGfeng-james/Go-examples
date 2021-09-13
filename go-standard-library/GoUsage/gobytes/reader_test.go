package gobytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBufferInit(t *testing.T) {
	buffer := &bytes.Buffer{}
	fmt.Printf("%d, %d\n", buffer.Len(), buffer.Cap())

	buffer.Grow(10)
	fmt.Printf("%d, %d\n", buffer.Len(), buffer.Cap())
}

func TestBufferReader(t *testing.T) {
	origin := []byte("foo")
	var buffer *bytes.Buffer = bytes.NewBuffer(origin)
	fmt.Printf(" %s\n", buffer.String())

	value := make([]byte, 2)
	n, err := buffer.Read(value)
	fmt.Printf("%d, %s, %s\n", n, value, buffer.String())
	check(err, t)

	origin[2] = 'T'
	fmt.Println(string(origin))
	fmt.Printf(" %s\n", buffer.String())
}

func TestBufferWriter(t *testing.T) {
	var buffer *bytes.Buffer = &bytes.Buffer{}

	content := "Michoi"
	n, err := buffer.Write([]byte(content))
	if n < len(content) {
		t.Fatal(err.Error())
	}
	fmt.Printf("%d, %s, %s\n", n, content, buffer.String())
}

func TestBufferByteReader(t *testing.T) {
	buffer := bytes.NewBufferString("Michoi")
	char, _ := buffer.ReadByte()
	fmt.Printf("%q\n", char)
}

func TestBufferByteScanner(t *testing.T) {
	buffer := bytes.NewBufferString("Michoi")
	char, _ := buffer.ReadByte()
	fmt.Printf("%q\n", char)

	buffer.UnreadByte()

	char, _ = buffer.ReadByte()
	fmt.Printf("%q\n", char)
}

func TestBufferRuneReader(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	r, n, _ := buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)
}

func TestBufferRuneScanner(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	r, n, _ := buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)

	buffer.UnreadRune()

	r, n, _ = buffer.ReadRune()
	fmt.Printf("%q, %d\n", r, n)
}

func TestBufferByteWriter(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	err := buffer.WriteByte('\x41')
	if err == nil {
		fmt.Println(buffer.String())
	}
}

func TestBufferStringWriter(t *testing.T) {
	buffer := bytes.NewBufferString("中Michoi国")
	buffer.WriteString("人，在中国")
	fmt.Println(buffer.String())
}

func TestBufferReaderFrom(t *testing.T) {
	buffer := &bytes.Buffer{}
	reader := bytes.NewReader([]byte("Michoi"))

	n, _ := buffer.ReadFrom(reader)
	fmt.Printf("%d, %s\n", n, buffer.String())
}

func TestBufferWriterTo(t *testing.T) {
	buffer := bytes.NewBufferString("Michio")
	writer := &bytes.Buffer{}
	n, _ := buffer.WriteTo(writer)
	fmt.Printf("%d, %s\n", n, writer.String())
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestNewLineContent(t *testing.T) {
	newLine := "\r\n"
	fmt.Println([]byte(newLine))
}
