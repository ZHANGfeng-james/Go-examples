

**表达式**的含义：在**操作符**中使用**操作数**和**函数**，计算得到了**一个值**。这种形式的内容，就称为是一个表达式。

# 1 Operands

**操作数**，意味着的是表达式中的元素值。

一个操作数，可能是一个**字面值**，也可能是一个表示常量、变量、函数、括号表达式的非空白的标识符。空白标识符只能出现在赋值语句中的左侧，用来代表一个操作数。

~~~go
Operand     = Literal | OperandName | "(" Expression ")" .
Literal     = BasicLit | CompositeLit | FunctionLit .
BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
OperandName = identifier | QualifiedIdent .
~~~

# 2 Qualified identifiers

**限定标识符**，表示的是一个**带有包名前缀限定**的标识符。包名和标识符都必须不是空白标识符：

~~~
QualifiedIdent = PackageName "." identifier .
~~~

一个限定标识符可以访问不同包下的，**已被导入的**标识符。这个标识符必须是**可导出的**，而且声明在其所在包的包代码块中。

~~~go
math.Sin	// denotes the Sin function in package math
~~~

# 3 Composite literals

复合字面值为结构体、数组、切片、map 构造值，每次计算这个复合字面值时都会得到一个新建的值。它们由文字的类型，以及紧随其后的元素列表（使用 `{}` 包围着的）组成。每个元素前，可选择写上对应的 key 值。

~~~go
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
~~~

字面值类型（LiteralType）的潜在类型必须是：结构体、数组、切片或 map 类型（语法强制执行此约束，但类型是 TypeName 时除外）。此处的 key，意味着结构体类型的字段，或者意味着数组、切片类型的索引，或者是 map 类型的 key 值。

对于**结构体类型**，需要遵循如下规则：

* key 必须是结构体类型的字段名称；
* 元素列表中，如果没有给出任何 key 内容，那么就需要按照结构体的字段声明的顺序给出其字段值；
* 如果元素列表中，某个元素是带有 key 的，那么全部元素都必须带有 key。否则编译器会提示 `mixture of field:value and value elements in struct literal`；
* 元素列表中，省略的部分将会被赋予字段的默认值；
* 一个结构体字面值可能会省略元素类表，这种形式相当于是结构体字段的默认值；
* 为属于不同包的结构体的非导出字段指定元素是错误的。

~~~go
type Point3D struct { x, y, z float64 }
type Line struct { p, q Point3D }
~~~

结构体字面值，可以这样写：

~~~go
origin := Point3D{}                            // zero value for Point3D
line := Line{origin, Point3D{y: -4, z: 12.3}}  // zero value for line.q.x
~~~

对于**数组和切片字面值**来说，遵循如下规则：

* 每一个元素都对应有整型位置索引；
* 每个元素有一个对应的 key，这个 key 值必须是非负常量值（常整型值）；如果是重定义类型，则必须是整型类型；
* 如果首个元素没有对应的 key 值，默认是 0；后续的其他元素则在其基础之上自增 1。

获取一个复合字面值的地址操作（Taking the address），将生成与该唯一变量（用这个复合字面值创建的）相关的指针。

~~~go
var pointer *Point3D = &Point3D{y: 1000}
~~~

对于数组和切片而言，该类型的零值（）和类型相同但已被初始化为空的值是不一样的。

~~~go
var slice []int // slice is nil
var slice = make([]int, 0) // slice is not nil!
~~~

与之相关的，获取一个空的数组或切片复合字面值的地址，和使用 new 分配一个新的 slice 和 map 类型的指针是不一样的：

~~~go
p1 := &[]int{}    // p1 points to an initialized, empty slice with value []int{} and length 0
p2 := new([]int)  // p2 points to an uninitialized slice with value nil and length 0
~~~

区别在于 `*p2` 的值是 nil，也就是 `[]int` 类型的零值（未被初始化的）。

关于**数组字面值**的长度：

~~~go
buffer := [10]string{}             // len(buffer) == 10
intSet := [6]int{1, 2, 3, 5}       // len(intSet) == 6
days := [...]string{"Sat", "Sun"}  // len(days) == 2
~~~

上述数组字面值中的**注解** `...` 表示的含义是：元素的最大位置索引值 + 1，其值就是数组的长度。

一个切片字面值指出了其底层所有的数组字面值，因此切片字面值的长度就是所给出元素的最大位置索引值再加上1计算出的值：

~~~go
[]T{x1, x2, … xn}
~~~

上述的切片字面值的长度等于：`xn` 元素的索引值 + 1

另外还可以在数组值的基础上，通过下列操作**获得切片**：

~~~go
tmp := [n]T{x1, x2, … xn}
tmp[0 : n] // tmp[:2] 切片操作，而不是数组操作
~~~

在数组、切片和 map 复合字面值中，可以省略元素的类型：

