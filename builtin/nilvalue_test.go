package builtin

import (
	"log"
	"testing"
)

func TestBeanCallMethod(t *testing.T) {
	pointerTest()
}

func normalTest() {
	var ptr Bean
	ptr.callMethod()
}

func pointerTest() {
	var bean *Bean
	if bean == nil {
		log.Println("bean is nil")
		bean.callMethod()
	}
}
