# Constants

常量，其属性在于：其值永远不会改变！

Go 中有各种不同的**常量**，比如：rune 常量、整型常量、浮点型常量、复数常量和字符串常量。其中，rune、整型、浮点型、复数型常量被归类为**数值型常量**。（rune 为什么会被归结为数值型常量？难道是因为 rune 是 int32 的**类型别名**？）

~~~go
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
~~~

一个**常量值**可以使用下面方式获取到：

* rune 类型、整型、浮点型、复数型、字符串型字面值；
* 使用常量限定的标识符；
* 常量表达式；
* 结果为常量的转化表达式；
* 某些 Go 语言内置函数比如 unsafe.Sizeof 的值；
* 应用于表达式的 cap 或 len 函数计算获得的值；
* real 或 imag 应用在复数常量计算得到的值；
* complex 应用在数值常量得到的值；
* 布尔真值由预先声明的常量 true 和 false 表示的值；
* 预声明的标识符 iota 表示的整数常量值。

作为示例，来看看 Go 标准库源代码中是如何定义常量的：

~~~go
// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	const max int = 1e6
	return x > max || x < -max
}
~~~

max 就是一种“使用常量限定的标识符”形式定义的常量值，定义 max 的类型是 int 类型，其值永远是 1e6 ——1,000,000。

~~~go
// Strings for use with buffer.WriteString.
// This is less overhead than using buffer.Write with byte arrays.
const (
	commaSpaceString  = ", "
	nilAngleString    = "<nil>"
	nilParenString    = "(nil)"
	nilString         = "nil"
	mapString         = "map["
	percentBangString = "%!"
	missingString     = "(MISSING)"
	badIndexString    = "(BADINDEX)"
	panicString       = "(PANIC="
	extraString       = "%!(EXTRA "
	badWidthString    = "%!(BADWIDTH)"
	badPrecString     = "%!(BADPREC)"
	noVerbString      = "%!(NOVERB)"
	invReflectString  = "<invalid reflect.Value>"
)
~~~

问题在于上述两种定义有和区别？而且第二种是没有给定类型的。类似的，看看下面的示例代码输出：

~~~go
package main

import (
	"fmt"
)

const (
	a uint = 1
	b      = 2
	c      = 3
	d      = 4
)

func main() {
	printlnTypeAndValue(a)
	printlnTypeAndValue(b)
	const max uint8 = 255
	printlnTypeAndValue(max)
}

func printlnTypeAndValue(value interface{}) {
	fmt.Printf("%T, %v.\n", value, value)
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
uint, 1.
int, 2.
uint8, 255.
~~~





一般来说，复数常量是常量表达式的一种。

数值型常量表示任意精度的精确值，并且不会溢出。相应的，没有用来表达 IEEE-754 标准中的  -0、无穷大值，以及 NaN 值。





**常量可能是有类型的，也有可能是无类型的**。字面值常量、true、false、iota，以及仅包含未类型化常量操作数的某些常量表达式是**未类型化的**（也就是没有类型的）。

一个常量可以通过下面的方式被赋予类型（**由无类型的值转化为有类型的值**）：显式的常量声明、转换；或者是隐式的变量声明或赋值，或者作为表达式中的操作数。如果这个常量值不能表示（represented）为与之对应的（respective）类型的值，则是错误的。

在需要类型的上下文中，**无类型的常量**具有**默认类型**（该类型会依据上下文做隐式类型转化）。比如，在短变量声明中没有显示类型的情况下，类似 `i := 0`。无类型的常量的默认类型分别是 bool、rune、int、float64、complex128 或 string，取决于该计算值是 boolean、rune、integer、float-point、complex、string 常量。

特别注意这些类型：bool、rune、int、float64、conplex128 和 string。因为这些类型对应的都是无类型值底层的默认类型（The default type of a untyped  constant...）！

实现约束：虽然数值型常量在语言层面上有任意的精度，但是编译器会使用带有限制精度的内部表示方式实现。这意味着每一个实现必须：

* 用至少 256 位表示整数常量；
* 表示浮点常数，包括复数常数的部分，尾数（mantissa）至少为 256 位，带符号的二进制指数至少为 16 位；
* 如果编译器无法精确表示整数常量，则会给出错误；
* 如果编译器因为溢出而无法表示浮点数或者复数常量，则会给出错误；
* 如果由于精度限制而无法表示浮点数或复数常数，则四舍五入到最接近的可表示的常数。

上述要求不仅应用于字面值常量，也对常量表达式的计算结果有效。