~~~go
[...]Point{{1.5, -3.5}, {0, 0}}     // same as [...]Point{Point{1.5, -3.5}, Point{0, 0}}
[][]int{{1, 2, 3}, {4, 5}}          // same as [][]int{[]int{1, 2, 3}, []int{4, 5}}
[][]Point{{{0, 1}, {1, 2}}}         // same as [][]Point{[]Point{Point{0, 1}, Point{1, 2}}}
map[string]Point{"orig": {0, 0}}    // same as map[string]Point{"orig": Point{0, 0}}
map[Point]string{{0, 0}: "orig"}    // same as map[Point]string{Point{0, 0}: "orig"}

type PPoint *Point
[2]*Point{{1.5, -3.5}, {}}          // same as [2]*Point{&Point{1.5, -3.5}, &Point{}}
[2]PPoint{{1.5, -3.5}, {}}          // same as [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}
~~~

下面是正确的数组、切片和 map 字面值：

~~~go
// list of prime numbers
primes := []int{2, 3, 5, 7, 9, 2147483647}

// vowels[ch] is true if ch is a vowel
vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}

// the array [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}

// frequencies in Hz for equal-tempered scale (A4 = 440Hz)
noteFrequency := map[string]float32{
	"C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
	"G0": 24.50, "A0": 27.50, "B0": 30.87,
}
~~~

# 4 Function literals

一个函数字面值，表示的是一个**匿名函数**。

~~~go
FunctionLit = "func" Signature FunctionBody .
~~~

比如这样的匿名函数：

~~~go
func(a, b int, z float64) bool { return a*b < int(z) }
~~~

该函数是没有名字的——匿名函数的本意。

一个函数字面量可以赋值给一个变量或者直接调用：

~~~go
f := func(x, y int) int { return x + y }

replyChan := make(chan int, 1)
ACK := 0
func(ch chan int) { ch <- ACK }(replyChan)
~~~

函数字面量是一个**闭包**：在这个函数字面量中，**可能会引用该函数字面量定义区域周边的变量**；这些变量会被函数字面量和其周边函数所共享。该变量在可被访问的时期内一直会存活。

# 5 Primary expressions

主表达式，指的是一元操作数和二进制表达式的操作数

~~~go
PrimaryExpr =
	Operand |
	Conversion |
	MethodExpr |
	PrimaryExpr Selector |
	PrimaryExpr Index |
	PrimaryExpr Slice |
	PrimaryExpr TypeAssertion |
	PrimaryExpr Arguments .

Selector       = "." identifier .
Index          = "[" Expression "]" .
Slice          = "[" [ Expression ] ":" [ Expression ] "]" |
                 "[" [ Expression ] ":" Expression ":" Expression "]" .
TypeAssertion  = "." "(" Type ")" .
Arguments      = "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" .
~~~

举例：

~~~go
x
2
(s + ".txt")
f(3.1415, true)
Point{1, 2}
m["foo"]
s[i : j + 1]
obj.color
f.p[i].x()
~~~

# 6 Selectors

对于主表达式来说，x 并不是包名，**选择器表达式**：

~~~go
x.f
~~~

其含义是 x 值的字段或方法 f（x 的含义有可能表示的是 `*x`，如果 x 是指针的话）。**标识符 f 就叫做（字段或方法的）选择器**；该标识符不能为空白标识符。选择器表达式的类型和 f 的类型一致。如果 x 是包名，那么这个表达式就是限定标识符。

一个选择器 f 可能代表的是类型 T 的一个字段或方法，或者引用的是 T 的嵌入字段的字段或方法。

下述是有关选择器的规则：

* 对于类型 T 或者 `*T`（其中 T 不能是一个指针或接口类型）的值 x 而言，`x.f` 的含义是：类型 T 的字段或方法，潜在的，类型 T 确实存在该字段或方法 f；如果不存在，那么该选择器表达式是非法的。
* 对于接口类型 I 的值 x 而言，`x.f` 的含义是：动态类型值 x 的实际方法 f。如果在 I 的方法集合中不存在方法 f，那么这个表达式非法。
* 一个例外的实例，如果 x 的类型是一个定义的指针类型，`(*x).f` 是一个合法的选择器表达式，其含义是一个字段（不是一个方法），`x.f` 是 `(*x).f` 的简短形式。
* 除了上述的其他写法，`x.f` 都是非法的。
* 如果 x 是一个指针类型，其值是 nil，那么 `x.f` 代表的是一个结构体字段；为 `x.f` 赋值或者计算 `x.f` 都会引发 run-time panic。
* 如果 x 是一个接口类型，其值是 nil，那么调用或计算 `x.f` 方法会引发 run-time panic。

举例而言：定义了如下的结构体类型，及其对应的方法

~~~go
type T0 struct {
	x int
}

func (*T0) M0()

type T1 struct {
	y int
}

