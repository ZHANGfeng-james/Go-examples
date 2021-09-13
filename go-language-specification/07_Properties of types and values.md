Properties of types and values 意为：类型和值的**属性特征**

类型和值的属性特征，其含义就是：

* 类型校验
* 赋值性
* 可表示性

# 1 Type identity

Type identity 类型一致性

两种类型要么是相同的，要么就是不一样的！

定义的类型始终与任何其他类型不同，比如 `type T int`。反之，如果两个类型的基础类型字面值在结构上相同，则它们就是相同的类型。（**以基本类型作为区分，类型 A 和类型 B 是同一类型的前提条件**）详细情况如下：

* **数组**：元素类型相同，且数组长度相同；
* **切片**：元素类型相同；
* **结构体**：两者有相同的字段**顺序**，对应的字段有相同的名称、类型、标签。来自不同包的非导出字段名，会造成结构体类型不同；
* **指针**：具有相同的基本类型，比如 `*int` 和 `*int8` 就不是相同的类型；
* **函数**：具有相同数量的参数、返回值，且对应的类型都是一样的；参数名和返回值名称不影响；
* **接口**：相同的方法集合，包括相同的方法名和方法类型；
* **map**：相同的 `key-元素` 类型；
* **channel**：相同的元素类型和方向。

给定如下的**声明（declaration）**：

~~~go
type ( // 类型重命名
	A0 = []string
	A1 = A0
	A2 = struct{ a, b int }
	A3 = int
	A4 = func(A3, float64) *A0
	A5 = func(x int, _ float64) *[]string
)

type ( // 类型重定义
	B0 A0
	B1 []string
	B2 struct{ a, b int }
	B3 struct{ a, c int }
	B4 func(int, float64) *B0
	B5 func(x int, y float64) *A1
)

type	C0 = B0
~~~

下面这些类型是**相同的**（identical）：

~~~go
A0, A1, and []string
A2 and struct{ a, b int }
A3 and int
A4, func(int, float64) *[]string, and A5

B0 and C0
[]int and []int
struct{ a, b *T5 } and struct{ a, b *T5 }
func(x int, y float64) *[]string, func(int, float64) (result *[]string), and A5
~~~

类型 B0 和类型 B1 是两种不同的类型，因为他们是通过类型定义创建了新的类型。此外，`func(int, float64) *B0` 和 `func(x int, y float64) *[]string` 是两种不同的类型，因为 B0 和 `[]string` 是不同的类型。

特别重要的是：

~~~go
package main

import "fmt"

type A struct {
	a int
}

type B struct {
	a int
}

func main() {
	var a A
	var b B

	a = A(b)
	fmt.Print(a)
}
~~~

如果要让 b 赋值给变量 a，必须要经过一次类型转换：`a = A(b)`，否则会编译出错。另外，对于结构体 A 和 B，必须是其中的各个字段**名称、类型、标签、顺序**是相同的，才能够相互进行转化。

~~~go
package main

import "fmt"

type A struct {
	a int
	b int
}

type B struct {
	b int
	a int
}

func main() {
	var a A
	var b B

	a = A(b) // cannot convert b (variable of type B) to A
	fmt.Print(a)
}
~~~

即便是 struct 中的字段顺序不一致，则会编译错误。

# 2 Assignability

Assignability 可赋值性

值 x 可被赋值给类型 T 的变量，需满足如下条件：

* x 的类型和 T 类型相同；
* x 的类型 v 和 T 类型具有相同的底层类型，并且 v 和 T 至少有一个不是定义类型（defined type）；
* T 是一个接口类型，x 实现了该接口；
* x 是一个双向 channel 值，T 是一个 channel 类型，x 的类型 v 和 T 有相同的元素类型，并且 v 和 T 至少有一个不是定义类型（defined type）；
* x 是预声明标识符 nil，T 是一个指针、方法、slice、map、channel 或接口类型；
* x 是一个 T 类型的无类型常量值。

上述情况下，无需进行强制类型转化，就可直接赋值。

# 3 Representability

representability 可表示性

如果满足以下条件之一，则常量 x 可以用类型 T 的值表示：

* x 的值在类型 T 的值的集合之中；
* T 是浮点类型，x 可以“四舍五入”而不会溢出地转化到 T 类型的精度。注意：常量值永远不可能为 -0、NaN 或无穷值。
* T 是复数类型，x 对应的实部和虚部能够被 T 对应的实部和虚部表示。

举例而言：

~~~go
x                   T           x is representable by a value of T because

'a'                 byte        97 is in the set of byte values
97                  rune        rune is an alias for int32, and 97 is in the set of 32-bit integers
"foo"               string      "foo" is in the set of string values
1024                int16       1024 is in the set of 16-bit integers
42.0                byte        42 is in the set of unsigned 8-bit integers
1e10                uint64      10000000000 is in the set of unsigned 64-bit integers
2.718281828459045   float32     2.718281828459045 rounds to 2.7182817 which is in the set of float32 values
-1e-1000            float64     -1e-1000 rounds to IEEE -0.0 which is further simplified to 0.0
0i                  int         0 is an integer value
(42 + 0i)           float32     42.0 (with zero imaginary part) is in the set of float32 values
~~~

反之：

~~~go
x                   T           x is not representable by a value of T because

0                   bool        0 is not in the set of boolean values
'a'                 string      'a' is a rune, it is not in the set of string values
1024                byte        1024 is not in the set of unsigned 8-bit integers
-1                  uint16      -1 is not in the set of unsigned 16-bit integers
1.1                 int         1.1 is not an integer value
42i                 float32     (0 + 42i) is not in the set of float32 values
1e1000              float64     1e1000 overflows to IEEE +Inf after rounding
~~~

比较特别的是：

~~~go
package main

import "fmt"

func main() {
	fmt.Printf("%T.\n", 42.0)
	var value byte
	value = 42.0
	fmt.Print(value)
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
float64.
42
~~~

对于**无类型常量** 42.0，其类型确实是 float64；可将其转化为 byte 类型，是因为 42 这个值确实是在 byte 类型的范围之内。

