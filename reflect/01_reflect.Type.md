~~~go
func TypeOf(i interface{}) Type
~~~

通过这个方法，获取到 i 的**动态类型信息**，也就是**一个 Type 实例**。当然 Type 本身是一个**接口**，定义了一系列的方法。

Type 接口代表的是 Go 中的一个类型。Type 接口中定义的所有方法，并非适用于 Go 中可用的所有类型。在调用类型相关方法前，可使用 Kind 方法获知其类型种类信息。如果调用了和该类型种类不合适的方法，会导致 panic。

Type 实例是可比较的，因此其可以作为 map 的 key。如果两个 Type 表示的是相同的 Go 类型，那么这两个 Type 实例就是相等的。

**指针变量的实例**和**普通变量的实例**的使用：

~~~go
type Account struct {
	username string
	age      int8
}

func (account *Account) GetAge() {
}

func (account Account) GetUsername() {
}

func main() {
	ptr := &Account{}
	fmt.Printf("%T\n", ptr) // *main.Accountp ptr是一个指针类型的变量

	typ := reflect.TypeOf(ptr)
	fmt.Println(typ.Name(), typ.NumMethod(), typ.Kind())
	// ValueOf returns a new Value initialized to the concrete value stored in the interface i.
	value := reflect.ValueOf(ptr)
    fmt.Println(value.Kind(), value.Type(), value.NumMethod())
}

*main.Account
 2 ptr
ptr *main.Account 2
~~~

特别的，程序第 6 行，typ 是无法获取到 Name 信息的：`Name returns the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.`

相对的：

~~~go
func normalTest() {
	obj := Account{}
    fmt.Printf("%T\n", obj)
    
	value := reflect.ValueOf(obj)
	fmt.Println(value.Kind(), value.Type().Name(), value.NumField(), value.NumMethod())

	typ := value.Type()
	fmt.Println(typ.Name(), typ.Kind(), typ.NumField(), typ.NumMethod())
}

main.Account
struct Account 2 1
Account struct 2 1
~~~

对于 `*main.Account` 类型来说，获取到的 `NumMethod()` 是 2 个；而对于 `main.Account` 类型来说，获取到的是 1 个。

如果程序中拿到的是 `*main.Account`，如何获取对应的类型信息：

~~~go
func ptrIndirect() {
	ptr := &Account{
		username: "Katyusha",
		age:      18,
	}
	typ := reflect.Indirect(reflect.ValueOf(ptr)).Type()
	fmt.Println(typ.Name(), typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name, field.Tag)
	}
}

Account 2
username geekorm:"PRIMARY KEY"
age 
~~~

`reflect.ValueOf(ptr)` 创建 `reflect.Value` 类型的值，`Value is the reflection interface to a Go value.`

通过执行 `reflect.Indirect(reflect.Value)` 相当于是：`Indirect returns the value that v points to.` 执行了一次指针变量的**间接访问操作**，即通过指针（变量的地址）访问到了对应的变量，得到的结果类型是 `reflect.Value`。

如果想通过一个 `reflect.Type` 构造出 `reflect.Value`：

~~~go
func reflectNew() {
	ptr := Account{}

	// start: reflect.Type，对应的是 mani.Account 类型
	typ := reflect.TypeOf(ptr)

	valuePtr := reflect.New(typ) // 其值对应的是 *main.Account 类型的值

	// end: reflect.Value
	value := reflect.Indirect(valuePtr) // 通过间接访问，即间接访问 *main.Account 变量
	fmt.Println(value.Kind(), value.NumField(), value.NumMethod())
}
~~~

间接使用：`New returns a Value representing a pointer to a new zero value for the specified type.`

解析出一个结构体变量的各个字段值：

~~~go
type Account struct {
	Username string `geekorm:"PRIMARY KEY"`
	Age      int8
}

