# 1 reflect.Value 的应用

仍然使用**获取 struct 标签信息**这个应用实例作为引子：

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

	valueOf := reflect.ValueOf(i) // the pointer variable
	fmt.Println(valueOf, fmt.Sprintf("%T", valueOf))
	ele := valueOf.Elem() // the struct variable
	fmt.Println(ele, fmt.Sprintf("%T", ele))

	for i := 0; i < ele.NumField(); i++ {
		fieldInfo := ele.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("json")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fmt.Println(name, ele.Field(i), fieldInfo)
	}
}
PS E:\goreflect> go run .\struct.go
&{10 Michio} reflect.Value
{10 Michio} reflect.Value
user_id 10 {UserId  int json:"user_id" 0 [0] false}
user_name Michio {UserName  string json:"user_name" 8 [1] false}
~~~

和 reflect.Type 不同，reflect.Value 类型值的输出信息类似于 `&{10 Michio}`，而前者输出的是 `*main.User`。这也就是对应了一个接口变量的**动态值**和**动态类型**的不同。



# 2 Value 类型

Go 中的 reflect 包，关于 Value 类型：

~~~go
// Value is the reflection interface to a Go value.
//
// Not all methods apply to all kinds of values. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of value before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run time panic.
//
// The zero Value represents no value.
// Its IsValid method returns false, its Kind method returns Invalid,
// its String method returns "<invalid Value>", and all other methods panic.
// Most functions and methods never return an invalid value.
// If one does, its documentation states the conditions explicitly.
//
// A Value can be used concurrently by multiple goroutines provided that
// the underlying Go value can be used concurrently for the equivalent
// direct operations.
//
// To compare two Values, compare the results of the Interface method.
// Using == on two Values does not compare the underlying values
// they represent.
type Value struct {
	// typ holds the type of the value represented by a Value.
	typ *rtype

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or typ.pointers() is true.
	ptr unsafe.Pointer

	// flag holds metadata about the value.
	// The lowest bits are flag bits:
	//	- flagStickyRO: obtained via unexported not embedded field, so read-only
	//	- flagEmbedRO: obtained via unexported embedded field, so read-only
	//	- flagIndir: val holds a pointer to the data
	//	- flagAddr: v.CanAddr is true (implies flagIndir)
	//	- flagMethod: v is a method value.
	// The next five bits give the Kind of the value.
	// This repeats typ.Kind() except for method values.
	// The remaining 23+ bits give a method number for method values.
	// If flag.kind() != Func, code can assume that flagMethod is unset.
	// If ifaceIndir(typ), code can assume that flagIndir is set.
	flag

	// A method value represents a curried method invocation
	// like r.Read for some receiver r. The typ+val+flag bits describe
	// the receiver r, but the flag's Kind bits say Func (methods are
	// functions), and the top bits of the flag give the method number
	// in r's type's method table.
}
~~~

和 reflect.Type 类似，并不是所有的 Value 都可调用如下的方法。比如 `Elem() Value`，获取接口 value 包含的或者指针指向的 Value 值，value 的 Kind 必须是 Interface 或 Ptr。Value 类型对应实现的方法：

* func (v Value) Addr() Value
* func (v Value) Bool() bool
* func (v Value) Bytes() []byte
* func (v Value) Call(in []Value) []Value
* func (v Value) CallSlice(in []Value) []Value
* func (v Value) CanAddr() bool
* func (v Value) CanInterface() bool
* func (v Value) CanSet() bool
* func (v Value) Cap() int
* func (v Value) Close()
* func (v Value) Complex() complex128
* func (v Value) Convert(t Type) Value
* func (v Value) Elem() Value：
* func (v Value) Field(i int) Value：获取结构体 Value 的第 i 序位的 reflect.Value 类型值；
* func (v Value) FieldByIndex(index []int) Value
* func (v Value) FieldByName(name string) Value
* func (v Value) FieldByNameFunc(match func(string) bool) Value
* func (v Value) Float() float64
* func (v Value) Index(i int) Value
* func (v Value) Int() int64
* func (v Value) Interface() (i interface{})
* func (v Value) InterfaceData() [2]uintptr
* func (v Value) IsNil() bool
* func (v Value) IsValid() bool
* func (v Value) IsZero() bool
* func (v Value) Kind() Kind：获取 Value 的 Kind；
* func (v Value) Len() int
* func (v Value) MapIndex(key Value) Value
* func (v Value) MapKeys() []Value
* func (v Value) MapRange() *MapIter
* func (v Value) Method(i int) Value
* func (v Value) MethodByName(name string) Value
* func (v Value) NumField() int
* func (v Value) NumMethod() int
* func (v Value) OverflowComplex(x complex128) bool
* func (v Value) OverflowFloat(x float64) bool
* func (v Value) OverflowInt(x int64) bool
* func (v Value) OverflowUint(x uint64) bool
* func (v Value) Pointer() uintptr
* func (v Value) Recv() (x Value, ok bool)
* func (v Value) Send(x Value)
* func (v Value) Set(x Value)
* func (v Value) SetBool(x bool)
* func (v Value) SetBytes(x []byte)
* func (v Value) SetCap(n int)
* func (v Value) SetComplex(x complex128)
* func (v Value) SetFloat(x float64)
* func (v Value) SetInt(x int64)
* func (v Value) SetLen(n int)
* func (v Value) SetMapIndex(key, elem Value)
* func (v Value) SetPointer(x unsafe.Pointer)
* func (v Value) SetString(x string)
* func (v Value) SetUint(x uint64)
* func (v Value) Slice(i, j int) Value
* func (v Value) Slice3(i, j, k int) Value
* func (v Value) String() string
* func (v Value) TryRecv() (x Value, ok bool)
* func (v Value) TrySend(x Value) bool
* func (v Value) Type() Type：获取 Value 的类型值；
* func (v Value) Uint() uint64
* func (v Value) UnsafeAddr() uintptr