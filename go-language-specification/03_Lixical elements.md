`Lixical` 语法元素包括：

* 注释
* 标记（Token）分为 4 大类：标识符；关键字；操作符和标点；字面值
* 分号
* 标识符
* 关键字
* 操作符和标点（Punctuation）
* 整型字面值、浮点字面值、虚数字面值、Rune字面值（Unicode code point）、字符串字面值

# Comments

注释相当于是程序代码的文档。有 2 种形式:

* 行注释：开始于 // 字符序列，直到一行的末尾
* 标准注释：开始于 `/*` 结束于首个 `*/` 序列

注释是不会在 rune 和 string 字面值开始，比如 `var name string = "中国/*"` 并没有定义注释，而是属于 string 的一部分。不包含换行的注释，对于编译器来说相当于是空格；其他类型的注释则相当于是换行。

不能使用在注释中嵌套注释，特别是 `/**/` 这种结构中。

# Tokens

token 对于一门课编程语言来说，是一个符号。

**Tokens 构成了 Go 语言的词汇表**，包含了 4 种类型：标识符、关键字、**操作符和标点**、字面值（literal）。

由空格（U+0020）、水平制表符（U+0009）、回车符（U+000D）和换行符（U+000A）构成的空白符会被忽略，除非是这样的情况：将不同的 Token 分隔开，否则会被当做是一个 Token。

此外，换行符或文件结尾可能会插入 semicolon 分号。

将输入字符串内容拆分成 Tokens 时，下一个 Token 内容会是构成有效（Valid） Token 的最长字符序列。

# Semicolons

分号（Semicolon）

一般的编程语言使用分号 `;` 作为 Production 的结尾标志。但在 Go 语言中，如下情况可以省略分号：

* 为了让复杂的语句（Statement）写在一行，可以在 `)` 或 `}` 之前省略分号；
* 当输入被拆分成多个 Token 时，分号会被自动插入到本行的最后一个 Token 之后，如果这个 Toke 是如下类型
  * 标识符；
  * 整型、浮点型、虚数类型、rune 或字符串字面值；
  * 操作符或标点符号：++,--,),],}。
  * break、continue、fallthrough、return 关键字；

# Identifiers

标识符用于命名程序实体，比如变量和类型。更加详细的，类型可以是结构体、接口、函数、方法等。一个标识符是有一个或多个字母和数字组成，且首字符必须是字母（letter）。

下面表示的 Production 就能明确表示 Identifiers 的构成：

~~~go
identifier = letter { letter | unicode_digit } .

letter        = unicode_letter | "_" .
unicode_letter = /* a Unicode code point classified as "Letter" */ .
unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
~~~

首字符必须是 letter，由 `letter` 和 `unicode_digit` 构成，比如：

~~~go
a
_x9
ThisVariableIsExported
αβ
~~~

# Keywords

下面这些关键字是 Go 语言的**保留字**，是**不允许被当做标识符的**：

~~~go
break    defalut    func    interface    select
case    defer    go    map    struct
chan    else    goto    package    switch
const    fallthrouth    if    range    type
continue    for    import    return    var
~~~

# Operators and punctuation

接下来的字符序列表示的是**操作符和标点符号**（Punctuation）：

~~~go
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
     &^          &^=
~~~

这些字符或字符序列自然是不能用来构成标识符的。

# Literals

## Integer literals

一个整型字面值（literal）是一系列的数字，表示一个整型常量。

可选的前缀用于表示非十进制表示值：0b 或 0B 代表**二进制**；0，0o或0O 表示**八进制**；0x 或 0X 表示**十六进制**。单独数字 0 表示十进制零值。

为了阅读的便利，下划线字符 `_` 会出现在进制前导字符之后，或者出现在相邻的数字之间，这种表示方式不会改变字面值：

~~~go
int_lit        = decimal_lit | binary_lit | octal_lit | hex_lit .

decimal_lit    = "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] .
binary_lit     = "0" ( "b" | "B" ) [ "_" ] binary_digits .
octal_lit      = "0" [ "o" | "O" ] [ "_" ] octal_digits .
hex_lit        = "0" ( "x" | "X" ) [ "_" ] hex_digits .

decimal_digits = decimal_digit { [ "_" ] decimal_digit } .
binary_digits  = binary_digit { [ "_" ] binary_digit } .
octal_digits   = octal_digit { [ "_" ] octal_digit } .
hex_digits     = hex_digit { [ "_" ] hex_digit } .

decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
~~~

比如：

~~~go
42
4_2
0600
0_600
0o600
0O600       // second character is capital letter 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
170141183460469231731687303715884105727
170_141183_460469_231731_687303_715884_105727

