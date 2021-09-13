`strconv` 包所具备的功能可从 2 个方面来看：

1. 类型角度：基本类型和 string 类型之间的相互转换功能，比如**下图**所表示的含义，Basic Data Types 和 string 之间的转化；
2. 转换角度：数值转换（字符串和数值之间的转换）和字符串转换（字符串或者字符串含义的类型值之间的转换）。

![](./img/Snipaste_2021-06-24_14-02-06.png)

# 1 数值转换

最常用的数值转换方法是如下 2 个：

`func Atoi(s string) (int, error)` 将 string 转化为 int 类型数值。`Atoi` 全称是：`ascii text string to integer`

~~~go
// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Atoi(s string) (int, error) {
	const fnAtoi = "Atoi"

	sLen := len(s)
	if intSize == 32 && (0 < sLen && sLen < 10) ||
		intSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0, &NumError{fnAtoi, s0, ErrSyntax}
			}
		}

		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0, &NumError{fnAtoi, s0, ErrSyntax}
			}
			n = n*10 + int(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid, big, or underscored integers.
	i64, err := ParseInt(s, 10, 0)
	if nerr, ok := err.(*NumError); ok {
		nerr.Func = fnAtoi
	}
	return int(i64), err
}
~~~

以及，`func Itoa(i int) string` 将 int 类型数值转化为 string。`Itoa` 全称：`integer to ascii text string`，以十进制的形式进行转换。

~~~go
// FormatInt returns the string representation of i in the given base,
// for 2 <= base <= 36. The result uses the lower-case letters 'a' to 'z'
// for digit values >= 10.
func FormatInt(i int64, base int) string {
	if fastSmalls && 0 <= i && i < nSmalls && base == 10 {
		return small(int(i))
	}
	_, s := formatBits(nil, uint64(i), base, i < 0, false)
	return s
}

// Itoa is equivalent to FormatInt(int64(i), 10).
func Itoa(i int) string {
	return FormatInt(int64(i), 10)
}
~~~

在上述转化中，默认是按照十进制进行转化，而且默认转化成的是 Go 中的 int 类型数值。

与之相关的是：

~~~go
func ParseBool(str string) (bool, error)
func ParseFloat(s string, bitSize int) (float64, error)
func ParseInt(s string, base int, bitSize int) (i int64, err error)
func ParseUint(s string, base int, bitSize int) (uint64, error)

func ParseComplex(s string, bitSize int) (complex128, error)
~~~

即：将 string 转化为对应的类型，有 bool、float64、int64、uint64 和 comple128 类型。比如：

~~~go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	flag, err := strconv.ParseBool("true")
	if err != nil {
		return
	}
	// true
	fmt.Println(flag)

	value, err := strconv.ParseInt("-42", 10, 64)
	if err != nil {
		return
	}
	// int64, -42
	fmt.Printf("%T, %d \n", value, value)
    
    // biggest int32
	i64, err := strconv.ParseInt("2147483647", 10, 32)
	if err != nil {
		return
	}
	// int64, 2147483647
	fmt.Printf("%T, %d \n", i64, i64)

	i := int32(i64)
	fmt.Printf("%T, %d \n", i, i)
}
~~~

其中上述函数的参数：

* base：进制值，比如默认的十进制、八进制等；
* `bitSize`：返回值的 bit 位值，标记的是返回值的类型。比如值为 0、8、16、32、64 的入参，与之对应的是：int、int8、int16、int32、int64 类型。从函数声明来看，Parse 函数**默认地**都将 string 转化成了较宽的数值类型，比如 float64、int64 类型。如果 `biteSize` 入参表示的是较窄的类型，则其结果是可以转化为该类型的且不会损失精度，比如 int32，但方法返回值仍然是较宽的数值类型。

另外，与之相反的是：

~~~go
func FormatBool(b bool) string
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
func FormatInt(i int64, base int) string
func FormatUint(i uint64, base int) string

func FormatComplex(c complex128, fmt byte, prec, bitSize int) string
~~~

就是将 bool、float64、int64、uint64、complex128 类型转化为 string 类型值。

另外，类似的数值转化功能：

~~~go
func AppendBool(dst []byte, b bool) []byte
func AppendFloat(dst []byte, f float64, fmt byte, prec, bitSize int) []byte
func AppendInt(dst []byte, i int64, base int) []byte
~~~

将指定的类型值转化为 string，并添加到目的字符串的后面。

# 2 字符串转换

在 Go 中，我们将**带有双引号的字符串**称之为 `The quoted Go string literals`，如果要输出这种带有双引号的字符串该怎么做？

~~~go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(`This is "studygolang.com" website`)
    
    // 转义
    fmt.Println("This is \"studygolang.com\" website")
    
    fmt.Println("This is", strconv.Quote("studygolang.com"), "website")
}
~~~

与之相关的转化函数有：

~~~go
// s 转化为 a quoted Go string literals
func Quote(s string) string
// 同上，不同之处：非 ASCII 字符输出的是 Unicode 编码值
func QuoteToASCII(s string) string
// 同上，不同之处：非 ASCII 字符按照原样输出，而不是 Unicode 编码值
func QuoteToGraphic(s string) string

// r 转化为 a single-quoted Go character literal
func QuoteRune(r rune) string
func QuoteRuneToASCII(r rune) string
func QuoteRuneToGraphic(r rune) string
~~~

比如：

~~~go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	origin := "Hello, \t世界\n"
	fmt.Println(origin)

	printTypeAndValue(strconv.Quote(origin))
	printTypeAndValue(strconv.QuoteToASCII(origin))
	printTypeAndValue(strconv.QuoteToGraphic(origin))
}

func printTypeAndValue(value interface{}) {
	fmt.Printf("%T, %v\n", value, value)
}

Hello,  世界

string, "Hello, \t世界\n"
string, "Hello, \t\u4e16\u754c\n"
string, "Hello, \t世界\n"
~~~
