> 我爸发给我的视频：
>
> * 孩子，父母的责任，是把你抚养成人，给你力所能及的关心，成长的路上，最重要的决定，其实往往需要你自己选择。
> * 孩子，走好人生路，确实不容易，但不要给自己太大压力。在外打拼，的确很辛苦，要懂得犒劳自己，要吃苦，但不要太受苦。
> * 孩子，爹娘给了你生命，希望你好好珍惜。一辈子太短暂，人生不能重来，好好珍重。家，永远是你的港湾，我们永远在等候你的归来。
> * 孩子，一个人的能力有高低，水平有大小，不要太过于攀比。树立好远大的目标和理想，朝着方向努力前进，尽最大的努力，做到问心无愧，人生就会是你想要的成功。
> * 孩子，失败并不可怕。可怕的是，你没信心和勇气去面对。别灰心丧气，一定要勇往直前，所向披靡。
> * 孩子，不要害怕在人生中遇到困难和挫折。一辈子，不如意十有八九，不会永远一帆风顺。吃得苦中苦，方为人上人。
> * 孩子，任何时候都不要放弃人生的目标。以一颗平常心，做好人生事，努力，一切都有可能。只要你不曾懈怠，你就定能走向美好的前程。
> * 孩子，为人处世的道理，你要努力学会。千万不要软弱了，人善被人欺，马善被人骑。你强大了，别人才会看得起。
> * 孩子，家家有本难念的经，人都有难言的苦，一家不知两家难。家庭，最重要的是和睦相处，与人为善。
> * 孩子，抱怨、泄气、逃避，都不是解决问题的办法。你只有抬起头来，打起精神，找准方向，才能所向披靡，一往无前。只要你勇敢拼搏奋斗，人生之路总会越走越宽广。

Go 的 bytes 包，其功能是处理 `[]byte` 类型实例的相关函数；其功能和 strings 包类似，可以做类比（analogous）。

bytes.TrimXxx 相关函数，用于在 s 中**删除指定的内容**。其中指定的内容就是入参的 cutset、prefix、suffix 等。在 doc 中说明该系列方法时，使用了一个词组 `slice off`，其中 slice 恰好就是 s 的类型，很是精妙。

> trim: to make sth neater, smaller, better, etc., by cutting parts from it; to cut away unnecessary parts from sth.
>
> to trim your hair; to tirm a hedge (back)

我们在调用 TrimXxx 相关方法时，必须**将 s 看作是 `UTF-8-encoded bytes`，也就是说 s 可能是任意的 Unicode 码**。归纳起来看，**TrimXxx 的功能就相当于是删除掉符合要求的 rune 值，并返回删除后剩下的部分。根据处理的区域（范围），以及匹配的规则，将 TrimXxx 方法分为如下 6 个部分**。

# 1 Trim

`func Trim(s []byte, cutset string) []byte`：删除 s 的头部和尾部的 cutset 集合中的 `UTF-8-encoded code points`，并返回子切片。也就是说，**只要是 cutset 集合中的 `UTF-8-encoded code points`，都需要删除**，而且删除的区域是在 s 的头部和尾部。

~~~go
func TestTrim(t *testing.T) {
	fmt.Printf("%q\n", bytes.Trim([]byte("michoim"), "mio"))
	fmt.Printf("[%q]\n", bytes.Trim([]byte(" !!! Achtung! Achtung! !!! "), "! "))
    fmt.Printf("%q\n", bytes.Trim([]byte("▓▽○中国国中○◇□▷★■●▽▒░▓"), "▽中▓"))
}
"ch"
["Achtung! Achtung"]
"○中国国中○◇□▷★■●▽▒░"
~~~

可以这样理解这个函数的结果：从首字节顺序遍历，以及从尾字节逆序遍历，如果发现 `UTF-8-encoded code points` 存在于 cutset 中，则将其删除。否则继续遍历，直到遇到了不在 cutset 中的 `UTF-8-encoded code points` 为停止遍历的信号。

需要注意的是：即然是 `UTF-8-encoded code points`，那就**支持所有的 Unicode 码**。

# 2 TrimFunc

`func TrimFunc(s []byte, f func(r rune) bool) []byte`：在 s 中的头部和尾部中，删除符合要求的 rune 值，即 func(r rune) 函数返回 true 的 rune 值。

~~~go
func TestTrimFunc(t *testing.T) {
	var trimFunc = func(r rune) bool {
		return true
	}
	val := bytes.TrimFunc([]byte("michoi"), trimFunc)
	fmt.Printf("[%q]\n", val)
}
[""]
return true --> return false
["michoi"]
~~~

因此，可以这么说，如果 trimFunc 函数直接返回 true 就表示，所有的 rune 都需要删除，反之，均保留。

那我们看看 Go 中提供能了哪些满足 `func(r rune) bool` 声明**类型**的函数，这些函数在 unicode 包下有很多：

~~~go
func IsControl(r rune) bool
func IsDigit(r rune) bool
func IsGraphic(r rune) bool
func IsLetter(r rune) bool
func IsLower(r rune) bool
func IsMark(r rune) bool
func IsNumber(r rune) bool
func IsPrint(r rune) bool
func IsPunct(r rune) bool
func IsSpace(r rune) bool
func IsSymbol(r rune) bool
func IsTitle(r rune) bool
func IsUpper(r rune) bool
~~~