_42         // an identifier, not an integer literal
42_         // invalid: _ must separate successive digits
4__2        // invalid: only one _ at a time
0_xBadFace  // invalid: _ must separate successive digits
~~~

## Floating-point literals

一个浮点字面值是可用**十进制**或**十六进制**表示浮点常量。

十进制表示的浮点数，其**指数部分用 `e` 或者 `E` 表示**，指数后的尾数表示的是：10^exp^。

十六进制表示的浮点数，其**指数部分用 `p` 或者 `P` 表示**，指数后的尾数表示的是：2^exp^。

为了方便阅读，下划线 `_` 字符和在整型字面值中的使用类似：

~~~go
float_lit         = decimal_float_lit | hex_float_lit .

decimal_float_lit = decimal_digits "." [ decimal_digits ] [ decimal_exponent ] |
                    decimal_digits decimal_exponent |
                    "." decimal_digits [ decimal_exponent ] .
decimal_exponent  = ( "e" | "E" ) [ "+" | "-" ] decimal_digits .

hex_float_lit     = "0" ( "x" | "X" ) hex_mantissa hex_exponent .
hex_mantissa      = [ "_" ] hex_digits "." [ hex_digits ] |
                    [ "_" ] hex_digits |
                    "." hex_digits .
hex_exponent      = ( "p" | "P" ) [ "+" | "-" ] decimal_digits .
~~~

比如：

~~~go
0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375
0x15e-2      // == 0x15e - 2 (integer subtraction)

0x.p1        // invalid: mantissa has no digits
1p-2         // invalid: p exponent requires hexadecimal mantissa
0x1.5e-2     // invalid: hexadecimal mantissa requires p exponent
1_.5         // invalid: _ must separate successive digits
1._5         // invalid: _ must separate successive digits
1.5_e1       // invalid: _ must separate successive digits
1.5e_1       // invalid: _ must separate successive digits
1.5e1_       // invalid: _ must separate successive digits
~~~

## Imaginary literals

一个虚数字面量代表的是一个复数常量。

~~~go
imaginary_lit = (decimal_digits | int_lit | float_lit) "i" .
~~~

作为示例，可以参考：

~~~go
0i
0123i         // == 123i for backward-compatibility
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i
~~~

## rune literals

> 为什么要以 UTF-8 作为 Unicode 字符的默认编码？
>
> UTF-8 编码格式是一种变长的编码格式，每个文字字符使用 1 ~ 4 个字节表示，也即是 `int32` 作为其底层存储格式是合理的，因此其存储形式紧凑而能节省内存空间。

rune 字面值表示的是 rune 类型**常量**，用来标识（identify）一个 Unicode code point 的**整型**值。rune 类型等价于 `int32` ，`rune` 是类型 `int32` 的类型别名（alias）：

~~~go
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
~~~

一个 rune 字面值使用**单引号**包括的单个或多个字符（此处并非指字符串，是需要和字符串区别的）：

~~~go
var runeChar rune = 'x'
var runeChar rune = '\n'
~~~

**在单引号中，除了换行符（直接键入 Enter 按键得到的不可见）和未转义的（`unescaped`）单引号（这种形式类似于在单引号内部直接写上单引号）外，任何字符都可以出现**。用单引号括起来的字符代表了字符本身的 Unicode 值，如果是以反斜线开头的多个字符序列，则表示的是以各种格式表示的 Unicode 值。

最简单表示 rune 字面值的形式是单引号括起来的单个字符，**因为 Go 的源文件是使用 UTF-8 编码的 Unicode 字符序列**，多个 UTF-8 编码的字节值可以用于表示单个整型值。比如字面值 'a' 使用单个字节表示的字面值 a，Unicode 值 U+0061，0x61。但是 'ä' 则使用的是 0xc3 0xa4 表示，U+00E4，0xe4。

~~~go
func main() {
	var value rune = '\u0061'
	fmt.Printf("%q, %#U.\n", value, value) // 'a', U+0061 'a'.
    
    var value = '\xff' // rune 类型
	printTypeAndValue(value) // 'ÿ', U+00FF 'ÿ'.

	var value1 = '\u00ff' // rune 类型
	printTypeAndValue(value1) // 'ÿ', U+00FF 'ÿ'.
}
~~~

综上所述：**表示 rune 有 2 种不同形式，其一是直接写出的字符，其二是对应的数值编码，表示 Unicode code point 值**。

