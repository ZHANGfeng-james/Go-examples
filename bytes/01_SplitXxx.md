> **刻意练习**，并不是盲目地做事情，而是带着**反馈**结果，有**计划**地实现**更高阶的目标**。

Go 中的 bytes 包下，包含的是处理 `[]byte` 类型实例的方法。bytes 包下包含的功能和 strings 包类似，可以做类比（analogous）。

SplitXxx 系列方法的功能：**按照指定的 sep 对 s 做分割，得到的结果类型是 `[][]byte`**。

# 1 SplitN

`func SplitN(s, sep []byte, n int) [][]byte`：**使用 sep 对 s 拆分操作，返回的切片数组是在 sep 两侧的字节子切片**。如果 sep 为空（比如：`[]byte("")`），则以每一个单个的 UTF-8 字节作为分隔符。函数返回的子切片个数按照如下条件确定：

* n > 0：最多返回 n 个子切片，其中最后一个子切片是不再分割的剩余部分；
* n == 0：返回值是 nil；
* n < 0：尽可能多的返回子切片，也就是尽可能进行分割，得到最大数量的子切片。

测试的**边界**情况，如果 s 是 nil 或者是**空字符串**：只要不是 `n == 0` 的情况，其返回结果都是一个非 nil 的切片类型值。

~~~go
func TestSplitNil(t *testing.T) {
	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN(nil, []byte(""), n)
		if val == nil {
			// n == 0
			fmt.Println("val is nil")
		} else {
			fmt.Println(val)
		}
	}
}
val is nil
[]
[]
[]
[]
[]
~~~

**如果 sep 是空串或者是 nil**，对于函数的返回值，其长度——子切片的个数特别的地方：

~~~go
func TestSplitNNull(t *testing.T) {
	value := "michoi"

	ns := []int{-1, 0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN([]byte(value), []byte(""), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}
["m" "i" "c" "h" "o" "i"]
val is nil
["michoi"]
["m" "ichoi"]
["m" "i" "choi"]
["m" "i" "c" "hoi"]
["m" "i" "c" "h" "o" "i"]
~~~

也就是说，sep 是 nil 或者是空串是相同的情况，其结果值的子串长度随 n 的值变化。而且其结果相当于是将 s 做 explode：

~~~go
// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subslices.
func genSplit(s, sep []byte, sepSave, n int) [][]byte {
	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		return explode(s, n) // 函数命名的意义！
	}
	if n < 0 {
		n = Count(s, sep) + 1
	}
    ...
}
~~~

explode 的含义很明显，就是将 s **按照单个 UTF-8 做拆分**。因此，从这个函数的实现来看，**首要条件是先去判断 n 的值**，然后再去判断 sep 的长度，如果 n 值是 0，直接返回 nil；否则，**再去判断 sep 的长度**，如果 sep 为空或者是 nil，就按照 explode 的规则拆分。n 的值相当于是**对结果 `[][]byte` 的长度**做了**约束**。

**正常情况**下：

~~~go
func TestSplitNPrefix(t *testing.T) {
	value := "michoim"

	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		// if n is 1, result len is 1, r[0] = s
		val := bytes.SplitN([]byte(value), []byte("m"), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%d, %q\n", n, val)
		}
	}
}
val is nil
1, ["michoim"]
2, ["" "ichoim"]
3, ["" "ichoi" ""]
4, ["" "ichoi" ""]
100, ["" "ichoi" ""]

func TestSplitNPrefix(t *testing.T) {
	value := "michoi"

	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		// if n is 1, result len is 1, r[0] = s
		val := bytes.SplitN([]byte(value), []byte("m"), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%d, %q\n", n, val)
		}
	}

	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c"), []byte(","), 2))
	z := bytes.SplitN([]byte("a,b,c"), []byte(","), 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}
val is nil
1, ["michoi"]
2, ["" "ichoi"]
3, ["" "ichoi"]
4, ["" "ichoi"]
100, ["" "ichoi"]
["a" "b,c"]
[] (nil = true)
~~~

特别注意，如果 sep 恰好是 s 的**最后一部分**，那么得到的结果中会包含一个空的子切片；如果 sep 恰好是 s 的开头一部分，那么得到的结果中会同样也会包含一个空的切片。

# 2 Split

`func Split(s, sep []byte) [][]byte`：同上，类似于 n 值是 -1 的返回结果。

~~~go
func TestSplitUsage(t *testing.T) {
	value := "michoi"
	fmt.Printf("%q\n", bytes.Split([]byte(value), []byte("")))
}
["m" "i" "c" "h" "o" "i"]
~~~

# 3 SplitAfterN

`func SplitAfterN(s, sep []byte, n int) [][]byte`：是和 `SplitN` 类似的，区别在于是在 sep 之后拆分。

~~~go
func TestSplitAfterN(t *testing.T) {
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c"), []byte(","), 3))
	fmt.Printf("%q\n", bytes.SplitAfterN([]byte("a,b,c"), []byte(","), 3))
}
["a" "b" "c"]
["a," "b," "c"]


func TestSplitAfterN(t *testing.T) {
	value := "michoi"

	ns := []int{0, 1, 2, 3, 4, 100}
	for n := range ns {
		val := bytes.SplitAfterN([]byte(value), []byte("i"), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}
val is nil
["michoi"]
["mi" "choi"]
["mi" "choi" ""]
["mi" "choi" ""]
["mi" "choi" ""]

// bytes.SplitAfterN --> bytes.SplitN
val is nil
["michoi"]
["m" "choi"]
["m" "cho" ""]
["m" "cho" ""]
["m" "cho" ""]
~~~

区别之处在于是否包含这个 sep 的内容。

# 4 SplitAfter

`func SplitAfter(s, sep []byte) [][]byte`：和 SplitAfter 类似，相当于 n 值为 -1。

~~~go
func TestSplitAfter(t *testing.T) {
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte(",")))
	fmt.Printf("%q\n", bytes.SplitAfter([]byte("a,b,c"), []byte(",")))
}
["a" "b" "c"]
["a," "b," "c"]
~~~

进一步看出 SplitAfter 和 Split 的差异！