`./builtin/builtin.go` 文件中定义了 Go 语言内置的函数：

1. `append`
2. `copy`
3. `delete`
4. `len`
5. `cap`
6. `make`
7. `new`
8. `close`
9. `panic`
10. `recover`
11. `print`
12. `println`

下面按照顺序依次分析：

# 1 append

Slice ==增长的算法（主要是 append 函数）==是怎样的？如果跟踪 Go 源代码，获取到 append 函数对 Slice 容量的影响？

append 函数的原型：

~~~go
// The append built-in function appends elements to the end of a slice. If
// it has sufficient capacity, the destination is resliced to accommodate the
// new elements. If it does not, a new underlying array will be allocated.
// Append returns the updated slice. It is therefore necessary to store the
// result of append, often in the variable holding the slice itself:
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
// As a special case, it is legal to append a string to a byte slice, like this:
//	slice = append([]byte("hello "), "world"...)
func append(slice []Type, elems ...Type) []Type
~~~

根据函数原型的注解可看出，如果原 slice 足够容纳新数据，则会向原 slice 中增加新数据；如果容量不够，则会**在底层重新开辟内存空间**。比如：

~~~go
func AppendByte(slice []byte, data ...[]byte) {
    m := len(slice)
    n := m + len(data)
    if n > cap(slice) {
        // If it does not, a new underlying array will be allocated.
        newSlice := make([]byte, (n + 1) * 2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n] // growth len
    copy(slice[m:n], data)
    return slice
}
~~~

Go 提供了一个内置的 append 函数，实现上述的 Slice 容量的增长。我们先通过示例程序来试探 append 函数：

~~~go
func main() {
	var sa []string
	if sa == nil {
		fmt.Println("sa == nil", len(sa))
	}

	for i := 0; i < 1001; i++ {
        // 仅增加一个元素时
		sa = append(sa, strconv.Itoa(i))
		printSliceData(sa)
	}

	printSliceData(sa)
}

func printSliceData(s []string) {
	fmt.Printf("地址：%p \t len:%d \t cap:%d \t \n", s, len(s), cap(s))
}
~~~

随着 slice 中元素的增加，cap(slice) 的增长有一套算法支持。

# 2 copy

下面是 copy 函数的定义：

~~~go
// The copy built-in function copies elements from a source slice into a
// destination slice. (As a special case, it also will copy bytes from a
// string to a slice of bytes.) The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
func copy(dst, src []Type) int
~~~

在上述定义中，隐藏了 `[]Type` 类型的 `len` 含义，如果 dst 的 len 值为 0，则不会发生任何拷贝动作。

为了增长一个切片的容量，我们必须创建一个新的并且容量更大的切片并且将原来切片中的数据复制到新的切片中。这就是其他语言动态数组背后进行的操作。下面这个例子就会通过创建一个新的切片 t，然后把 s 切片中的数据复制到 t 之中，再把 t 赋值给 s 来对 s 进行容量翻倍操作：

~~~go
t := make([]byte, len(s), (cap(s) + 1) * 2)
for i := range s {
    t[i] = s[i]
}
s = t
~~~

向代码中的 for 循环来赋值的操作可以通过 Go 内置的函数 copy 来操作：

~~~go
t := make([]byte, len(s), (cap(s) + 1) * 2)
copy(t, s)
s = t
~~~

用 2 个例子来说明：

~~~go
func main() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}

	copy(slice1, slice2)
	fmt.Println(slice1, slice2)  // [5 4 3 4 5] [5 4 3]
}
~~~

上述代码中：只会复制 slice2 的 3 个元素到 slice1 的前 3 个位置。

~~~go
func main() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}

	copy(slice2, slice1)
	fmt.Println(slice1, slice2)  // [1 2 3 4 5] [1 2 3]
}
~~~

上述代码中：只会复制 slice1 的前 3 个元素到 slice2 中。

另一个能说明问题的示例代码：

~~~go
func main() {
	slice1 := []int{1, 2, 3, 4, 5}
	var slice2 []int

	copy(slice2, slice1)
	fmt.Println(slice1, slice2) // [1 2 3 4 5] []

	fmt.Printf("%p, %p\n", slice1, slice2) // 0xc00000a4b0, 0x0
}
~~~

`slice2` 是一个 nil 切片，其 cap(slice2) 和 len(slice2) 都是 0，作为 copy() 的 dst 时不会发生拷贝。