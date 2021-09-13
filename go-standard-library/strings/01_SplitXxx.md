Go 中的标准库 strings，用于处理 UTF-8 编码的字符串，可将 strings 包当作是一个（UTF-8编码格式的）字符串处理工具箱。



`func Split(s, sep string) []string`：使用 sep 对 s 拆分，拆分得到的是 []string 类型的结果。比如 `中-国-人`，字符串中各个中文字符中间间隔的是 `-` 字符，可用 `-` 字符对上述字符拆分

~~~go
func TestSplitNormal(t *testing.T) {
	values := "中-国-人"

	// 1 0
	fmt.Printf("%s\n", strings.Split(values, ""))
	// 1 1
	fmt.Printf("%s\n", strings.Split(values, "-"))
	fmt.Printf("%s\n", strings.Split("中国人", "-"))

	// 0 1
	fmt.Printf("%q\n", strings.Split("", "-"))

	// 0 0
	fmt.Printf("%q\n", strings.Split("", ""))

	result := strings.Split("", "-")
	fmt.Println(len(result))

	result = strings.Split("", "")
	fmt.Println(len(result))
    
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))
	fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))
}
[中 - 国 - 人]
[中 国 人]
[中国人]
[""]
[]
1
0
["a" "b" "c"]
["" "man " "plan " "canal panama"]
[" " "x" "y" "z" " "]
[""]
~~~

由此可以归纳出上述**拆分的规则**：

1. 如果 sep 为空，分为 2 种情况：
   * 如果 s 不为空，会将 s 拆分成单个 UTF-8 编码的字符；
   * 如果 s 为空，返回的是一个确定的值，内容为空的 []string 值。
2. 如果 sep 不为空，分为 2 种情况：
   * 如果 s 不为空，则按照 sep 进行拆分。如果 s 中不包含 sep，则返回包含 s 的 []string 值；
   * 如果 s 为空，返回包含有元素个数为 1的——内容为空的字符串——的 []string 值。

另外一个测试示例：

~~~go
func TestSplitFoo(t *testing.T) {
	values := "fof"
	fmt.Printf("%q\n", strings.Split(values, "f"))
}
["" "o" ""]
~~~

很具有代表性，使用 `f` 字符串对 `fof` 拆分，其结果为：`["" "o" ""]`

`func SplitAfter(s, sep string) []string`：大部分规则和 Split 是相同的，不同之处是拆分出的 []string 元素值会包含 sep

~~~go
func TestSplitAfter(t *testing.T) {
	values := "中-国-人"

	// 1 0
	fmt.Printf("%s\n", strings.SplitAfter(values, ""))
	// 1 1
	fmt.Printf("%s\n", strings.SplitAfter(values, "-"))
	fmt.Printf("%s\n", strings.SplitAfter("中国人", "-"))

	// 0 1
	fmt.Printf("%q\n", strings.SplitAfter("", "-"))

	// 0 0
	fmt.Printf("%q\n", strings.SplitAfter("", ""))

	result := strings.SplitAfter("", "-")
	fmt.Println(len(result))

	result = strings.SplitAfter("", "")
	fmt.Println(len(result))

	fmt.Printf("%q\n", strings.SplitAfter("a,b,c", ","))
	fmt.Printf("%q\n", strings.SplitAfter("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.SplitAfter(" xyz ", ""))
	fmt.Printf("%q\n", strings.SplitAfter("", "Bernardo O'Higgins"))
}
[中 - 国 - 人]
[中- 国- 人]
[中国人]
[""]
[]
1
0
["a," "b," "c"]
["a " "man a " "plan a " "canal panama"]
[" " "x" "y" "z" " "]
[""]
~~~

`func SplitN(s, sep string, n int) []string`：拆分方式和 Split 相同，不同之处是限定了拆分后结果元素个数

~~~go
func TestSplitN(t *testing.T) {
	values := "中-国-人"
	// 1 0
	fmt.Printf("%q\n", strings.SplitN(values, "", 10))
	// 1 1
	fmt.Printf("%q\n", strings.SplitN(values, "-", 10))
	fmt.Printf("%q\n", strings.SplitN("中国人", "-", 10))

	// 0 1
	fmt.Printf("%q\n", strings.SplitN("", "-", 10))

	// 0 0
	fmt.Printf("%q\n", strings.SplitN("", "", 10))

	fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))
	z := strings.SplitN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}
["中" "-" "国" "-" "人"]
["中" "国" "人"]
["中国人"]
[""]
[]
["a" "b,c"]
[] (nil = true)
~~~

结果值的元素个数判断标准：

* n > 0 时：最多返回 n 个 string 值的 []string；
* n == 0 时：nil
* n < 0 时：不做限制，对 s 做全拆分，类似于调用了 Split 函数。

`func SplitAfterN(s, sep string, n int) []string`：拆分方式和 SplitAfter 类似，不同之处是限定了拆分后结果元素个数

~~~go
func TestSplitAfterN(t *testing.T) {
	values := "中-国-人"
	// 1 0
	fmt.Printf("%q\n", strings.SplitAfterN(values, "", 10))
	// 1 1
	fmt.Printf("%q\n", strings.SplitAfterN(values, "-", 10))
	fmt.Printf("%q\n", strings.SplitAfterN("中国人", "-", 10))

	// 0 1
	fmt.Printf("%q\n", strings.SplitAfterN("", "-", 10))

	// 0 0
	fmt.Printf("%q\n", strings.SplitAfterN("", "", 10))

	fmt.Printf("%q\n", strings.SplitAfterN("a,b,c", ",", 2))
	z := strings.SplitAfterN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}
["中" "-" "国" "-" "人"]
["中-" "国-" "人"]
["中国人"]
[""]
[]
["a," "b,c"]
[] (nil = true)
~~~

SplitXxx 系列包括如下函数：

* `func Split(s, sep string) []string`：使用 sep 对 s 拆分，拆分得到的是 []string 类型的结果；
* `func SplitAfter(s, sep string) []string`：大部分规则和 Split 是相同的，不同之处是拆分出的 []string 元素值会包含 sep；
* `func SplitAfterN(s, sep string, n int) []string`：拆分方式和 SplitAfter 类似，不同之处是限定了拆分后结果元素个数；
* `func SplitN(s, sep string, n int) []string`：拆分方式和 Split 相同，不同之处是限定了拆分后结果元素个数。