func recordValues() {
	account := &Account{
		Username: "Katyusha",
		Age:      18,
	}

	// 获取 Fields 数组
	typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type()
	fmt.Println(typ.NumField())

	value := reflect.Indirect(reflect.ValueOf(account))
	for i := 0; i < typ.NumField(); i++ {
		// 依据 field.Name 获取对应的值，此次 Field 必须是可导出的
		v := value.FieldByName(typ.Field(i).Name).Interface()
		fmt.Printf("fieldName[%d]=%s, value:%v\n", i, typ.Field(i).Name, v)
	}
}
~~~

`value.FieldByName(typ.Field(i).Name)` 返回值类型是 `reflect.Value`，也就是对应结构体变量的对应字段的 `reflect.Value` 值。为了拿到这个值，还要作一次转换：`reflect.Value` 转化成 `interface{}`，使用的是 `Interface()` 方法：`Interface returns v's current value as an interface{}.` 相当于是：

~~~go
var i interface{} = (v's underlying value)
~~~

`reflect.Value` 和 `reflect.Type` 都有 `Elem()`，区别是什么：

~~~go
func elemAndInterface() {
	var accounts []Account

	var ptr interface{}
	ptr = &accounts

	destSlice := reflect.Indirect(reflect.ValueOf(ptr)) // reflect.Value []Acccount
	destType := destSlice.Type().Elem()                 // Account
	fmt.Println(destType.Name())

	value := reflect.New(destType).Elem() // reflect.Value
	value.Interface()                     // interface{}
}
~~~

对于 `reflect.Type` 来说，`Elem()` 方法获取到的是 `Array, Chan, Map, Ptr, or Slice` 元素类型，返回值类型是 `reflect.Type`；对于 `reflect.Value` 来说，`Elem()` 方法是 `Elem returns the value that the interface v contains or that the pointer v points to.` 类似于获取其底层的值，而返回值类型是 `reflect.Value`。

`reflect.Indirect` 函数的迷惑性：

~~~go
func indirect() {
	account := Account{
		Username: "Katyusha",
		Age:      18,
	}
	value := reflect.Indirect(reflect.ValueOf(account))
	fmt.Println(value.FieldByName("Username"))

	fmt.Println(value.Type().Name())
}
~~~

不过该函数的注释部分写得很明确：

~~~go
// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func Indirect(v Value) Value {
	if v.Kind() != Ptr {
		return v
	}
	return v.Elem()
}
~~~



# 1 结构体的标签信息

在处理 json 格式字符串的时候，经常会看到声明 struct 结构的时候，字段属性的右侧还有一些描述信息。比如：

~~~go
type User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
~~~

这些描述信息是用**反引号**括起来的。**这其中的内容有什么作用？**

要比较详细的了解这个，要先了解一下 golang 的基础，在 golang 中，命名都是推荐都是用驼峰方式，并且在首字母大小写有特殊的语法含义：包外无法引用。但是由于经常需要和其它的系统进行数据交互，例如转成 json 格式，存储到 mongodb 啊等等。这个时候，如果用属性名来作为键值可能**不一定会符合项目要求**。所以，就多了反引号的内容，在 golang 中叫**标签（Tag）**，在转换成其它数据格式的时候，会使用其中特定的字段作为键值。例如上例在转成 json 格式：

~~~go
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func main() {
	user := &User{
		UserId:   1,
		UserName: "Michoi",
	}

	j, _ := json.Marshal(user)
	fmt.Println(string(j))
}
PS E:\GoUsage> go run main.go
{"user_id":1,"user_name":"Michoi"}
user_id
user_name
~~~

如果在属性中**不增加标签说明**，则输出：`{"UserId":1,"UserName":"Michoi"}`，可以看到直接用 struct 的**属性名**做**键值**。

这样，如果处理的是 HTTP 请求，那可直接将 Marshal 之后的结果输出到 Client 端，这个格式是符合 JSON 格式的：`{"user_id":1,"user_name":"Michoi"}`。

另外，当我们需要自己封装一些操作，需要**用到 Tag 中的内容时，咋样去获取呢？**这边可以使用**反射包（reflect）**中的方法来获取：

