nil 值有类型吗？

~~~go
func main() {
	var x interface{}
	fmt.Printf("%T, %v.\n", x, x) // <nil>, <nil>.
}
~~~

上述这种情况输出的是变量 x 的动态类型，而其动态类型是 nil。接口类型是不能实例化的。

~~~go
// nil is a predeclared identifier representing the zero value for a
// pointer, channel, func, interface, map, or slice type.
var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int
~~~





**类型是什么？**

**类型决定了一组值以及在这些值上的一系列操作和方法**。一个类型可以使用一个类型名表示，或者是使用特定的类型字面值（该类型字面值由现有类型组成 compose）。

~~~go
Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent .
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType | SliceType | MapType | ChannelType .
~~~

Go 语言预声明（predeclare）了一些类型名称，其他的类型则使用的是**类型声明语句**（type declarations，包括 2 种形式：类型别名和类型重定义）。复合类型——array、struct、pointer、function、interface、slice、map 和 channel——可以使用**类型字面值构造**。

~~~go
type Person struct {
    name string
    age uint8
}
~~~

每一种类型 T 都有一个**底层类型**（underlying type）：如果类型 T 是预声明的布尔、数值、字符串类型，或者是类型字面值，对应的底层类型就是类型 T 自身。反之，类型 T 的底层类型是 T 在其类型声明语句引用的类型的底层类型。

~~~go
type (
	A1 = string // 类型别名
	A2 = A1
)

type (
	B1 string // 类型重定义
	B2 B1
	B3 []B1
	B4 B3
)
~~~

string、A1、A2、B1 的底层类型是 string；[]B1、B3、B4 的底层类型是 []B1。

下面依次就各个小主题（**Go 语言中的各种不同的类型**）讲述：

* bool
* numeric
* string
* array
* struct
* slice
* pointer
* function
* method
* interface
* map
* channel

# 1 bool

Go 语言中的布尔类型表示的是 Boolean 真值的集合，该集合包含有预声明的 true 和 false 两个值。

预声明的布尔类型使用 bool 表示，它是一个**类型重定义**的类型：

~~~go
package builtin

// bool is the set of boolean values, true and false.
type bool bool

// true and false are the two untyped boolean values.
const (
	true  = 0 == 0 // Untyped bool.
	false = 0 != 0 // Untyped bool.
)
~~~

在 App 开发中可以使用的就是 bool 类型，实际上述代码中 true 和 false 是**无类型的布尔值**！

~~~go
package main

import (
	"fmt"
)

func main() {
	printTypeAndValue(true)
}

