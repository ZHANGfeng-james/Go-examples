==结构体指针==的问题：

~~~go
func main() {
	stuPtr := &Stu{"abc", 123}
	fmt.Println(*stuPtr)
	fmt.Printf("%p\n", stuPtr) // 0xc0000044c0

	stuPtr = &Stu{"abc", 123}
	fmt.Printf("%p\n", stuPtr) // 0xc000004520
}
~~~

`stuPtr` 类型是 `*Stu` 类型，通过输出格式 `fmt.Printf("%p\n", stuPtr)` 输出的是该指针指向的内存地址。上述代码可得到的结论是：`stuPtr` 指向的是不同的 `Stu` 变量，也就是不同的内存地址。

==结构体变量的零值==问题：

~~~go
type Stu struct {
	name *string
	age  int
}

func main() {
	var s1 Stu

	fmt.Println(s1)

	var s2 *string
	fmt.Printf("%p, %v\n", s2, s2)  // 0x0, <nil>
}
~~~

上面程序中声明了 s1 变量，其成员初始化为其类型的零值，也就是说对于 string 类型初始化为空字符串，对于 age 而言初始化为 0；而对于指针类型，则相当于是没有指向任何内存地址，**也就是 0x0（地址值）**，其值为 `nil`。

==结构体变量的赋值==问题：

~~~go
type Stu struct {
	name string
	age  int
}

func main() {
	var s1 = Stu{name: "z", age: 10}
	fmt.Println(s1)

	s2 := s1
	s2.age = 20
	fmt.Println(s1)
}
~~~

定义结构体（并不是引用类型）变量后，重新定义了另外的结构体变量并赋值为 s1。此时对 s2 进行修改，并再次打印 s1，得到了相同的结果。由此说明了：**结构体变量的赋值，仅仅做的是结构体中成员的（值）拷贝**。但是如果 `Stu` 结构体中**==成员是指针类型，情况就会不一样==**：

~~~go
type Stu struct {
	name *string  // 想在将 name 声明为 *string 类型
	age  int
}

func main() {
	var s1 Stu
	name := "z"

	s1.name = &name
	s1.age = 10
	fmt.Println(*(s1.name))

	s2 := s1   // 同样是做了值拷贝，s2.name 和 s1.name 指向了相同内存地址
	*(s2.name) = "a"
	s2.age = 20
	fmt.Println(*(s1.name))
}
~~~

由此可以看出：如果结构体成员是指针类型时，结构变量的赋值操作会让指针指向相同的内存地址，需要小心修改对全局造成的影响。

==复杂结构体内容的理解==：

~~~go
type Stu struct {
	addr *Stu
	buf  []byte
}

func main() {
	var stu Stu

	// 0xc0000044c0（变量的地址）, 0x0（变量的地址）, <nil>
	fmt.Printf("%p, %p, %v\n", &stu, stu.addr, stu.addr)
	
	// 0x0, []（stu.buf变量的内容）
	fmt.Printf("%p, %v\n", stu.buf, stu.buf)

	// 凡是在 Printf 中指定为 %p 格式输出时，得到的是变量的地址（指针值）；如果不是指针，则取其指针，必须要让类型一致！
}
~~~