多个反斜线转义符允许将任意值编码为 ASCII 文本。有 4 种方式将整数值表示为数字常量：\x 后跟随 2 个十六进制的数字；\u 后跟随 4 个十六进制的数字；\U 后跟随 8 个十六进制的数字；以及反斜线 \ 后跟随 3 个 8 进制的数字。上述这些情形中，字面值都是由相应基数中的数字表示的值。

虽然所有这些表示方法的结果都是一个整型值，但有不同的有效范围。八进制转义表示的范围是 [0, 255]；十六进制转义通过构造满足此条件。\u 和 \U 表示的 Unicode 码点值有可能是非法的，特别是那些大于 0x10FFFF 的值。

~~~go
0xxxxxxx                             runes 0-127    (ASCII)
110xxxxx 10xxxxxx                    128-2047       (values <128 unused)
1110xxxx 10xxxxxx 10xxxxxx           2048-65535     (values <2048 unused)
11110xxx 10xxxxxx 10xxxxxx 10xxxxxx  65536-0x10ffff (other values unused)

上述左侧表示的是底层字节数组的存储内容，右侧部分表示的是 Unicode code point 值；
对于 U+00110000(0x110000) 这个就是一个非法的 Unicode code point 值，超出了最大的可表示范围：65536-0x10ffff
~~~

比如 `\xFF` 和 `\377` 表示的是同一个 Unicode 字符：`ÿ`！

在反斜线后，下面的一些单个字符转义表示的是特殊的值：

~~~go
\a   U+0007 alert or bell
\b   U+0008 backspace
\f   U+000C form feed
\n   U+000A line feed or newline
\r   U+000D carriage return
\t   U+0009 horizontal tab
\v   U+000b vertical tab
\\   U+005c backslash
\'   U+0027 single quote  (valid escape only within rune literals)
\"   U+0022 double quote  (valid escape only within string literals)
~~~

