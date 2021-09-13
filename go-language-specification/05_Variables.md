# Variables

在 Go 中，**一个变量是一个用于存放值的内存（存储）空间**。一个变量，不仅包含了**类型信息**，还包括了其中**要存储的实际值**。变量允许（permissible）存储的值的集合取决于变量的类型。

一个变量声明，或者是函数参数或结果，或者一个函数声明的签名，或者函数字面值都为命名变量保留存储空间。内置函数 new 的调用，或者是持有复合（composite）字面值都会在运行期**为变量分配内存空间**。隐式（anonymous）变量会通过一个指针被间接持有。

结构化类型的变量，比如 array、slice 和 struct 具有可以单独（individually）寻址的元素和字段，其中每个这样的元素都像一个变量。

**变量的静态类型**在声明语句中给出，或者是调用 new 方法时给出，或者是复合字面值，或者是结构化变量的元素类型。接口类型变量具有一个具体的（concrete）动态类型，实际就是在运行期被赋予的值的类型，**除非这个值是预先定义的 nil，它没有类型**。动态类型在运行期会发生改变，但是存储在接口变量中的值始终是可被赋值给该变量的静态类型的值。

~~~go
type T = string

func main() {
	var x interface{}
	fmt.Printf("%T, %v.\n", x, x) // <nil>, <nil>.
	var v *T
	fmt.Printf("%T, %v.\n", v, v)
	if v == nil {
		fmt.Println("v is nil!")
	}

	x = 42
	printTypeAndValue(x) // int, 42.
	x = v
	printTypeAndValue(x) // *string, <nil>.
	if x == nil {
		fmt.Println("x is nil!")
	}

	str := "Katyusha"
	v = &str
	x = v
	printTypeAndValue(x) // *string, 0xc0000321f0.
}

func printTypeAndValue(value interface{}) {
	fmt.Printf("%T, %v.\n", value, value)

	if nil == value {
		fmt.Println("value is nil!")
	}
}

*string, <nil>.
v is nil!
int, 42.
*string, <nil>.
*string, 0xc0000321f0.
~~~

下面一行行代码予以解释：

* `var x interface{}`：x 值为 nil，具备的静态类型是 `interface{}`
* `var v *T`：v 值是 nil，静态类型是 `*string`
* `x = 42`：x 变量值是 42，动态类型是 `int`
* `x = v`：x 值是 `(*T)(nil)`，动态类型是 `*string`；虽然 v 的值是 nil，但是 x 的值并不是 nil

通过引用表达式中的变量来检索变量的值，它是分配给变量的最新值。如果一个变量没有被赋予值，那么这个变量就会被赋予这个类型的零值（zero value）。