func printTypeAndValue(value interface{}) {
	fmt.Printf("%T, %v.\n", value, value)

	if nil == value {
		fmt.Println("value is nil!")
	}
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
bool, true.
~~~

但是

1. **问题在于 true 和 false 是如何与 bool 类型关联的？**
2. bool 类型占用多少个字节？
3. `type bool bool` 右侧的 bool 又指代的是什么？

# 2 numeric

数值类型表示的整型、浮点型值的集合。预声明的**独立于架构体系**的数值类型包括：

~~~go
uint8       the set of all unsigned  8-bit integers (0 to 255)
uint16      the set of all unsigned 16-bit integers (0 to 65535)
uint32      the set of all unsigned 32-bit integers (0 to 4294967295)
uint64      the set of all unsigned 64-bit integers (0 to 18446744073709551615)

int8        the set of all signed  8-bit integers (-128 to 127)
int16       the set of all signed 16-bit integers (-32768 to 32767)
int32       the set of all signed 32-bit integers (-2147483648 to 2147483647)
int64       the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

float32     the set of all IEEE-754 32-bit floating-point numbers
float64     the set of all IEEE-754 64-bit floating-point numbers

complex64   the set of all complex numbers with float32 real and imaginary parts
complex128  the set of all complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32
~~~

一个 n-bit 的整型数值，具备有 n 个比特位宽，并使用**二进制补码**表示。

Go 语言中同样包含有**特定于实现（implementation-specific）宽度**的预定义数值类型，也就是和具体的平台相关：

~~~go
uint     either 32 or 64 bits
int      same size as uint
uintptr  an unsigned integer large enough to store the uninterpreted bits of a pointer value

// uint is an unsigned integer type that is at least 32 bits in size. It is a
// distinct type, however, and not an alias for, say, uint32.
type uint uint

// int is a signed integer type that is at least 32 bits in size. It is a
// distinct type, however, and not an alias for, say, int32.
type int int

// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
~~~

为了避免移植性（portability）问题，所有的数值类型都是已定义类型（defined types），除了 byte 和 rune 均应该被区分。byte 类型是 uint8 类型的别名类型，rune 类型是 int32 的别名类型。当表达式或者赋值语句中混合有不同的数值类型时，需要有显式的（explicit）类型转换（conversion）。例如，int32 和 int 类型是不相同的类型，虽然在特定的计算架构上它们有相同的位宽。

# 3 string

string 类型表示的是字符串值的集合。一个字符串值是一个（可能是空的）字节序列。**字节的数量被称之为字符串的长度**，是一个非负的值。预声明的字符串类型名是 string，是一个已定义的类型。

~~~go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
~~~

一个 string 值可能是空串，但是不能为 nil。此外，string 值一旦创建，就是不可变的（不可能去修改字符串的内容）。

~~~go
package main

import "fmt"

func main() {
	var value string
	value = "中国"
	fmt.Println(len(value)) // 6
}
~~~

可以使用内置函数 len 来获取字符串值 s 的长度。如果字符串值是常量，那么该长度在编译期就是常量值。一个字符串的字节值可以通过整数索引 0 ~ len(s) - 1 获取到。获取字符串中字节索引值的地址是非法的，比如 `&s[i]`。

~~~go
// The len built-in function returns the length of v, according to its type:
//	Array: the number of elements in v.
//	Pointer to array: the number of elements in *v (even if v is nil).
//	Slice, or map: the number of elements in v; if v is nil, len(v) is zero.
//	String: the number of bytes in v.
//	Channel: the number of elements queued (unread) in the channel buffer;
//	         if v is nil, len(v) is zero.
// For some arguments, such as a string literal or a simple array expression, the
// result can be a constant. See the Go language specification's "Length and
// capacity" section for details.
func len(v Type) int
~~~

# 4 array

一个数组就是一系列顺序的元素类型的值，元素的数量称之为数组的长度。

~~~go
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .
~~~

**数组的长度是数组类型的一部分**，它必须求值为可以由 int 类型的值表示的非负常数。可以使用内建函数 `len` 获取到数组 a 的长度，同样也可以使用整数位置索引 `0 ~ len(a) - 1` 获取到元素值。一般情况使用的数组都是一维数组，可以由此构建多维数组类型。

~~~go
[32]byte
[2*N] struct { x, y int32 } // 结构体数组，每个数组元素都是一个（匿名的）结构体
[1000]*float64 // 指针数组
[3][5]int
[2][2][2]float64  // same as [2]([2]([2]float64))
~~~

数组是值类型：

~~~go
package main

import "fmt"

func main() {
	var value = [3]int8{1, 2, 3}
	// 0xc000014080, 0xc000014081, 0xc000014082.
	fmt.Printf("%p, %p, %p.\n", &value[0], &value[1], &value[2])
}
~~~

不能够直接通过 value[0] 获得其地址（数组不是引用类型，而是值类型），而只能通过取地址操作符 & 得到。从上面可以看出，数组元素是 int8 也就是 1 个字节时，底层地址值增加 1（从 0xc000014080 --> 0xc000014081）。

# 5 struct

一个结构体是一系列的名为字段的命名元素，每一个字段都有一个**名字**和**类型**。字段的名字可以被明确指定**（标识符列表）**或者隐式指定（**嵌入字段**）。**在结构体中，非空白字段名称必须唯一**。

~~~go
StructType    = "struct" "{" { FieldDecl ";" } "}" .
FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
EmbeddedField = [ "*" ] TypeName .
Tag           = string_lit .
~~~

比如下面的结构体定义示例：

~~~go
// An empty struct.
struct {}

// A struct with 6 fields.
struct {
	x, y int
	u float32
	_ float32  // padding
	A *[]int
	F func()
}
~~~

**【结构体类型中嵌入类型字段】**使用类型而不显式给定字段名称，这种方式被称为**嵌入字段**。一个嵌入字段必须是使用类型名称 T 指定，或者一个指向非接口类型的指针 `*T` 类型指定（且此时的 T 类型不能是接口类型）。此时**这种类型名就是字段名**：

~~~go
// A struct with four embedded fields of types T1, *T2, P.T3 and *P.T4
struct {
	T1        // field name is T1
	*T2       // field name is T2
	P.T3      // field name is T3
	*P.T4     // field name is T4
	x, y int  // field names are x and y
}
~~~

下面的这些字段声明是非法的，因为在结构体类型中字段名必须唯一：

~~~go
struct {
	T     // conflicts with embedded field *T and *P.T
	*T    // conflicts with embedded field T and *P.T
	*P.T  // conflicts with embedded field T and *T
}
~~~

下面的这个知识点是：**struct 类型、方法、接口 3 个知识点的汇合点**！

* 如果 S 包含了 T 作为嵌入类型：S 和 `*S` 的方法集合将包含 T 类型值作为接收者的方法集合；`*S` 的方法集合包含了 `*T` 类型指针方法集合。
* 如果 S 包含的是 `*T` 作为嵌入类型：S 和 `*S` 的方法集合将包含 T 类型值方法和 `*T` 类型指针方法集合。

关于是否被嵌入类型是否“继承”了嵌入类型的方法，不能使用 `selector` 检验，因为有可能编译器自动做了转换——“语法糖”。而只能使用下面这种接口赋值的方式进行。关于内嵌字段对应的方法集合，可以做这样的示例程序：

~~~go
package main

import "fmt"

type AType struct {
	a int
	BType
}

type BType struct {
	x, y int
}

type ToolInterface interface {
	function()
}

func (b BType) function() {
	fmt.Println(b.y)
}

type ToolInterfacePtr interface {
	assertPtr()
}

func (b *BType) assertPtr() {
	fmt.Println(b.x, b.y)
}

func main() {
	b := BType{
		x: 1,
		y: 2,
	}
	a := AType{
		BType: b,
	}

	var obj ToolInterface
	obj = b
	obj = &a
	fmt.Println(obj)

	var assertPtr ToolInterfacePtr
	assertPtr = &b
	assertPtr = &a
	fmt.Println(assertPtr)
}
~~~

这个程序说明了如下结论：如果 S 包含了 T 作为嵌入类型

* S 的方法集合包含 T 类型值作为接收者的方法集合，但不包含 *T 类型值作为接收者的方法集合；
* `*S` 的方法集合包含了 T 类型和 `*T` 类型值作为接收者的方法集合。

~~~go
package main

import "fmt"

type AType struct {
	a int
	*BType
}

type BType struct {
	x, y int
}

type ToolInterface interface {
	function()
}

func (b BType) function() {
	fmt.Println(b.y)
}

type ToolInterfacePtr interface {
	assertPtr()
}

func (b *BType) assertPtr() {
	fmt.Println(b.x, b.y)
}

func main() {
	b := BType{
		x: 1,
		y: 2,
	}
	a := AType{
		BType: &b,
	}

	var obj ToolInterface
	obj = b
	obj = &a
	fmt.Println(obj)

	var assertPtr ToolInterfacePtr
	assertPtr = &b
	assertPtr = &a
	fmt.Println(assertPtr)
}
~~~

这个程序说明了如下结论：如果 S 包含了 `*T` 作为嵌入类型

* S 的方法集合包含了 T 类型和 `*T` 类型值作为接收者的方法集合；
* `*S` 的方法集合包含了 T 类型和 `*T` 类型值作为接收者的方法集合。

结构体中一个字段的声明后面可能会跟着一个**可选的字符串字面量标签**，将会作为其**对应的属性**存在。一个空的标签字符串，等价于一个空白的标签。这些标签通过**反射接口**是可见的，并参与结构体的类型标识，但在其他情况下将被忽略。

结构体中一个字段的声明后面可能会跟着一个可选的字符串字面量标签，将会作为其对应的属性存在。一个空的标签字符串，等价于一个空白的标签。这些标签通过反射接口是可见的，并参与结构体的类型标识，但在其他情况下将被忽略。

~~~go
struct {
	x, y float64 ""  // an empty tag string is like an absent tag
	name string  "any string is permitted as a tag"
	_    [4]byte "ceci n'est pas un champ de structure"
}

// A struct corresponding to a TimeStamp protocol buffer.
// The tag strings define the protocol buffer field numbers;
// they follow the convention outlined by the reflect package.
struct {
	microsec  uint64 `protobuf:"1"`
	serverIP6 uint64 `protobuf:"2"`
}
~~~

# 6 slice

切片是底层数组连续片段的描述符，并提供了对该数组中元素序列的访问方法。切片元素的数量称为切片的长度，并且永远不会为负数。一个未初始化的 slice，其值是 nil：

~~~go
SliceType = "[" "]" ElementType .
~~~

切片的长度使用内建函数 len 获取到。不像数组，切片的长度在程序运行过程中是可能改变的。和数组一样，都是可使用整数位置索引获取到其各个元素。一个 slice 的给定元素的索引值，有可能比其他有相同底层数组的切片元素索引值要小。

~~~go
package main

import "fmt"

func main() {
	var value = []int32{1, 2, 3}

	slice := value[1:3]
	fmt.Println(slice)              // [2 3]
	fmt.Println(value[1], slice[0]) // 2 2，获取到的是相同的元素值，但是位置索引不同
}
~~~

一个 slice，一旦初始化了，都会关联上底层用于保存元素的数组。这个 slice 会和底层的 array 享有共同的存储空间，同样其他的 slice 也是可以引用这个相同的底层数组的。比较而言（by contrast），不同的底层数组都表示有不同的底层存储区域。

切片底层的数组可能会延伸到切片的末尾，其容量则是这种延伸的度量。可以通过从原始切片中 slicing 一个长度达到该容量的新 slice。一个 slice 的容量，可以使用内建函数 cap() 获取到。

使用给定的类型 T 新建一个已初始化的 slice 值，可以使用内建的 make 函数。在使用这个函数时，包含了元素的类型，以及指定的 len 值，和备选的 cap 值。使用这种方式创建的 slice，都会分配一个新的、隐含的底层数组，这个数值由 slice 引用：

~~~go
make([]T, length, capacity)
~~~

分配一个数组，同时使用切片操作，可以产生相同的切片，下述表达式是等价的：

~~~go
package main

import "fmt"

func main() {
	slice1 := make([]int, 50, 100)
	printInfo(slice1)

	slice2 := new([100]int)[0:50]
	printInfo(slice2)
}

func printInfo(slice []int) {
	fmt.Printf("%p, %d, %d.\n", slice, len(slice), cap(slice))
}
~~~

和数组类似，一般使用 slice 都是一维的，也是可以创建多维的切片。

# 7 pointer

一个指针类型表示的一类指向给定类型 T 变量的指针集合，被称之为类型 T 的指针。一个未初始化的指针类型值是 nil：

~~~go
PointerType = "*" BaseType .
BaseType    = Type .
~~~

比如：

~~~go
*Point // 类型 Point 的指针
*[4]int // 数组类型 [4]int 的指针
~~~

指针类型的名称是一个整体：`*Point`，可以这样表述：**Point 的指针类型**。

# 8 function

函数类型意味着是一类具有**相同参数类型和返回值类型**的函数。函数类型也是一种引用类型，其未初始化的函数类型变量值是 nil。函数类型的组成类似：

~~~go
FunctionType   = "func" Signature .
Signature      = Parameters [ Result ] .
Result         = Parameters | Type .
Parameters     = "(" [ ParameterList [ "," ] ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
~~~

在函数类型的参数或结果列表中，标识符的名称**要么全部显式写明，要么全部不存在**。如果这些标识符是存在的，每个名称代表指定类型的一项参数或结果，并且签名中的所有非空名称都必须是唯一的。如果不存在，则每种类型代表该类型的一项。

参数和结构列表始终带有括号（parenthesized），但如果有一个确切的未命名结果，则可以将其写为非括号类型。比如下述的 `func(a, _ int, z float32) bool`

在函数签名中，参数列表的最后的参数可以带有一个 `...` **类型前缀符**，这种情况称之为**可变参数**（variadic）。对于可变参数，**既可以传入 0 个参数，也可以传入多个参数**：

~~~go
func()
func(x int) int
func(a, _ int, z float32) bool
func(a, b int, z float32) (bool)
func(prefix string, values ...int)
func(a, b int, z float64, opt ...interface{}) (success bool)
func(int, int, float64) (float64, *[]int)
func(n int) func(p *T)
~~~

# 9 method

一个类型可能有一系列与之关联的方法集合。一个接口类型的方法集合是它的接口。

* T 的方法集合由所有声明使用类型 T 作为接收者的方法组成；
* `*T` 的方法集合既包含了 *T 作为接收者的方法，也包含了 T 作为接收者的方法。

可以用下面这种方法进行验证：

~~~go
package main

import "fmt"

type BType struct {
	x, y int
}

type ToolInterface interface {
	function()
}

func (b BType) function() {
	fmt.Println(b.y)
}

type ToolInterfacePtr interface {
	assertPtr()
}

func (b *BType) assertPtr() {
	fmt.Println(b.x, b.y)
}

func main() {
	b := BType{
		x: 1,
		y: 2,
	}
	ptr := &b

	var obj ToolInterface
	obj = ptr

	var assertPtr ToolInterfacePtr
	assertPtr = ptr

	fmt.Println(obj)
	fmt.Println(assertPtr)
}
~~~

更多的应用在 struct 类型的嵌入字段上的规则，将在 struct 类型这部分进行描述。其他的类型则有一个空的方法集合。在方法集合中，每一个方法都必须具有唯一的、不为空的方法名。

**关于一个类型的方法集合，决定了这个类型实现的接口类型（到底是否实现了某个特定的接口），以及使用指定的接收者调用的方法**。也就是说，通过类型的方法集合，可以知道这个类型到底实现了哪些接口。

# 10 interface

**一个接口类型指代的是一个系列的方法集合，这些方法称之为接口**。一个接口类型的变量可以存储任何一个包含有该接口方法的类型值。这种类型意味着就是“实现”了该接口。

接口是一种引用类型，这种类型的变量的未初始化值是 nil，其组成形式：

~~~go
InterfaceType      = "interface" "{" { ( MethodSpec | InterfaceTypeName ) ";" } "}" .
MethodSpec         = MethodName Signature .
MethodName         = identifier .
InterfaceTypeName  = TypeName .
~~~

接口类型可以通过方法规范显式指定方法，也可以通过接口类型名嵌入其他接口的方法：

~~~go
// A simple File interface.
interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
~~~

每一个显式指定的方法名必须是唯一且非空白的：

~~~go
interface {
	String() string
	String() string  // illegal: String not unique
	_(x int)         // illegal: method must have non-blank name
}
~~~

可能会有多个类型都实现了同一个接口，比如下面的 S1 和 S2 都有这些方法集合：

~~~go
func (p T) Read(p []byte) (n int, err error)
func (p T) Write(p []byte) (n int, err error)
func (p T) Close() error
~~~

因此，对于接口 File 而言，类型 S1 和 S2 都实现了这个接口，不管 S1 和 S2 可能共享或拥有其他的方法。

一个类型可以实现任何接口的子集方法，因此可以实现几个不同的接口。比如，**所有类型都实现下面的空接口类型**：

~~~go
interface {}
~~~

正如在 Java 中，所有类型的（最顶级）父类都是 Object 一样。

类似的，考虑以下此接口规范，该规范出现在类型声明中以定义一个称为 Locker 的接口：

~~~go
type Locker interface {
	Lock()
	Unlock()
}
~~~

如果 S1 和 S2 都实现了：

~~~go
func (p T) Lock() { … }
func (p T) Unlock() { … }
~~~

因此 S1 和 S2 既实现了 File 接口，也实现了 Locker 接口。

一个接口 T 可能会使用类型另一个接口类型名 E 作为方法签名，这种形式称之为**将接口 E 内嵌到接口 T 中**。接口类型 T 的方法集合不仅包含了 T 中显示声明的方法，还包括了嵌入于的接口 E 的方法：

~~~go
type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type Writer interface {
	Write(p []byte) (n int, err error)
	Close() error
}

// ReadWriter's methods are Read, Write, and Close.
type ReadWriter interface {
	Reader  // includes methods of Reader in ReadWriter's method set
	Writer  // includes methods of Writer in ReadWriter's method set
}
~~~

合并的方法集仅包含每个方法集中的方法一次，并且具有相同名称的方法必须具有相同的签名。

~~~go
type ReadCloser interface {
	Reader   // includes methods of Reader in ReadCloser's method set
	Close()  // illegal: signatures of Reader.Close and Close are different
}
~~~

接口类型 T 不能递归嵌入自己或嵌入 T 的任何接口类型：

~~~go
// illegal: Bad cannot embed itself
type Bad interface {
	Bad
}

// illegal: Bad1 cannot embed itself using Bad2
type Bad1 interface {
	Bad2
}
type Bad2 interface {
	Bad1
}
~~~

# 11 map

Map 类型是一种无序的“键-元素”对序列，其中包含有键类型和元素类型。map 类型是一种引用类型，未初始化的类型变量值是 nil，其组成如下：

~~~go
MapType     = "map" "[" KeyType "]" ElementType .
KeyType     = Type .
~~~

键类型必须是可以支持 `==` 和 `!=` 操作符运算的（也就是**可比较的**），因此，键类型一定不能是 function、map、slice。如果键类型是接口类型，键值的动态类型必须支持 `==` 和 `!=` 比较操作符的，否则会导致运行时 panic。

~~~go
map[string]int
map[*T]struct{ x, y float64 }
map[string]interface{}
~~~

map 中“key-元素”对的个数，称之为 map 的长度。对于 map 类型的变量 m 而言，可以使用 len(m) 获取其长度，且该值可能会在程序运行中改变。“key-元素”对可能会在运行中使用赋值表达式添加到 map 中，同时也可以用内建的 delete 函数删除。

可以使用内建的 make 函数新创建一个空的 map：

~~~go
make(map[string]int)
make(map[string]int, 100) // 可选的容量
~~~

初始容量并不是用来限制 map 的可容纳的“key-元素”对的数量，map 的容量会动态增长以容量更多的“key-元素”对。nil 值的 map 变量，意味着空，不同之处在于其不能添加任何“key-元素”对。

~~~go
package main

import "fmt"

func main() {
	maps := make(map[int]string, 3)

	maps[1] = "1"
	maps[2] = "2"
	maps[3] = "3"

	maps[4] = "4"

	fmt.Printf("%d.\n", len(maps)) // 4
}
~~~

# 12 channel

channel 类型被用于 Go 语言中的并发编程模型，提供了一种通过 communicating 共享内存的方式。

channel 类型是一种引用类型，其未初始化的变量值是 nil，其组成如下：

~~~go
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
~~~

可选的操作符 `<-` 用于表示 channel 的方向：发送或接收。如果没有给定任何方向信息，那 channel 则是双向的。在声明表达式或者通过显式转化的方式，可以限制 channel 仅用于发送或接收。

~~~go
chan T          // can be used to send and receive values of type T
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
~~~

操作符 `<-` 尽可能地与它最左侧（leftmost）的 channel 进行组合使用：

~~~go
chan<- chan int    // same as chan<- (chan int)
chan<- <-chan int  // same as chan<- (<-chan int)
<-chan <-chan int  // same as <-chan (<-chan int)
chan (<-chan int)
~~~

channel 类型值的初始化可使用内建的 make 函数，并附带有可选的参数：

~~~go
make(chan int, 100)
~~~

channel 类型值的容量，表示的是 channel 中 buffer 的大小。如果此时的 capacity 为 0 或者不写，channel 将是非缓冲的，只有当 sender 和 receiver 都准备就绪才能通信。反之，channel 称之为缓冲通道，而且当 buffer 未满或未空时都可进行正常的通信而不会阻塞。一个 nil 值的 channel 将永远阻塞，不会用于通信。

channel 可以使用内建的 close() 函数执行关闭操作。从 channel 中接收值的多变量赋值语句可以检测到 channel 是否已经关闭了。

单个 channel 变量在任意多个数量的 goroutine 中执行发送、接收、调用内置 cap 和 len 函数功能时，无需进一步的同步操作。channel 的数据进出类似于 FIFO 结构，比如：在一个 goroutine 中向 channel 中发送了数据，同时在另一个 goroutine 中接收数据，将会按照发送顺序进行接收。

# 13 FAQ

* 有这样的编程问题，将 uint32 类型的值转化为 int 类型的值，实际上是**有风险的**！

  ~~~go
  func (gt *myGoTickets) init(count uint32) bool {
      ...
  	ch := make(chan struct{}, count)
  	i := int(count)
  	for index := 0; index < i; index++ {
  		ch <- struct{}{}
      }
      ...
  	return true
  }
  ~~~

  具体风险在于 uint32 值的范围是：[0x00, 0xFFFFFFFF]，而对于 int 类型来说，不同计算架构下使用了不同长度存储，比如 32 位系统架构下范围是：-0x80000000 ~ 0x7FFFFFFF，如果是 64 位架构下的范围是 -0x8000,0000,0000,0000 ~ 0x7FFF,FFFF,FFFF,FFFF。如果此时 uint32 和 int 都是 32 位的，而且该值大于 0x7FFFFFFF 时，很遗憾，此时转化为 int 类型时，该值就是一个负数了！

  ~~~go
  func main() {
  	var value uint32
  	value = 0x8FFFFFFF
  
  	fmt.Println(value)
  	fmt.Println(int32(value))
  }
  PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
  2415919103
  -1879048193
  ~~~

  此外，关于类型的自增和自减，关于有符号类型也是需要注意的：

  ~~~go
  package main
  
  import (
  	"fmt"
  	"time"
  )
  
  func main() {
  	var value int8
  	value = 0x70
  
  	for {
  		value++
  		fmt.Printf("%d.\n", value)
  		time.Sleep(500 * time.Millisecond)
  	}
  }
  PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
  113. --> 127 --> -128 --> 0 --> 127
  ~~~

  也就是说：如果超出了 int8 类型所能表示的范围，则直接从其最小的负值开始增加。

* Go 中的哪些类型是引用类型，哪些是值类型？

* 什么时候会使用到接口的指针类型？比如 *T，而且 T 是一个接口类型。

* 哪些数据类型是并发安全的？

* 在 buildin.go 源代码中，有如下的类型定义

  ~~~go
  // bool is the set of boolean values, true and false.
  type bool bool
  ~~~

  在 Go 源代码中可以使用 bool 类型（也就是左侧的类型），那右侧的类型在哪里声明的？

* 结构体中各个字段的顺序有什么影响？如果顺序变了，类型也就变了吗？