~~~go
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func main() {
	user := &User{
		UserId:   1,
		UserName: "Michoi",
	}
	t := reflect.TypeOf(user)
	idField := t.Elem().Field(0)
	fmt.Println(idField.Tag.Get("json"))

	nameField := t.Elem().Field(1)
	fmt.Println(nameField.Tag.Get("json"))
}
PS E:\GoUsage> go run main.go
user_id
user_name
~~~

# 2 reflect.Type 类型的应用

上述获取 struct 的**标签信息**，可以通过如下代码实现：

~~~go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func main() {
	ptr := &User{
		UserId:   10,
		UserName: "Michio",
	}

	var i interface{} = ptr
	// 调用 Elem() 时，是接口或者指针（对调用类型有要求）
	typeof := reflect.TypeOf(i)
	fmt.Println("kind:", typeof.Kind())            // kind: ptr
	fmt.Println(typeof, fmt.Sprintf("%T", typeof)) // *main.User *reflect.rtype

	v := typeof.Elem()
	fmt.Println(v, fmt.Sprintf("%T", v)) // main.User *reflect.rtype
	for i := 0; i < v.NumField(); i++ {
		// 返回 reflect.StructField 类型值，类似于：{UserId  int json:"user_id" 0 [0] false}
		fieldInfo := v.Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("json")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fmt.Println(name, v.Field(i))
	}
}
PS E:\goreflect> go run .\struct.go
kind: ptr ; name: 0
*main.User *reflect.rtype
main.User *reflect.rtype
user_id {UserId  int json:"user_id" 0 [0] false}
user_name {UserName  string json:"user_name" 8 [1] false}
~~~

下面一步步做解释：

1. `reflect.TypeOf(i)`：注意此处 i 是一个 `interface{}` 类型的值，其中存放的是**一个结构体指针值**，也就是 ptr。该函数的返回值，是一个 `*reflect.rtype` 类型（其本身是实现了 reflect.Type 接口），其含义表示了 i 的**动态类型值**——`*main.User`。
2. `typeof.Elem()`：调用 `Elem()` 方法，**获取其元素的类型**。其调用者 typeof 必须是 Array, Chan, Map, Ptr, or Slice 类型中的一种，此时 typeof 是 `*main.User` 类型。其返回值类型仍然是 `*reflect.rtype` 类型，其真实类型是 `main.User`。
3. `v.NmuField()`：获取 v 这个结构体类型的字段数量。v 的类型**必须是结构体类型**，因此，此处必须调用 `typeof.Elem()`，其返回值类型是 `main.User` 结构体类型。
4. `v.Field(i)`：获取第 i 个结构体类型值，其返回值类型是 StructField。潜在的，v 的类型必须是结构体类型，否则会 panic。
5. `fieldInfo.Tag`：获取当前 Field 的标签信息，其返回值类型是 StructTag，对应的就是 string 类型。
6. `tag.Get("json")`：StructTag 类型有 Get 和 Lookup 方法，用于获取 key 对应的 value。比如 `json:"user_id"` 这个 Tag，key 值就是 json，对应的 value 是 user_id。

# 3 reflect.Type 结构体定义

下面，详细看看 reflect.Type 的结构体：

