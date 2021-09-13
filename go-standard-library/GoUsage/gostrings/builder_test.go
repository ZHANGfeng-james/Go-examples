package gostrings

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuilder(t *testing.T) {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("ignition")
	b.WriteRune('中')
	b.WriteByte('\xcc')
	b.Write([]byte("Michoi"))
	fmt.Println(b.String(), "cap:", b.Cap(), "; len:", b.Len())

	b.Reset()
	fmt.Println(b.String())
}
