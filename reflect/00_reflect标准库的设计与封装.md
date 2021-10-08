Go 的 reflect 标准库实现的是**运行时反射**机制，允许应用程序**操作（修改）**任意类型实例。

> 运行时的反射机制，关键词：
>
> 1. 运行时：runtime
> 2. 反射机制：“反”这个字，由运行时的实例（对象）获取到这个对象的类型、值信息。

reflect 包的典型使用方式：使用**静态类型 `interface{}` 装载**一个值，调用 reflect.TypeOf 函数获取其**动态类型信息**，其返回值是**一个 Type 实例**。调用 reflect.ValueOf 函数可以返回**一个 Value 实例**，表示的是其**动态值信息**。如果返回的是零值，那么表示对应类型的零值。

reflect 机制的通常用法是**和 `interface{}` 类相关的**，能获取其**动态 Type 和动态 Value**。

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

为什么将 Type 设计成 interface？

# 1 The Laws of Reflection

> https://go.dev/blog/laws-of-reflection

在计算科学中，反射是计算机程序的一种的能力：特别是通过**类型**，检查**自身结构**。这种能力是元编程 metaprogramming 的组成部分。当然，反射的很多内容让人很困惑。

在这篇文章中，我试图通过**解释 Go 中反射的运行机制和原理**，以此理清这些让人困惑的内容。每一种编程语言的**反射模型**都是不同，当然还有一些编程语言是不支持反射的，但这篇文章的场景是 Go 语言，因此，接下来的内容中 reflection 的含义就是 Go 语言中的反射。

## 1.1 类型和接口 Types and interfaces

**反射机制的基础**是**类型系统**，因此，我们从 Go 语言中的类型开始。

Go 语言是一种**静态类型**的编程语言。每个变量都有一个静态类型，这个类型在编译期被确定（已知、固定），比如 int、float32、*MyType、[]byte 等等。如果我们声明如下：

~~~go
type MyInt int

var i int
var j MyInt
~~~

这样，变量 i 的类型是 int；变量 j 的类型是 MyInt。变量 i 和 j 有**不相同的静态类型**，虽然它们具有**相同的底层类型**，但如果没有显式地转换，相互之间是不能被赋值的。

关于类型的一个重要类别是 interface 的类型，其中 interface 表示的是一系列固定方法的集合。一个接口变量可以用于保存任何固定值（非接口值），但这个值必须实现了 interface 的所有方法。一个很显然的实例是来自 io 包中的接口 io.Reader 和 io.Writer：

~~~go
// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}
~~~

任何一个实现了这种签名的 Read 方法或 Write 方法都称之为是实现了 io.Reader 或 io.Writer 接口。讨论这部分的目的是，**一个 io.Reader 类型的变量**能够装载任何实现了 Read 方法的类型对应的值：

~~~go
var r io.Reader
r = os.Stdin
r = bufio.NewReader(r)
r = new(bytes.Bufer)
// and so on
~~~

有一点必须要明确的是：不管变量 r 中装载的内容是什么，比如上面的 os.Stdin、bufio.NewReader(r) 等等，变量 r 的类型始终是 io.Reader。**Go 语言是一种静态类型语言，变量 r 的静态类型是 io.Reader**。

另外一个极度重要的关于 interface 类型的例子是：**空 interface**

~~~go
interface{}
~~~

`interface{}` 在 Go 中是**一种类型**，它代表的是**一个空的方法集合**，Go 语言中的任何值都满足条件（因为任何值都存在 0 个或更多数目的方法）。也就是说，`interface{}` 类型的变量，可以装载任何 Go 语言中的值。

一些开发人员说：Go 中的接口是**动态类型**的，这实际上是**误导**！它们仍然是**静态类型**的：一个接口类型的变量始终具有相同的静态类型（比如上面说的 io.Reader 类型），虽然在 Runtime 时期保存在接口变量中的值可能改变其类型，但是这个值仍然是满足于接口类型的（值的真实类型是实现了接口方法）。

我们需要更加详细的理解上面的内容，因为**反射**和**接口**是非常接近的。

## 1.2 一个接口的表示

一个接口类型的变量存储了**一对值**：赋值给变量的固定（明确）**值**，以及关于值的**类型描述符**。也就是：`(value - type descriptor)`。更加详细的：value 是实现了这个接口的类型值，type 描述符是关于类型的完整描述。比如，下述示例程序：

~~~go
var r io.Reader
tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil {
    return nil, err
}
r = tty
~~~

程序在运行后，变量 r 包含了 `(value - type descriptor)` —— `(tty, *os.File)`。需要注意的是 `*os.File` 类型实现的方法远多于 io.Reader 接口中包含的方法，虽然这个接口（io.Reader）仅提供了执行 Read 方法的途径，但是 value 中保存了所有关于这个 `*os.File` 类型的信息。因此，我们可以这样使用：

~~~go
var w io.Writer
w = r.(io.Writer)
~~~

