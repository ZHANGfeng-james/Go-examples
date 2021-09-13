Go 中的标准库 strings，用于处理 UTF-8 编码的字符串，可将 strings 包当作是一个（UTF-8编码格式的）字符串处理工具箱。



下面改变描述方式：归纳汇总所有待阐述方法，相当于是一个列表，紧接着是对其**更加细致的测试代码和说明**。



ContainsXxx 系列方法，包括：

* `func Contains(s, substr string) bool`：判断 s 中是否（完整）包含 substr；
* `func ContainsAny(s, chars string) bool`：判断 s 中是否包含任意的 chars 中的 `Unicode code points`；
* `func ContainsRune(s string, r rune) bool`：判断 r 这个特定的 `Unicode code point` 是否在 s 中。

示例程序如下：

~~~go
func TestContains(t *testing.T) {
	// 1 1
	fmt.Println(strings.Contains("seafood", "foo"))
	fmt.Println(strings.Contains("seafood", "bar"))
	// 1 0
	fmt.Println(strings.Contains("seafood", ""))
	// 0 0
	fmt.Println(strings.Contains("", ""))
	// 0 1
	fmt.Println(strings.Contains("", "foo"))
}
true
false
true
true
false
~~~

需要理解的是当 substr 和 s 分别为空字符串时的判断结果。

~~~go
func TestContainsAny(t *testing.T) {
	// 1 1
	fmt.Println(strings.ContainsAny("seafood", "foo"))
	fmt.Println(strings.ContainsAny("seafood", "bar"))
	// 1 0
	fmt.Println(strings.ContainsAny("seafood", ""))
	// 0 0
	fmt.Println(strings.ContainsAny("", ""))
	// 0 1
	fmt.Println(strings.ContainsAny("", "foo"))
}
true
true
false
false
false
~~~

和 Contains 函数的结果是有区别的，包含 2 部分：

1. ContainsAny 判断 s 中是否包含任意的 chars 中的 `Unicode code points`，而 Contains 函数则是对 substr 的整体判断；
2. ContainsAny 中边界值的判断和 Contains 是不一样的，比如当 chars 为空字符串时。

~~~go
func TestContainsRune(t *testing.T) {
    // \u706b 火 Unicode值
	fmt.Println(strings.ContainsRune("众里寻他千百度，蓦然回首，那人却在灯火阑珊处！", '\u706b'))

	// Finds whether a string contains a particular Unicode code point.
	// The code point for the lowercase letter "a", for example, is 97.
	fmt.Println(strings.ContainsRune("aardvark", 97))
	fmt.Println(strings.ContainsRune("timeout", 97))
}
true
true
false
~~~