作为示例：

~~~go
func TestTrimFunc(t *testing.T) {
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("\"go-gopher!\""), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}
-gopher!
"go-gopher!"
go-gopher
go-gopher!
~~~

其中 `unicode.IsPunct` 表示判断是否是标点符号（`punctuation /ˌpʌŋktʃuˈeɪʃn/ `: the marks used in writing that divide sentences and phrases; the system of using these marks）。

# 3 TrimLeft 和 TrimLeftFunc

`func TrimLeft(s []byte, cutset string) []byte`：在 Trim 函数中，处理了 s 的首部和尾部。而在 TrimLeft 函数中，则仅处理的是首部。

`func TrimLeftFunc(s []byte, f func(r rune) bool) []byte`：和 TrimFunc 函数类似，仅处理的是 s 的首部。

~~~go
func TestTrimLeft(t *testing.T) {
	fmt.Printf("%q\n", bytes.TrimLeft([]byte("michoim"), "mio"))
	fmt.Printf("[%q]\n", bytes.TrimLeft([]byte(" !!! Achtung! Achtung! !!! "), "! "))
	fmt.Printf("%q\n", bytes.TrimLeft([]byte("▓▽○中国国中○◇□▷★■●▽▒░▓"), "▽中▓"))
}
"choim"
["Achtung! Achtung! !!! "]
"○中国国中○◇□▷★■●▽▒░▓"
~~~

从和调用了 `bytes.Trim` 的结果比较来看，`bytes.TrimLeft` 确实仅删除了 s 的左侧（头部）区域的 rune。

~~~go
func TestTrimLeftFunc(t *testing.T) {
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher!"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("\"go-gopher!\""), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}
-gopher!
"go-gopher!"
go-gopher!
go-gopher!567
~~~

同理！

# 4 TrimRight 和 TrimRightFunc

`func TrimRight(s []byte, cutset string) []byte`：和 TrimLeft 恰好相反，仅删除的是 s 的尾部区域。

`func TrimRightFunc(s []byte, f func(r rune) bool) []byte`：和 TrimLeftFunc 恰好相反，仅删除的是 s 的尾部区域。

# 5 TrimPrefix 和 TrimSuffix

`func TrimPrefix(s, prefix []byte) []byte`：s 和 prefix 都是 []byte 类型的，该函数的作用是删除 s 中以 prefix 开头的部分，并返回 s 剩下的部分。

`func TrimSuffix(s, suffix []byte) []byte`：同理，返回的是删除 s 中以 suffix 的结尾部分。

~~~go
func TestTrimPrefixAndSuffix(t *testing.T) {
	fmt.Println(string(bytes.TrimPrefix([]byte("michoi"), []byte("mi"))))
	fmt.Println(string(bytes.TrimSuffix([]byte("michoi mi"), []byte(" mi"))))

	var b = []byte("Goodbye,, world!")
	b = bytes.TrimPrefix(b, []byte("Goodbye,"))
	b = bytes.TrimPrefix(b, []byte("See ya,"))
	fmt.Printf("Hello%s\n", b)
}
choi
michoi
Hello, world!
~~~

再比如：

~~~go
func TestTrimPrefixAndSuffix(t *testing.T) {
	var b = []byte("Hello, goodbye, etc!")
	b = bytes.TrimSuffix(b, []byte("goodbye, etc!"))
	b = bytes.TrimSuffix(b, []byte("gopher"))
	b = append(b, bytes.TrimSuffix([]byte("world!"), []byte("x!"))...)
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}
Hello, world!
~~~

因此，**TrimPrefix/TrimSuffix 和 TrimLeft/TrimRight 其含义和作用是有很大不同的**，TrimPrefix 和 TrimSuffix 直接匹配的是 prefix 和 suffix 的整体内容，而不是 `UTF-8-encoded code points`，相当于是**整体匹配**！

> prefix: n. a letter or group of letters added to the **beginning** of a word to change its meaning; a word, letter or number that is put before another.
>
> suffix: n. a letter or group of letters added to the **end** of a word to make another word.

# 6 TrimSpace

`func TrimSpace(s []byte) []byte`：删除 s 的头部和尾部的所有空白 `UTF-8-encoded code points`。

~~~go
func TestTrimSpace(t *testing.T) {
	values := " \t\n a lone gopher \n\t\r\n"
	fmt.Printf("%s.\n", bytes.TrimSpace([]byte(values)))

	fmt.Printf("%s.\n", bytes.TrimFunc([]byte(values), unicode.IsSpace))
}
a lone gopher.
a lone gopher.
~~~

`bytes.TrimFunc([]byte(values), unicode.IsSpace)` 是和 `bytes.TrimSpace` 相同的功能，不同之处是前者传递了一个函数实例。

在判断是否是 Space 字符时，是这样判断的：

~~~go
// IsSpace reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
// Other definitions of spacing characters are set by category
// Z and property Pattern_White_Space.
func IsSpace(r rune) bool {
	// This property isn't the same as Z; special-case it.
	if uint32(r) <= MaxLatin1 {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			return true
		}
		return false
	}
	return isExcludingLatin(White_Space, r)
}
~~~