func (T1) M1()

type T2 struct {
	z int
	T1
	*T0
}

func (*T2) M2()

type Q *T2

var t T2     // with t.T0 != nil
var p *T2    // with p != nil and (*p).T0 != nil
var q Q = p
~~~

那么在这些类型基础之上的选择器表达式有：

~~~go
t.z          // t.z
t.y          // t.T1.y
t.x          // (*t.T0).x

p.z          // (*p).z
p.y          // (*p).T1.y
p.x          // (*(*p).T0).x

q.x          // (*(*q).T0).x        (*q).x is a valid field selector

p.M0()       // ((*p).T0).M0()      M0 expects *T0 receiver
p.M1()       // ((*p).T1).M1()      M1 expects T1 receiver
p.M2()       // p.M2()              M2 expects *T2 receiver
t.M2()       // (&t).M2()           M2 expects *T2 receiver, see section on Calls
~~~

特别的 `(*t.T0).x` 实际上就是：

1. `t.T0`
2. `*(t.T0)`
3. `(*(t.T0)).x`

验证如下：

~~~go
func Init() {
	var t T2 // with t.T0 != nil

	t = T2{
		z: 10,
		T1: T1{
			y: 20,
		},
		T0: &T0{
			x: 30,
		},
	}
	if t.T0 == nil {
		fmt.Println("nil")
	}

	t.x = 90
	(*t.T0).x = 80
	(*(t.T0)).x = 70
	fmt.Println(t.x)
}
~~~

但是对于如下是非法的：

~~~go
q.M0()       // (*q).M0 is valid but not a field selector
~~~

![](./pics/Snipaste_2021-05-21_16-15-45.png)

# 7 Method expressions

如果方法 M 是类型 T 方法集合中的一员，那么 T.M 可以作为一个被调用的普通方法：方法参数不仅包含相同的参数，另外还需要在参数列表的首位置增加一个作为方法接收者的参数。方法表达式的构成：

~~~go
MethodExpr    = ReceiverType "." MethodName .
ReceiverType  = Type .
~~~

如下所述的结构体类型 T 包含有 2 个方法：