~~~go
// Type is the representation of a Go type.
//
// Not all methods apply to all kinds of types. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of type before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run-time panic.
//
// Type values are comparable, such as with the == operator,
// so they can be used as map keys.
// Two Type values are equal if they represent identical types.
type Type interface {
	// Methods applicable to all types.

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It panics if i is not in the range [0, NumMethod()).
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver,
	// and only exported methods are accessible.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	//
	// Methods are sorted in lexicographic order.
	Method(int) Method

	// MethodByName returns the method with that name in the type's
	// method set and a boolean indicating if the method was found.
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	MethodByName(string) (Method, bool)

	// NumMethod returns the number of methods accessible using Method.
	//
	// Note that NumMethod counts unexported methods only for interface types.
	NumMethod() int

	// Name returns the type's name within its package for a defined type.
	// For other (non-defined) types it returns the empty string.
	Name() string

	// PkgPath returns a defined type's package path, that is, the import path
	// that uniquely identifies the package, such as "encoding/base64".
	// If the type was predeclared (string, error) or not defined (*T, struct{},
	// []int, or A where A is an alias for a non-defined type), the package path
	// will be the empty string.
	PkgPath() string

	// Size returns the number of bytes needed to store
	// a value of the given type; it is analogous to unsafe.Sizeof.
	Size() uintptr

	// String returns a string representation of the type.
	// The string representation may use shortened package names
	// (e.g., base64 instead of "encoding/base64") and is not
	// guaranteed to be unique among types. To test for type identity,
	// compare the Types directly.
	String() string

	// Kind returns the specific kind of this type.
	Kind() Kind

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	// AssignableTo reports whether a value of the type is assignable to type u.
	AssignableTo(u Type) bool

	// ConvertibleTo reports whether a value of the type is convertible to type u.
	ConvertibleTo(u Type) bool

	// Comparable reports whether values of this type are comparable.
	Comparable() bool

	// Methods applicable only to some types, depending on Kind.
	// The methods allowed for each kind are:
	//
	//	Int*, Uint*, Float*, Complex*: Bits
	//	Array: Elem, Len
	//	Chan: ChanDir, Elem
	//	Func: In, NumIn, Out, NumOut, IsVariadic.
	//	Map: Key, Elem
	//	Ptr: Elem
	//	Slice: Elem
	//	Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

	// Bits returns the size of the type in bits.
	// It panics if the type's Kind is not one of the
	// sized or unsized Int, Uint, Float, or Complex kinds.
	Bits() int

	// ChanDir returns a channel type's direction.
	// It panics if the type's Kind is not Chan.
	ChanDir() ChanDir

	// IsVariadic reports whether a function type's final input parameter
	// is a "..." parameter. If so, t.In(t.NumIn() - 1) returns the parameter's
	// implicit actual type []T.
	//
	// For concreteness, if t represents func(x int, y ... float64), then
	//
	//	t.NumIn() == 2
	//	t.In(0) is the reflect.Type for "int"
	//	t.In(1) is the reflect.Type for "[]float64"
	//	t.IsVariadic() == true
	//
	// IsVariadic panics if the type's Kind is not Func.
	IsVariadic() bool

	// Elem returns a type's element type.
	// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	Elem() Type

	// Field returns a struct type's i'th field.
	// It panics if the type's Kind is not Struct.
	// It panics if i is not in the range [0, NumField()).
	Field(i int) StructField

	// FieldByIndex returns the nested field corresponding
	// to the index sequence. It is equivalent to calling Field
	// successively for each index i.
	// It panics if the type's Kind is not Struct.
	FieldByIndex(index []int) StructField

	// FieldByName returns the struct field with the given name
	// and a boolean indicating if the field was found.
	FieldByName(name string) (StructField, bool)

	// FieldByNameFunc returns the struct field with a name
	// that satisfies the match function and a boolean indicating if
	// the field was found.
	//
	// FieldByNameFunc considers the fields in the struct itself
	// and then the fields in any embedded structs, in breadth first order,
	// stopping at the shallowest nesting depth containing one or more
	// fields satisfying the match function. If multiple fields at that depth
	// satisfy the match function, they cancel each other
	// and FieldByNameFunc returns no match.
	// This behavior mirrors Go's handling of name lookup in
	// structs containing embedded fields.
	FieldByNameFunc(match func(string) bool) (StructField, bool)

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type

	// Key returns a map type's key type.
	// It panics if the type's Kind is not Map.
	Key() Type

	// Len returns an array type's length.
	// It panics if the type's Kind is not Array.
	Len() int

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	NumField() int

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Type

	common() *rtype
	uncommon() *uncommonType
}
~~~

在深入查看 reflect.Type 接口中各个方法之前，**特别需要注意**：

并非所有的 reflect.Type 接口值能调用上述的所有方法，在调用其中的方法之前，最好查看 Type 的 Kind 信息。比如 `Elem() ` 这个接口方法，其调用的限制是，Type 必须是 Array, Chan, Map, Ptr, or Slice。如果不是上述这些类型，就会报出运行时 panic。

下面从可调用的角度出发，看看 Type 接口中的方法：

