Notation 其含义是 Go 的**符号**系统。

Go 的语法（syntax）使用的是 `Extended Backus-Naur Form` （EBNF） 扩展巴科斯-瑙尔范式。

> EBNF 是一种用于描述计算机编程语言的、与上下文无关语法的**元语法符号表示法**，也就是一种描述语言的语言。

EBNF 的基本语法形式叫做是 Production：

~~~go
Production  = production_name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .

生成式 = 生成式名 "=" [ 表达式 ] ".".
表达式 = 选择项 { "|" 选择项 } .
选择项 = 条目 { 条目 } .
条目   = 生成式名 | 标记 [ "…" 标记 ] | 分组 | 可选项 | 重复项 .
分组   = "(" 表达式 ")" .
可选项 = "[" 表达式 "]" .
重复项 = "{" 表达式 "}" .
~~~

生成式由表达式构成，表达式通过条目及以下操作符构成，其优先级由低到高：

~~~go
|   选择项并联
()  分组
[]  选项
{}  重复
~~~

`a…b` 这种结构代表的是从 a 到 b 的字符集合作为可选项（Alternative）。在 Go 语言规范中的其他地方也使用 `…` 来非正式地表示各种枚举或未进一步明确的代码片段。`…` 区别于 `...`！

`Option` 和 `Alternative` 的区别是什么？