上面表达式是**类型断言**表达式，断言：r 内部包含的 item 也实现了 io.Writer 接口，因此，我们能够将 value 赋值给 w。经过这个赋值表达式后，变量 w 中也包含了 `(value - type descriptor)` —— `(tty, *os.File)`，这其中的 pair 是和 r 相同的内容。**接口的静态类型决定了接口变量可以调用的方法**，虽然其内部的真实值（concrete value）包含有更多的方法。

继续，我们还可以这样做：

~~~go
var empty interface{}
empty = w
~~~

这个 `interface{}` 类型的值中也包含了相同的  `(value - type descriptor)` —— `(tty, *os.File)`。**很方便**，一个空接口可以承载任何值，而且还包含了和这个值相关的任何信息。

上面的这个赋值是不需要类型断言的，因为 w 实现了 `interface{}` 接口。在之前的示例程序中，我们将值从一个 Reader 赋值给一个 Writer，需要一个显式的类型断言，因为 Writer 的方法并不是 Reader 的子集。

一个重要的细节是： 一个接口值中的`(value - type descriptor)` 具有严格的要求，必须是 `(value - concrete type)` 而不能是 `(value, interface type)`。也就是说，**接口不能用来承载接口值**。

## 1.3 反射规则：从 interface 值到 reflect 实体值

在“新手村”——最开始认识 reflection 时，我们可以这样认为：反射仅仅就是一种**检查 interface 变量**中  `(value - type descriptor)` 的机制。最开始，我们需要知道在 reflect 标准库中的 2 种类型：Type 和 Value。这两个类型提供了访问一个接口变量中的  `(value - type descriptor)` 的方式，以及 2 个简单的函数：`reflect.TypeOf` 和 `reflect.ValueOf`，对应的能够解析出一个 interface 变量中的 reflect.Type 和 reflect.Value 值。（当然了，通过 reflect.Value 也是有办法获取到对应的 reflect.Type 的，但是在这里，保持 Value 和 Type 概念的独立性）

我们从 TypeOf 开始：

~~~go
package main

import (
	"fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    fmt.Println("type:", reflect.TypeOf(x))
}
type: float64
~~~

你可能会疑惑：哪个有一个 interface 接口值，因为程序看起来仅仅只是传递了 x 变量值，并不是一个 interface 值。但是在 reflect 的 doc 中，可以看到函数声明：

~~~go
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
~~~

在 TypeOf 的参数中，确实是一个 interface{} 类型的变量等着接收入参值。当我们调用 `reflect.TypeOf(x)` 时，x 值首先会作为入参被存储到一个空的接口中，紧接着 `reflect.TypeOf` 会 **unpack 这个空接口值**，并**从中恢复出类型信息**。

同样的，对于 `reflect.ValueOf` 函数来说，会从这个空接口值中恢复出值：

~~~go
fmt.Println("ValueOf:", reflect.ValueOf(x).String()) // ValueOf: <float64 Value>
~~~

此处，我们显式调用了 String() 方法，因为默认情况下 fmt 包会深入到 reflect.Value 调用该类型的 String() 方法，输出接口值中的类型信息，比如上述的 `<float64 Value>`。

reflect.Type 和 reflect.Value 都有**很多**能够**检查和操作它们的方法**。一个很重要的例子是：Value 有一个名为 Type 的方法，能够返回 reflect.Value 的 Type 值。

~~~go
// Type returns v's type.
func (v Value) Type() Type
~~~

另外一个是 Type 和 Value 都有 Kind 的方法，用于获取一个表示 item 类型的值，比如 Uint、Float64、Slice 等等。reflect.Value 类型上的名为 Int 和 Float 的方法能够获取到 item 的值（以 int64 和 float64 类型）：

~~~go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())
fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
fmt.Println("value:", v.Float())
~~~

reflect 标准库中有一些属性值得挑出来：

1. Value 的 `getter` 和 `setter` 方法是用于操作最长类型，比如 int64 是用于所有有符号整型值。也就是说，`Int` 方法能够返回一个 `int64`，`SetInt` 方法能够设置一个 int64 的值。因此，有些情况下可能需要做类型转换。

   ~~~go
   var x uint8 = 'x'
   v := reflect.ValueOf(x)
   fmt.Println("type:", v.Type())
   fmt.Println("kind is uint8:", v.Kind() == reflect.Uint8)
   x = uint8(v.Uint()) // v.Uint returns a uint64
   ~~~

2. Kind 值表示的是 `(value - type descriptor)`中**底层类型信息**，而**不是静态类型**。如果一个反射实例包含的是一个用户自定义的整型类型，比如：

   ~~~go
   type MyInt int
   var x MyInt = 7
   v := reflect.ValueOf(x)
   ~~~

   Kind 方法返回的依然是 reflect.Int，即使 x 的静态类型是 MyInt，而不是 int。换句话说，**Kind 获取的是底层类型，而 Type 会对其做区分**。

   ~~~go
   func main() {
   	type MyInt int
   
   	var x MyInt = 3
   	value := reflect.ValueOf(x)
   	fmt.Println(value.Type(), value.Kind() == reflect.Int)
   }
   Main.MyInt true
   ~~~

   因为在 Go 语言中我们可以使用 type 关键字构造很多自定义类型，而种类 Kind 就是指**底层的类型**。特别的，可以获取到的种类有：**Array、Chan、Func、Map、Ptr、Slice、String、Struct** 等。