* **所有的 Type 均可调用的方法**
  * Align() int：获取定义该类型的值时，内存占用的字节数；比如 `var ele int32 = 100` 这个变量，`reflect.TypeOf(i).Align()` 对应的值是 4。
  * FieldAlign() int：当此类型在结构体中作为字段时，需要使用的字节数；
  * Method(int) Method：获取第 i 序位的方法；
  * MethodByName(string) (Method, bool)：查找指定名称的方法，其返回的 bool 值表示是否找到。
  * NumMethod() int：获取类型对应的方法数量，并且是可导出方法；
  * Name() string：获取类型名字符串（附加上包路径）；
  * PkgPath() string：获取该类型值对应的包路径名；
  * Size() uintptr：返回存储该类型值所需要的内存字节数；
  * String() string：获取类型的描述字符串；
  * Kind() Kind：获取此类型值特有的 Kind 信息；
  * Implements(u Type) bool：类型的值是否实现了指定的 u 类型；
  * AssignableTo(u Type) bool：类型的值是否可被赋值给 u 类型；
  * ConvertibleTo(u Type) bool：类型的值是否可转换成指定的 u 类型；
  * Comparable() bool：类型的值是否可比较；
* **依据 Kind 的不同，其可调用不同的方法**
  * `Int*, Uint*, Float*, Complex*`，比如 int8、int32 等：Bits
  * `Array`：Elem, Len
  * `Chan`：ChanDir, Elem
  * `Func`：In, NumIn, Out, NumOut, IsVariadic.
  * `Map`：Key, Elem
  * `Ptr`：Elem
  * `Slice`：Elem
  * `Struct`：Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

上述根据 Kind 的不同，可调用的方法功能如下：

* Bits() int：获取类型对应的二进制位数，比如 int32 类型，对应的返回值是 32；
* ChanDir() ChanDir：获取 channel 类型值的方向；
* IsVariadic() bool：判断其类型的入参是否包含一个可变参数；
* Elem() Type：获取类型值的**元素类型**；
* Field(i int) StructField：获取结构体类型中，第 i 序位的结构体字段类型；
* FieldByIndex(index []int) StructField：按照 index 的值，连续调用 Field(i int)，相当于是多层成员结构访问；
* FieldByName(name string) (StructField, bool)：获取指定 name 的结构体字段类型，其中 name 实际上就是字段名；
* FieldByNameFunc(match func(string) bool) (StructField, bool)：获取满足指定方法签名的结构体字段类型；
* In(i int) Type：返回函数类型值的第 i 序位的入参类型值；
* Key() Type：获取 Map 类型值的 key 类型值；
* Len() int：获取数组类型值的长度；
* NumField() int：获取结构体类型值的字段数量；
* NumIn() int：获取函数类型值的入参数量；
* NumOut() int：获取函数类型值的返回值数量；
* Out(i int) Type：获取函数类型值的第 i 序位的返回值类型。

# 4 StructField

对于结构体类型值，还封装了 StructField 类型，即：结构体字段的类型

~~~go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag string
~~~

其中：

* Name：字段的名称；
* Type：字段的类型信息；
* Tag：字段标签信息，其类型是 StructTag。其中 StructTag 定义了 `Get(key string)` 和 `Lookup(key string)` 方法。

# 5 Method

一个类型，对应会定义该类型的方法。其中方法的类型封装，如下：

~~~go
// Method represents a single method.
type Method struct {
	// Name is the method name.
	// PkgPath is the package path that qualifies a lower case (unexported)
	// method name. It is empty for upper case (exported) method names.
	// The combination of PkgPath and Name uniquely identifies a method
	// in a method set.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	Name    string
	PkgPath string

	Type  Type  // method type
	Func  Value // func with receiver as first argument
	Index int   // index for Type.Method
}
~~~

其中：

* Name：方法的名称；
* Type：方法的类型
* Func：方法的接收者，是 reflect.Value 类型值。

# 6 Kind

reflect.Type 代表的类型值，在 Go 中对应的类型索引，如下所示。各个类型的 Kind 值都是一个 uint 类型值：

~~~go
// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
~~~

