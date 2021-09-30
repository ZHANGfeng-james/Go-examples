Go 的 reflect 标准库实现的是**运行时反射**机制，允许应用程序操作（修改）任意类型实例。reflect 包的典型使用方式：使用**静态类型 `interface{}` 装载**一个值，调用 reflect.TypeOf 函数获取其**动态类型信息**，其返回值是**一个 Type 实例**。调用 reflect.ValueOf 函数可以返回**一个 Value 实例**，表示的是其**动态值信息**。如果返回的是零值，那么表示对应类型的零值。reflect 机制的通常用法是**和 `interface{}` 类相关的**，能获取其**动态 Type 和动态 Value**。

Go 语言中重要的函数和类型：

~~~go
ant@MacBook-Pro Go-examples-with-tests % go doc reflect |grep "^func"              
func Copy(dst, src Value) int
func DeepEqual(x, y interface{}) bool
func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool)
func Swapper(slice interface{}) func(i, j int)
~~~

类型：

~~~go
ant@MacBook-Pro Go-examples-with-tests % go doc reflect |grep "^type"|grep "struct"
type MapIter struct{ ... }
type Method struct{ ... }
type SelectCase struct{ ... }
type SliceHeader struct{ ... }
type StringHeader struct{ ... }
type StructField struct{ ... }
type Value struct{ ... }
type ValueError struct{ ... }

ant@MacBook-Pro Go-examples-with-tests % go doc reflect |grep "^type"|grep "interface"
type Type interface{ ... }
~~~

