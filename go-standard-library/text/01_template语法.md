这部分是从数据库表结构转对应 struct 结构体得来的：

~~~go
const strcutTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`

type StructTemplateDB struct {
	TableName string          // 和 structTpl 中的 .TableName 对应
	Columns   []*StructColumn // 其元素是 StructColumn，和 structTpl 中的 .Columns 对应
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}
~~~

对应得到了如下的结构：

~~~go
type COLUMNS struct {
         // TABLE_CATALOG 
         TABLECATALOG   string  `json:"TABLE_CATALOG"`
         // TABLE_SCHEMA 
         TABLESCHEMA    string  `json:"TABLE_SCHEMA"`
         // TABLE_NAME
         TABLENAME      string  `json:"TABLE_NAME"`
         // COLUMN_NAME
         COLUMNNAME     string  `json:"COLUMN_NAME"`
         // ORDINAL_POSITION
         ORDINALPOSITION        int32   `json:"ORDINAL_POSITION"`
         // COLUMN_DEFAULT
         COLUMNDEFAULT  string  `json:"COLUMN_DEFAULT"`
         // IS_NULLABLE
         ISNULLABLE     string  `json:"IS_NULLABLE"`
         // DATA_TYPE
         DATATYPE       string  `json:"DATA_TYPE"`
         // CHARACTER_MAXIMUM_LENGTH
         CHARACTERMAXIMUMLENGTH int64   `json:"CHARACTER_MAXIMUM_LENGTH"`
         // CHARACTER_OCTET_LENGTH
         CHARACTEROCTETLENGTH   int64   `json:"CHARACTER_OCTET_LENGTH"`
         // NUMERIC_PRECISION
         NUMERICPRECISION       int64   `json:"NUMERIC_PRECISION"`
         // NUMERIC_SCALE
         NUMERICSCALE   int64   `json:"NUMERIC_SCALE"`
         // DATETIME_PRECISION
         DATETIMEPRECISION      int32   `json:"DATETIME_PRECISION"`
         // CHARACTER_SET_NAME
         CHARACTERSETNAME       string  `json:"CHARACTER_SET_NAME"`
         // COLLATION_NAME
         COLLATIONNAME  string  `json:"COLLATION_NAME"`
         // COLUMN_TYPE
         COLUMNTYPE     string  `json:"COLUMN_TYPE"`
         // COLUMN_KEY
         COLUMNKEY      string  `json:"COLUMN_KEY"`
         // EXTRA
         EXTRA  string  `json:"EXTRA"`
         // PRIVILEGES
         PRIVILEGES     string  `json:"PRIVILEGES"`
         // COLUMN_COMMENT
         COLUMNCOMMENT  string  `json:"COLUMN_COMMENT"`
         // GENERATION_EXPRESSION
         GENERATIONEXPRESSION   string  `json:"GENERATION_EXPRESSION"`
         // SRS_ID
         SRSID  int32   `json:"SRS_ID"`
}

func (model COLUMNS) TableName() string {
        return "COLUMNS"
}
~~~

其中涉及到：数值填充、管道、判断、循环等语法结构。这就是使用 text/template 包下的功能：**数据驱动式**模板渲染功能，按照指定的模板输出成特定的格式。

下面对它们进行具体的讲解：

* **双层大括号**：在 template 中，所有的动作 `Actions`、数据评估、控制流转都需要用标识符双层大括号包裹，**其余**的模板内容均**全部原样输出**。
* **点**（DOT）：会根据点标识符进行**模板变量**的渲染，其参数可以为任何值，但特殊的复杂类型需进行特殊处理。例如，当为指针时，内部会在必要时自动表示为指针所指向的值。如果执行结果生成了一个函数类型的值，如结构体的函数类型字段，那么该函数不会自动调用。
* **函数调用**：在前面的代码中，通过 `FuncMap` 方法注册了名 title 的**自定义函数**。在模板渲染中一共用了两类处理方法，使用 `{{title .Name1}}` 和管道符对 `.Name3` 进行处理。在 template 中，会把管道符前面的运算结果**作为参数**传递给管道符后面的函数，最终，**命令的输出结果就是这个管道的运算结果**。

上面的 Template 模板，可以这样理解：

~~~go
const strcutTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`
~~~

其中有如下要点：

1. `{{range .Columns}}...{{end}}` 是一个循环结构；
2. `{{if gt $length 0}}...{{else}}...{{end}}` 是一个判断结构；
3. `{{.TableName | ToCamelCase}}` 使用管道，`.TableName` 作为管道的输入，让 `ToCamelCase` 作用在管道的输入上，并由此得到输出。

