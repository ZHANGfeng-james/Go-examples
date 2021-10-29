package builtin

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func byteSlice2string() {
	bytes := []byte("123")
	fmt.Println(bytes)

	bytes_str := string(bytes)
	fmt.Println(bytes_str)

	for i := range bytes_str {
		fmt.Println(bytes_str[i])
	}

	bytes[0] = byte(50)           // 修改 original 字节数组内容
	fmt.Println(bytes_str, bytes) // []byte->string string内容并没有改变
	log.Printf("bytes:%s, str:%s", bytes, bytes_str)
}

func stringCopy() {
	str1 := "123"
	str2 := str1
	log.Printf("str1:%s, str2:%s", str1, str2)

	// 只是获取到了 str1 和 str2 变量的地址
	log.Printf("%p, %p", &str1, &str2)
	// unsafe.Pointer --> uintptr
	dst1 := (uintptr)(unsafe.Pointer(&str1))
	dst2 := (uintptr)(unsafe.Pointer(&str2))
	log.Printf("dst1:%#x, dst2:%#x", dst1, dst2)

	hdr1 := (*reflect.StringHeader)(unsafe.Pointer(&str1))
	hdr2 := (*reflect.StringHeader)(unsafe.Pointer(&str2))
	log.Printf("hdr1:%#x, hdr2:%#x", hdr1.Data, hdr2.Data)

	str1 = "234"
	log.Printf("dst1:%#x, dst2:%#x", dst1, dst2)
	log.Printf("hdr1:%#x, hdr2:%#x", hdr1.Data, hdr2.Data)

	log.Printf("%s, %s", str1, str2)
}
