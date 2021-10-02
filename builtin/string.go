package builtin

import "fmt"

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
}
