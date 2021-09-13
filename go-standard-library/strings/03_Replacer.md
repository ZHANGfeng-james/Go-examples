Go 中的标准库 strings，用于处理 UTF-8 编码的字符串，可将 strings 包当作是一个（UTF-8编码格式的）字符串处理工具箱。

strings.Replacer 中结构体定义：**使用 old-new 替换对字符串 s 执行替换**，这就是 Replacer 的含义

~~~go
// Replacer replaces a list of strings with replacements.
// It is safe for concurrent use by multiple goroutines.
type Replacer struct {
	once   sync.Once // guards buildOnce method
	r      replacer
	oldnew []string
}
~~~

可看到里面包含了 sync.Once，即能够确保在 goroutine 中安全使用。包含如下 3 个方法：

* `func NewReplacer(oldnew ...string) *Replacer`：根据 old-new 参数构造 `*Replacer` 实例，如果是奇数个参数，会直接 panic；
* `func (r *Replacer) Replace(s string) string`：执行替换操作，即对 s 中的字符串做 old-new 的替换；
* `func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)`：将执行 Replace 的结果写入到 io.Writer 实例中。

实例程序：

~~~go
func TestReplacer(t *testing.T) {
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt")
	fmt.Println(replacer.Replace("This is <b>HTML</b>!"))

	replacer.WriteString(os.Stdout, "This is <b>HTML</b>!\n")
}
This is &lt;b&gtHTML&lt;/b&gt!
~~~

