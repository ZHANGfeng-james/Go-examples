



~~~go
type string string
~~~

string 在 Go 中是**一种 builtin 内置的数据结构**，**代表了 8 位字节字符串的集合**，通常但不一定代表的是 UTF-8 编码的文本内容。一个 string 类型的变量内容可能是空的，但并不是 nil。另外一个特点是：string 类型值是**不可变** immutable 的。

runtime.go 中 string 在**运行时**对应的类型是：

~~~go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
~~~

根据字符串创建 string：

~~~go
//go:nosplit
func gostringnocopy(str *byte) string {
	ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)}
	s := *(*string)(unsafe.Pointer(&ss))
	return s
}
~~~

先创建的是 stringStruct，然后转化成 string。



**[]byte 值转化成 string 类型值**：

~~~go
package main

import "fmt"

func main() {
	bytes := []byte("123")
	fmt.Println(bytes)

	bytes_str := string(bytes) // []byte->string
	fmt.Println(bytes_str)

	for i := range bytes_str {
		fmt.Println(bytes_str[i])
	}

	bytes[0] = byte(50)           // 修改 original 字节数组内容
	fmt.Println(bytes_str, bytes) // []byte->string string内容并没有改变
}

[49 50 51]
123
49
50
51
123 [50 50 51] []byte->string string内容并没有改变
~~~

[]byte 值转化成 string 类型值的方法：`string(bytes)`

整个转化过程最关键的一个特点是：**需要执行一次内存拷贝**，原理是这样的：

1. 根据 []byte 申请到内存；
2. 构建 string；
3. **执行内存拷贝**：将 []byte 的内容完全拷贝到 string 申请到的底层内存中。