所有其他以 `\` 反斜线（backslash）开头的字符序列存在于 rune 字面值内部的情况，都是非法的。

~~~go
rune_lit         = "'" ( unicode_value | byte_value ) "'" .
unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\` octal_digit octal_digit octal_digit .  --> 3 个八进制数值
hex_byte_value   = `\` "x" hex_digit hex_digit .  --> 2 个十六进制值
little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .  --> 4 个十六进制数值
big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
                           hex_digit hex_digit hex_digit hex_digit .  --> 8 个十六进制数值
escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) . --> 转义字符
~~~

比如：

~~~go
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
'\''         // rune literal containing single quote character
'aa'         // illegal: too many characters
'\xa'        // illegal: too few hexadecimal digits
'\0'         // illegal: too few octal digits
'\uDFFF'     // illegal: surrogate half
'\U00110000' // illegal: invalid Unicode code point
~~~

**补充 Unicode 字符的 UTF-8 编码转换原理**：

每个符号编码后第一个字节的高端 bit 位用于表示编码总共有多少个字节。如果第一个字节的高端 bit 为 0，则表示对应 7 bit 的 ASCII 字符，ASCII 字符每个字符依然是一个字节。如果第一个字节的高端 bit 是 110，则说明需要 2 个字节；后续的每个高端 bit 都以 10 开头。**更大的 Unicode 码点**也是采用类似的策略处理：

~~~go
0xxxxxxx                             runes 0-127    (ASCII)
110xxxxx 10xxxxxx                    128-2047       (values <128 unused)
1110xxxx 10xxxxxx 10xxxxxx           2048-65535     (values <2048 unused)
11110xxx 10xxxxxx 10xxxxxx 10xxxxxx  65536-0x10ffff (other values unused)
~~~

> 上述说的存储规则，其含义就是 Unicode point code 转化为底层存储值。

“麦”这个字的 Unicode 码点值（默认是 UTF-8 编码格式的）是 `\u9EA6` 的内容，底层字节数组存储的是 `e9 ba a6`；也就是说 Unicode code point 值是一个转换之后的值，是经过了 UTF-8 编码之后的值。

对于“麦”而言，可以得到如下结论：

* 过程原理和顺序：**通过底层存储的字节数组的剥离，按照 Unicode 编码的编码原则，解析出 Unicode code point 值**；
* 使用 3 个字节编码，符合 `1110xxxx 10xxxxxx 10xxxxxx` 模式；
* 根据 UTF-8 的编码规则，也就是 `1110xxxx 10xxxxxx 10xxxxxx` 规则，上述确定的位值是固定的，剩下的 `x` 是需要填充的。**把这些需要填充的 x “剥离”出来，就得到了 Unicode code point 值**；
* 对于二进制结果的第一个字节：`11101001`，剩余部分恰好是十六进制 `9` 的二进制表示结果；对于二进制结果的第二和第三个字节：`10111010`，`10100110`，将符合模板的值去除，余下结果为：`1110 1010 0110` 恰好就是和 `EA6` 的二进制相匹配了。综上所述，对于“麦”这个值 Unicode 码是 U+9EA6，对应的底层存储值是 `e9 ba a6`，其转换规则就是上述的注释内容。

再比如 `ä` 这个 Unicode 字符，底层存储的是 0xc3 0xa4 字节数组，对一个的 Unicode code point 值是什么？

* 底层字节数组：`11000011 10100100`，符合 `110xxxxx 10xxxxxx` 规则；
* 剥离出 `x` 部分的填充内容，并使用 0 填充高位：`0000 0000 1110 0100`，其结果就是 `U+00E4`。

## string literals

string 字面值表示的字符串常量，其表示形式有 2 种：

* raw string literals：原生的字符串字面量，也就是说，除了反引号本身，其他任何字符均可以，是没有转义操作的；
* interpreted string literals：经过解释后的字符串常量，潜在包含了转义。

raw string literals（可译为：原生的字符串字面量）由**反引号**（和单引号区分开）包括的字符序列，比如 `foo`。

在反引号内部中除了反引号本身，其他任何字符均可以包含在内。在原生的字符串字面值中，**没有转义操作**，全部的内容都是字面的意思。比如：

~~~go
func main() {
	var value string = `\n.
	\n.`

	fmt.Print(value)
}

PS C:\Users\Developer\sample> go run main.go
\n
        \n
~~~

特别的，**原生的字符串字面值中反斜线没有特殊含义**，并且字符串中可能包含换行符。

interpreted string literals 是使用双引号扩起来的字符序列（经过解释后的字符串字面量），比如 `"bar"`。在双引号内部，除了换行符和未转义的双引号之外（如果键入换行符时，会出现 compile 异常），任何字符都可以。双引号之间的内容构成了字面量的值，在其中包含的反斜线转义的含义和 rune 字面量的转换是一样的规则。其中有如下例外情况：

~~~go
func main() {
	var value string = "\""

	var value string = "\'" // compile error

	fmt.Print(value)
}
~~~

3 个八进制字符（`\nnn`）和 2 个十六进制字符（`0xnn`）转义表示结果字符串的各个字节。所有其他转义符表示单个字符的 UTF-8 编码（可能是多字节）：

~~~go
func main() {
	var value = '\xff' // rune 类型 
	printTypeAndValue(value) // 'ÿ', U+00FF 'ÿ'.

	var str = "\u00FF" // string 类型，或者直接写 var str = "\xc3\xbf"
	fmt.Printf("%v, length:%d.\n", str, len(str)) // ÿ, length:2.
    fmt.Println(utf8.RuneCountInString(str)) // 1 个 Rune
}
~~~

对于 `“\u00FF”` 来说，string 类型将其值转义了，也就是说将该值作为单个 UTF-8 编码解析。

但是问题在于：

~~~go
func main() {
	var value = '\xff' // rune 类型
	printTypeAndValue(value) // 'ÿ', U+00FF 'ÿ'.

	var str = "\xff" // string 类型
	fmt.Printf("%v, length:%d.\n", str, len(str)) / �, length:1.
	fmt.Println(utf8.RuneCountInString(str)) // 1
}
~~~

通过上面这个例子可以看出：在静态编程语言中，类型属性对于编译器来说是很重要的属性！对于 rune 类型来说，值为 `\xff` 需要被解析为 rune；而对于 string 类型来说，值为 `\xff` 需要被解析为单个字节，并不作为 rune 来解析。如果将 `\xff` 换成 `\u00FF` 则已经指明了作为 rune 来解析。

~~~go
string_lit             = raw_string_lit | interpreted_string_lit .
raw_string_lit         = "`" { unicode_char | newline } "`" .
interpreted_string_lit = `"` { unicode_value | byte_value } `"` .
~~~

比如：

~~~go
`abc`                // same as "abc"
`\n
\n`                  // same as "\\n\n\\n"
"\n"
"\""                 // same as `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
"\uD800"             // illegal: surrogate half
"\U00110000"         // illegal: invalid Unicode code point
~~~

下面的所有内容都表示相同的字符串字面值：

~~~go
"日本語"                                 // UTF-8 input text
`日本語`                                 // UTF-8 input text as a raw literal
"\u65e5\u672c\u8a9e"                    // the explicit Unicode code points
"\U000065e5\U0000672c\U00008a9e"        // the explicit Unicode code points
"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // the explicit UTF-8 bytes
~~~