~~~go
type T struct {
	a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
~~~

表达式 T.Mv 表示的就是一个和 Mv 相同的方法，但其包含有一个显式的接收者作为其参数：

~~~go
func(tv T, a int) int
~~~

这个方法可以这样被调用，下述调用方式都是等价的：

~~~go
t.Mv(7)
T.Mv(t, 7)
(T).Mv(t, 7)
f1 := T.Mv; f1(t, 7)
f2 := (T).Mv; f2(t, 7)
~~~

于此类似的，表达式：

~~~go
(*T).Mp
~~~

是一个方法值，表达的是 Mp：

~~~go
func(tp *T, f float32) float32
~~~

对于基本类型 T 的方法 Mv（类型 T 是其值接收者类型），可将其转化为显式的指针接收者：

~~~go
(*T).Mv
~~~

代表的是 Mv 方法值，其方法签名：

~~~go
func(tv *T, a int) int
~~~

作为验证：

~~~go
func function() {
	var t T
	ptr := &t

	fmt.Printf("%T.\n", (*T).Mp) // func(*method.T, float32) float32.
	fmt.Printf("%T.\n", (*T).Mv) // func(*method.T, int) int.

	(*T).Mp(ptr, 2)
}
~~~

但是这种现实是不合法的，因为 `func(tv *T, a int) int` 并不是方法集合中的一员。

# 8 Method values

如果表达式 x 的静态类型是 T，而且 T 类型的方法集合中有方法 M。那么 x.M 就叫做一个方法值。

方法值 x.M 是一个函数值，这个函数值可以类似于调用 x.M 一样（包含相同的参数）被调用。其中，表达式 x 会被计算，同时表达式 x 的值在计算方法值的过程中被保存下来；被保存下来的值的副本会作为方法调用的接收者。

类型 T 可能是一个接口，或者一个非接口类型。

~~~go
type T struct {
	a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
var pt *T
func makeT() T
~~~

表达式 `t.Mv`（不是 `T.Mv`）是一个函数值，其类型是：

~~~go
func(int) int
~~~

如下两种调用是等价的：

~~~go
t.Mv(7)
f := t.Mv; f(7)
~~~

向类似的，表达式 `pt.Mp` 表达的是一个具有如下类型的方法：

~~~go
func(float32) float32
~~~

和选择器的使用一样，使用一个指针引用了一个非接口类型的值接收器的方法，将自动将 `pt.Mv` 转化为 `(*pt).Mv`；和方法的调用一样，一个可取址的变量傻姑娘调用一个非接口类型的指针接收者方法，会自动将 `t.Map` 转化为 `(&t).Mp`

~~~go
f := t.Mv; f(7)   // like t.Mv(7)
f := pt.Mp; f(7)  // like pt.Mp(7)
f := pt.Mv; f(7)  // like (*pt).Mv(7)
f := t.Mp; f(7)   // like (&t).Mp(7)
f := makeT().Mp   // invalid: result of makeT() is not addressable
~~~

虽然上述例子都使用的是非接口类型，但同样可使用一个接口类型值创建方法值：

~~~go
var i interface { M(int) } = myVal
f := i.M; f(7) // like i.M(7)
~~~

# 9 Index expressions

`a[x]` 表达式表达的含义是：数组、指向数组指针、slice、string 的元素，或者是 map 类型的变量 a 使用 x 作为索引的 value 值。在这个表达式中，x 值被称之为位置索引或者 map 的 key 值。遵循如下规则：

如果 a 不是 map 类型的值，则有如下规则：

* x 值必须是整型，或者是无类型的常量值；
* 一个常量的索引值，必须是非负的而且可被表示为 int 类型的值；
* 一个常量的索引值，其类型是无类型的整型值；
* x 值的范围是 `[0, len(a))`

如果 a 是数组类型 A：

* 一个常量的索引值必须是在范围内的；
* 如果 x 超出了范围，会在运行时抛出运行时 panic；
* `a[x]` 是一个第 x 序位的元素，而且其类型是 A 的元素类型。

如果 a 是指向数组类型的指针：

* `a[x]` 就是 `(*a)[x]` 的简化表示方式。

~~~go
package main

import "fmt"

func main() {
	value := [...]int{1, 2, 3, 4}
	ptr := &value
	fmt.Println(ptr[0])
}
~~~

如果 a 是切片类型 s：

* 如果 x 在运行时超出了范围，会抛出运行时 panic；
* `a[x]` 是一个第 x 序位的元素，而且其类型是 s 的元素类型。

如果 a 是 string 类型：

* 如果 a 是常量，那么常量的位置索引值必须是在范围内的；
* 如果 x 在运行时超出范围，则会抛出运行时 panic；
* `a[x]` 是一个非常量的 byte 类型值，其类型是 byte；
* `a[x]` **不可被赋值**。

如果 a 是 map 类型 M：

* x 的类型必须是可被赋值给 M 的 key 类型的；
* 如果 map 包含一个 key 值是 x 的条目，那么 a[x] 的类型就是对应的 value 类型；
* 如果 map 是 nil，或者不包含任何条目，a[x] 就会是零值。

对于其他形式，`a[x]` 就是非法的。

在赋值表达式或初始化表达式中，map 值 a 的特殊形式：

~~~go
v, ok = a[x]
v, ok := a[x]
var v, ok = a[x]
~~~

可以得到一个另外的无类型 boolean 类型值。当 map 中存在该 key 值时，ok 的值为 true，反之则为 false。

# 10 Slice expressions

切片表达式可以从 string、array、pointer to array、slice 类型的值中构造出一个子字符串或者切片值。有 2 种变量形式：

1. **简单形式**：指明高低边界；
2. **完整形式**：除了指明高低边界，还给出了容量值。

## 10.1 Simple slice expression

对于 string、array、pointer to array、slice 类型的值，

~~~go
a[low : high]
~~~

表达式构造出了**一个子字符串或切片**。

~~~go
package main

import "fmt"

func main() {
	v := "michoi"
	fmt.Printf("%T.\n", v[1:]) // string.

	vArray := [...]int{1, 2, 3}
	fmt.Printf("%T.\n", vArray[1:2]) // []int.

	vSlice := []int{1, 2, 3, 4}
	fmt.Printf("%T.\n", vSlice[1:2]) // []int.
}
~~~

low 和 high 索引值表示的是结果值中元素的范围。其结果的元素索引从 0 开始，其**长度**为 high - low。

~~~go
a := [5]int{1, 2, 3, 4, 5}
s := a[1:4]
~~~

切片 s 的类型是 `[]int`，其长度是 3，容量是 4，对应的元素值分别是：

~~~go
s[0] == 2
s[1] == 3
s[2] == 4
~~~

很方便的是，`a[low : high]` 中的位置索引都可以省略。如果省略的是 low 值，那么**默认值**是 0；如果省略的是 high 值，那么**默认值**就是 a 的长度：

~~~go
a[2:]  // same as a[2 : len(a)]
a[:3]  // same as a[0 : 3]
a[:]   // same as a[0 : len(a)]
~~~

如果 a 是一个指向数组的指针，那么 `a[low : high]` 就是 `(*a)[low : high]` 的简化形式。

对于数组或字符串，位置索引的**范围和条件**是：`0 <= low <= high <= len(a)`；对于 slice 来说，**索引的上限**是切片容量 `cap(a)` 而不是长度。

如果是对一个数组类型值做切片操作，得到的将会是一个切片，而且 slice 元素类型和数组元素类型一致。

如果对一个值是 nil 的切片值做切片操作，得到的仍然是 nil 值。

~~~go
package main

import "fmt"

func main() {
	var a []int
	if a == nil {
		fmt.Println("a is nil!")
	}

	slice := a[:]
	fmt.Println(slice)
}
~~~

否则，其结果将会是一个切片，而且两者具备有相同的底层数组：

~~~go
var a [10]int
s1 := a[3:7]   // underlying array of s1 is array a; &s1[2] == &a[5]
s2 := s1[1:4]  // underlying array of s2 is underlying array of s1 which is array a; &s2[1] == &a[5]
s2[1] = 42     // s2[1] == s1[2] == a[5] == 42; they all refer to the same underlying array element
~~~

## 10.2 Full slice expression

对于 array、pointer to array 或者切片值 a 来说，

~~~go
a[low : high : max]
~~~

表达的含义是：构造了一个切片值。该切片值的长度和 `a[low:high]` 相同，而且有相同的元素类型。

此外，其 cap 值等于 `max - low`，仅仅第一个值 low 是可以省略的，默认值是 0.

~~~go
a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
~~~

得到的切片值 t 的类型是 `[]int`，长度是 2，容量是 4，各个元素值：

~~~go
t[0] == 2
t[1] == 3
~~~

如果 a 是指向数组的指针，那么 `a[low : high : max]` 是 `(*a)[low : high : max]` 的简化形式。如果是对一个数组做切片操作，那么必须是可被寻址的。

其中 low/high/max 的范围和大小关系：`0 <= low <= high <= max <= cap(a)`，否则是超出了范围。

# 11 Type assertions

Type assertion 意为：**类型断言**

对于一个关于接口类型类型 x 的表达式，以及类型 T 来说，表达式：

~~~go
x.(T)
~~~

其含义是，断言 x 的值不是 nil，而且接口类型 x 值中存储的结果类型是 T。

注解 `x.(T)` 被称之为**类型断言**。特别注意的是，上述类型断言表达式中，x 必须是接口！

比如示例程序：

~~~go
package main

import "fmt"

func main() {
	var a *aStruct
	if a == nil {
		fmt.Println("a is nil!")
	}

	a = &aStruct{
		age: 18,
	}
    var x interface{}
    x = a
	result, ok := x.(A)
	fmt.Println(ok, result)
}

type A interface {
	getAge() int
}

type aStruct struct {
	age int
}

func (a *aStruct) getAge() int {
	return a.age
}
~~~

更加精确的，如果 T 不是一个接口类型，`x.(T)` 则断言 x 的动态类型就是类型 T；在这种情况中，T 必须实现 x 的接口类型，否则这个类型断言是不合法的，因为 x 是不可能保存了 T 类型的值。如果 T 是一个接口类型，`x.(T)` 断言 x 的动态类型实现了 T 接口。

~~~go
func main() {
	a := aStruct{
		age: 18,
	}
	var x interface{}
	x = a
	result, ok := x.(aStruct)
	fmt.Println(ok, result) // true {18}
}
~~~

如果类型断言是成立的，那么表达式的值就是 x 中存储的值，而且值的类型是 T；反之，将会抛出运行时 panic。换句话说，即便 x 的动态仅在运行时才知道，但只有在正确的程序中 `x.(T)` 的类型才是 T。

~~~go
func main() {
	a := aStruct{
		age: 18,
	}
	var x interface{}
	x = a
	result := x.(int)
	fmt.Println(result)
}
~~~

使用上述的 `result := x.(int)` 类型断言，会在此处抛出 panic：`panic: interface conversion: interface {} is main.aStruct, not int`

~~~go
var x interface{} = 7          // x has dynamic type int and value 7
i := x.(int)                   // i has type int and value 7

type I interface { m() }

func f(y I) {
	s := y.(string)        // illegal: string does not implement I (missing method m)
	r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
	…
}
~~~

比如示例程序：

~~~go
package main

import (
	"fmt"
	"io"
)

func main() {
	a := &aStruct{
		age: 18,
	}
	f(a) // a 变量值实现了 A 接口类型
}

func f(a A) {
	s := a.(io.Reader) // 此时接口变量 a 的动态类型是 *aStruct 类型，而且也实现了 io.Reader 接口
	fmt.Println(s)
}

type A interface {
	getAge() int
	Read(p []byte) (n int, err error)
}

type aStruct struct {
	age int
}

func (a *aStruct) getAge() int {
	return a.age
}

func (a *aStruct) Read(p []byte) (n int, err error) {
	return 0, nil
}
~~~

在赋值表达式或者初始化表达式中使用类型声明：

~~~go
v, ok = x.(T)
v, ok := x.(T)
var v, ok = x.(T)
var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
~~~

得到了一个额外的无类型的 boolean 值。这种类型断言的写法中，如果类型断言成功，则 ok 的值是 true；反之，ok 的值是 false，且 v 的值是类型 T 的零值。这种情况下，是不会抛出运行时 panic 的。

# 12 Calls

给定函数类型 F 的表达式 f：

~~~go
f(a1, a2, … an)
~~~

表示的是使用参数 a1，a2，...，an 调用函数。除一种特殊情况外，参数必须是可分配给 F 参数类型的单值表达式，并在调用函数之前对其求值。上述函数调用表达式的结果类型就是 F 的结果类型。方法的调用和函数的调用类似，但方法会根据其接收者类型指定为一个选择器。

~~~go
math.Atan2(x, y)  // function call
var pt *Point
pt.Scale(3.5)     // method call with receiver pt
~~~

在函数调用中，函数的值和其参数会按照常规的顺序求值。在参数被求值后，这些参数值会被传给被调函数，此时被调函数开始执行。函数的返回值会依照返回参数返回给调用方。

调用 nil 值的函数值，会引发 panic 。

一种特殊情况，函数或方法 g 的返回参数的个数和另外的函数或方法的参数个数相等，此时函数或方法调用 `f(g(parameters_of_g))` 会在计算了 g 后触发调用 f。在这种情况中，f 的调用除 g 外不得包含任何其他参数，并且 g 至少应该具有一个返回值。如果 f 的参数列表末尾有 `...` 的可变参数，这些会使用 g 的返回值赋值。

~~~go
func Split(s string, pos int) (string, string) {
	return s[0:pos], s[pos:]
}

func Join(s, t string) string {
	return s + t
}

if Join(Split(value, len(value)/2)) != value {
	log.Panic("test fails")
}
~~~

如果类型 x 的方法列表包含 m，而且参数列表也能够传递给 m，那么 `x.m()` 是合法的；如果 x 是可寻址的，而且 `&x` 的方法集合中包含 m，那么 `(&x).m()` 的简化形式就可写成：`x.m()`。

~~~go
var p Point
p.Scale(3.5)
~~~

# 13 Passing arguments to ... parameters

如果函数 f 是一个可变参的函数，其最后一个参数 p 的类型是 `...T`。此时，参数 p 的类型相当于就是 `[]T`。如果在调用函数 f 时，没有传递任何实际的参数（也就是参数列表为空）时，相当于是传递了 nil 给参数 p；否则，传递给参数 p 的就会是一个新创建的 slice，slice 的类型就是 `[]T`。此时，slice 底层的数组元素依次就是调用函数 f 的参数，而且这些参数都必须是能够复制给类型 T 的。

可变参数类型对应的 slice 的长度、容量随每一次调用中实际传递的参数而变化。

比如，有这样的函数和调用：

~~~go
func Greeting(prefix string, who ...string)
Greeting("nobody")
Greeting("hello:", "Joe", "Anna", "Eileen")
~~~

在上述 `Greeing` 函数中，在第一次调用中，可变参数的值是 nil；第二次的值则是 `[]string{"Joe", "Anna", "Eileen"}`。

如果传递给最后的可变参数的是一个 slice（其类型是 `[]T`），那么如果在这个参数其后跟上 `...` 的话，就**不会新建任何 slice 变量而是直接作为函数的参数**。

~~~go
package main

import "fmt"

func main() {
	value := make([]int, 2)
	fmt.Printf("%p.\n", value)
	test(value...)

	fmt.Println(value)
}

func test(p ...int) {
	fmt.Printf("%p.\n", p)
	if p == nil {
		fmt.Println("nil")
	}
	p = append(p, 1)
	fmt.Printf("%p.\n", p)
}
~~~

在上面这个例子中，传递给 test 函数的实际上就是 `[]int` 类型的切片值，变量 p 和实际的参数“享有”相同的底层数组数据。

# 14 Operators

操作符用于在**表达式**中连接操作数。

~~~go
Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .

binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" | "|" | "^" .
mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .

unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
~~~

其中 `binary_op` 是指二元操作符；`unary_op` 是指一元操作符。

对于二元操作符来说，操作数的类型必须是相同的，除非这个操作中包含了移位或者无类型的常量值。对于仅包含有常量的操作，可以看看**常量表达式**的部分。

除了移位操作，如果其中一个操作数是无类型的常量，而另外一个操作数不是的，那么此时这个常量会隐式转化为另一个操作数的类型。

在移位表达式中，右侧操作数必须是整型的或者是一个可表示为 uint 类型的无类型常量值。在一个非整数移位的表达式中，如果左侧操作数是一个无类型常量，那么这个左侧操作数会隐式转化为其声明的类型：

~~~go
var a [1024]byte
var s uint = 33

// The results of the following examples are given for 64-bit ints.
var i = 1<<s                   // 1 has type int
var j int32 = 1<<s             // 1 has type int32; j == 0
var k = uint64(1<<s)           // 1 has type uint64; k == 1<<33
var m int = 1.0<<s             // 1.0 has type int; m == 1<<33
var n = 1.0<<s == j            // 1.0 has type int; n == true
var o = 1<<s == 2<<s           // 1 and 2 have type int; o == false
var p = 1<<s == 1<<33          // 1 has type int; p == true
var u = 1.0<<s                 // illegal: 1.0 has type float64, cannot shift
var u1 = 1.0<<s != 0           // illegal: 1.0 has type float64, cannot shift
var u2 = 1<<s != 1.0           // illegal: 1 has type float64, cannot shift
var v float32 = 1<<s           // illegal: 1 has type float32, cannot shift
var w int64 = 1.0<<33          // 1.0<<33 is a constant shift expression; w == 1<<33
var x = a[1.0<<s]              // panics: 1.0 has type int, but 1<<33 overflows array bounds
var b = make([]byte, 1.0<<s)   // 1.0 has type int; len(b) == 1<<33

// The results of the following examples are given for 32-bit ints,
// which means the shifts will overflow.
var mm int = 1.0<<s            // 1.0 has type int; mm == 0
var oo = 1<<s == 2<<s          // 1 and 2 have type int; oo == true
var pp = 1<<s == 1<<33         // illegal: 1 has type int, but 1<<33 overflows int
var xx = a[1.0<<s]             // 1.0 has type int; xx == a[0]
var bb = make([]byte, 1.0<<s)  // 1.0 has type int; len(bb) == 0
~~~

（**操作符优先级**）一元操作符有最高的优先级。因为 ++ 和 -- 操作符能构成语句（statements），而不是表达式（expressions），这两个操作符是脱离了操作符的结构。因此，在类似这样的语句中 `*p++` 实际上是等价于 `(*p)++`。

在 Go 语言中，有 5 个关于二元操作符的优先级别。而对于更加复杂的操作符，由更多的操作符组成，比如比较操作符、逻辑与、逻辑或操作符：

~~~go
Precedence    Operator
    5             *  /  %  <<  >>  &  &^
    4             +  -  |  ^
    3             ==  !=  <  <=  >  >=
    2             &&
    1             ||
~~~

&& 是 Go 语言中的逻辑与操作符，|| 则是逻辑或操作符。上述**等级 5 的优先级是最高的！**

表达式中包含了相同的二元操作符，操作数从左到右的结合。比如，`x / y * z` 等价于 `(x/y) * z`：

~~~go
+x
23 + 3*x[i]
x <= f()
^a >> b
f() || g()
x == y+1 && <-chanInt > 0
~~~

# 15 Arithmetic operators

算术操作符会在数值计算中使用，并阐述产生了一个和第一个操作数相同类型的值。有 4 种标准的算术操作可以用在整型、浮点型和复述型值中，它们是：+、-、*、/；在 string 类型值中，还可以使用 +。二进制操作符和移位操作符只能用在整型值中：

~~~go
+    sum                    integers, floats, complex values, strings
-    difference             integers, floats, complex values
*    product                integers, floats, complex values
/    quotient               integers, floats, complex values
%    remainder              integers

&    bitwise AND            integers
|    bitwise OR             integers
^    bitwise XOR            integers
&^   bit clear (AND NOT)    integers

<<   left shift             integer << unsigned integer
>>   right shift            integer >> unsigned integer
~~~

在整型操作符中，对于两个整型数值：x 和 y，`q = x / y` **取整操作**和 `r = x % y` **取余数操作**有如下的关系：

~~~go
x = q*y + r  and  |r| < |y|
~~~

关于取整操作，有这样的特性：

~~~go
 x     y     x / y     x % y
 5     3       1         2
-5     3      -1        -2
 5    -3      -1         2
-5    -3       1        -2
~~~

其中 `x % y` 的值可以用这样的公式计算得到：`x - q * y = r` 其中 r 值就是计算得到的余数。

对于上述计算规律有一个例外：如果被除数 x 是 int 类型的最小负数值，那么 `q = x / -1` 的结果等于 x（ r 值等于 0），因为整数溢出：

~~~go
			 x, q
int8                     -128
int16                  -32768
int32             -2147483648
int64    -9223372036854775808
~~~

可以用下述代码验证：

~~~go
package main

import "fmt"

func main() {
	var value int8 = -128
	fmt.Println(value / -1) // -128
}
~~~

如果除数是一个常量值，该值一定不能是 0。如果在程序的运行期中遇到除数是 0 的情况，会导致 run-time panic。如果被除数是非负数，除数是二进制的幂数常量值（2^n^），此时可以等价于是向右移位操作：

~~~go
x     x / 4     x % 4     x >> 2     x & 3
 11      2         3         2          3
-11     -2        -3        -3          1
~~~

对于移位操作来说，移位的计数值（count）必须是一个非负数；如果在 run-time 期间发现移位的计数值是负数，会导致 run-time panic。在移位操作中，被移位的数值如果是一个**无符号的数值（unsigned）**，会执行逻辑移位；如果是一个**有符号的数值（signed）**，会执行算数移位。对于移位的计数值（count），是没有上限的。

对于整型操作数，一元操作符 +、- 和 ^ 操作对应的定义是：

~~~go
+x                          is 0 + x
-x    negation              is 0 - x
^x    bitwise complement    is m ^ x  with m = "all bits set to 1" for unsigned x
                                      and  m = -1 for signed x
~~~

比如对 `int8` 类型的值 5，做上述运算：

~~~go
package main

import "fmt"

func main() {
	var value int8 = 5  // 1001
	fmt.Println(^value) // 1111 1111 ^ 0000 1001 --> 1111 0110 --> 0000 1001 + 1 --> -6
}
~~~

**整数值溢出**——**对于无符号的整型值来说**，操作符 +、-、* 和 << 会在 2^n^ 数值上做相应的操作，其中 n 值是无符号整型值的位宽。非严格的来说，对于这个无符号整型数值的操作中，会忽略高位字节的溢出。**对于有符号的整型数值来说**，操作符 +、-、* 和 << 会考虑是否有溢出，其最终计算的结果都只会在该有符号类型表示的数值范围之内。数值溢出是不会导致 run-time panic 的。编译器不会对这种计算做优化，因为编译器预期的计算是不会导致数值溢出的。

**浮点操作符**——浮点数操作会涉及到一个概念：`fused multiply and add`（FMA），对于 Go 来说，有如下可能性的实现：

~~~go
// FMA allowed for computing r, because x*y is not explicitly rounded:
r  = x*y + z
r  = z;   r += x*y
t  = x*y; r = t + z
*p = x*y; r = *p + z
r  = x*y + float64(z)

// FMA disallowed for computing r, because it would omit rounding of x*y:
r  = float64(x*y) + z
r  = z; r += float64(x*y)
t  = float64(x*y); r = t + z
~~~

**字符串拼接操作**——字符串的拼接操作，可以使用 + 或者 `+=` 操作符：

~~~go
s := "hi" + string(c)
s += " and good bye"
~~~

字符串的拼接操作，会产生一个新的字符串！

# 16 Comparison operators

比较操作符会比较 2 个操作数，其结果是一个无类型的 boolean 值：

~~~go
==    equal
!=    not equal
<     less
<=    less or equal
>     greater
>=    greater or equal
~~~

在上述操作符的任何写法中，**第一个操作数必须是能够赋值给第二个操作数的类型的**，反之亦然。否则会被视为非法的表达式。比如：`value := 1 > true` 编译器会报错：`cannot convert 1 (untyped int constant) to untyped bool`

（**有 2 种比较规则**）在使用判断相等性的操作符：`==` 或 `!=` 时，操作数必须是可以**判断相等性**的；在使用判断大小的操作符：`<`、`<=`、`>`、`>=` 时，操作数必须是可以**判断大小**的。与之相关的规则如下：

* 布尔类型值只能比较相等性；
* 整型值既是可比较相等性，也能比较大小；
* 浮点值既是可比较相等性，其比较规则依据 IEEE-754 标准；
* 复数值只能比较相等性，如果两个复数的实部和虚部都相等，就代表两个复数相等；
* 字符串类型值既是可比较相等性，也能比较大小，按字节顺序进行比较；
* 指针类型值只能比较相等性，如果两个指针指向相同的变量，或者都是 nil 则两个指针值相等；
* channel 类型值只能比较相等性：其底部数据相同（channel 是引用类型值），或者值都是 nil，则认为是相等的；
* 接口类型值只能比较相等性：如果两个接口值的动态类型和动态值都相等，或者都是 nil，则认为是相等的；
* 非接口类型 X 的值 x，和接口类型 T 的值 t，两者是可比较的，条件是 X 类型的值是可比较的，而且 X 实现了 T 接口。如果 t 的动态类型是 X，其动态值是 x，那么认为 t 和 x 是相等的。
* struct 类型值是可比较的，其条件在于：struct 的各个域都是可比较的。
* array 类型值是可比较的，其条件在于：array 的元素是可比较的。如果 array 的各个元素值都相等，那么认为 array 是相等的。





# 17 Logical operators

逻辑操作符：

~~~go
&&    conditional AND    p && q  is  "if p then q else false"
||    conditional OR     p || q  is  "if p then true else q"
!     NOT                !p      is  "not p"
~~~























# 18 Address operators





# 19 Receive operator





# 20 Conversions





# 21 Constant expressions



# 22 Order of evaluation





解答下述疑惑：

~~~go
func (gen *myGenerator) Start() bool {
	logger.Infoln("Starting load generator...")

	if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}
    ...
}
~~~

`&gen.status` 优先级是怎样的？其中 myGenerator 是一个结构体。