package bytes

import (
	"bytes"
	"log"
)

func bytesUsage() {
	buffer := new(bytes.Buffer)

	slice := make([]byte, 128)

	buffer.WriteString(string(slice))
	log.Println(buffer.Len(), buffer.Cap())

	buffer.Reset()
	log.Printf("reset len:%d, cap:%d", buffer.Len(), buffer.Cap())
}