## 1.4 反射：从 reflect 实体值到 interface 值

类似于物理中的反射，Go 语言中的反射会存在一个**逆向**过程。

给定一个 reflect.Value 值，我们可以使用 Interface 方法获取到一个 interface 值。也就是意味着：这个方法能够从 reflect.Value 中获取 `<value - type descriptor>` 值，并重新组装成一个 interface 值：

~~~go
// Interface returns v's current value as an interface{}.
// It is equivalent to:
//	var i interface{} = (v's underlying value)
// It panics if the Value was obtained by accessing
// unexported struct fields.
func (v Value) Interface() (i interface{}) {
	return valueInterface(v, true)
}
~~~

作为示例：

~~~go
func main() {
	type MyInt int

	var x MyInt = 3
	value := reflect.ValueOf(x)

	tmp := value.Interface() // interface{}
	fmt.Println(tmp)

	v, ok := tmp.(MyInt) // true, is not int
	if ok {
		fmt.Println(v)
	}
}
~~~

换句话说，Interface 方法相当于是 ValueOf 函数的逆过程，不同之处在于，Interface 方法的返回值的静态类型始终是 `interface{}`。

## 1.5 修改 reflect 值，value 必须是可被修改的

下面是一个不能运行的代码，但很值得学习：

~~~go
var x float64 = 3.4
v := reflect.ValueOf(x)
v.SetFloat(7.1) // Error: will panic
~~~

我们得到的 panic 信息是：

~~~
panic: reflect: reflect.Value.SetFloat using unaddressable value
~~~

其原因并不是 7.1 不是一个可定位的值，而是变量 v 不是。**可被设置值**是 reflect.Value 的一个属性，并不是所有的 reflect.Value 都具有这个属性。

reflect.Value 的 CanSet 方法用于测试 reflect.Value 是否具有该属性：

~~~go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("settability of v:", v.CanSet())
~~~

因此，在一个不具有可设置属性的值上调用 Set 方法，显然是不合适的。那么什么是**可设置属性**？

**可设置属性**就像是可取地址类似，是一个比特位。可设置属性的含义就是能够使用 reflect.Value 值改变其底层实际存储的值。**可设置属性取决于 reflect.Value 是否保存其源数据变量**。比如：

~~~go
var x float64 = 3.4
v := reflect.ValueOf(x)
~~~

我们传递了 x 的一份拷贝到 reflect.ValueOf 函数中，因此，reflect.ValueOf 函数使用的是 x 的拷贝而不是 x 本身创建了 reflect.Value 值。那如果：

~~~go
v.SetFloat(7.1)
~~~

运行成功，也不会改变原先 x 变量的值。这种情况下，只会去改变现在保存在 reflect.Value 中的值——x 的一份拷贝，而 x 不会有任何改变。因此，这样做是**非法的**，可设置属性就是为了避免这种情况的发生。

这种场景就和传递 x 给一个函数是类似的：

~~~go
f(x)
~~~

我们并没有想过要让函数 f 去修改变量 x 的值，因为我们传递的是 x 的一份拷贝，而不是 x 本身。如果我们想要达到此目的，我们必须穿进去的是 x 的指针：

~~~go
f(&x)
~~~

相同的道理：

~~~go
var x float64 = 3.4
p := reflect.ValueOf(&x)
fmt.Println("type of p:", p.Type())
fmt.Println("settability of p:", p.CanSet())
~~~

上述示例程序输出结果是：

~~~shell
type of p: *float64
settability of p: false
~~~

反射实体 p 并不具有可设置属性，但是**我们并不需要修改 p 的值，实际上是要修改 `*p` 的内容**。为了获取到 p 指向的实体，我们需要调用 reflect.Value 的 Elem 方法：

~~~go
v := p.Elem()
fmt.Println("settability of v:", v.CanSet())
~~~

此时 v 就是一个具有可设置属性的反射实体（reflect.Value），并且此时代表的就是变量 x，我们可以通过 v 修改 x 的值。

总之，我们只需要记住：**如果想要修改源数据，就需要使用变量的指针（地址）**。

## 1.6 结构体

在上面示例程序中，v 并不是指针本身，但是是从指针衍生出来的 reflect.Value 值。一个更加通用的场景是使用反射机制**修改 struct 的各个字段值**。因为我们已经拿到了结构体变量的地址，就可以修改各个字段的值。

下面是一个简单的例子用于解析结构体类型值 t。我们使用结构体变量的指针创建反射值，是因为接下来要求修改结构体的字段值。

~~~go
type T struct {
    A int
    B string
}
t := T{23, "skidoo"}
s := reflect.ValueOf(&t).Elem()
typeOfT := s.Type()
for i := 0; i < s.NumField(); i++ {
    f := s.Field(i)
    fmt.Printf("%d: %s %s = %v \n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
}
~~~



