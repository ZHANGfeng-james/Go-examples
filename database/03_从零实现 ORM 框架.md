# 1 ORM 框架的疑惑

关于 ORM 的疑惑：

1. ORM 框架是什么？
2. ORM 框架解决了什么问题？在没有 ORM 框架时，开发人员面临什么问题？
3. ORM 框架是如何解决这些问题的？

对象关系映射（Object Relational Mapping）是通过使用描述对象和数据库之间**映射**的**元数据**，将面向对象语言程序中的对象自动**持久化**到关系数据库中。那首先第一个问题：**面向对象程序设计中的对象和数据库是如何映射的？**

| 数据库 |      面向对象程序设计中的对象      |
| :----: | :--------------------------------: |
|   表   |     类（类型，struct、class）      |
|  记录  |          对象（类的实例）          |
|  字段  | 对象的属性（类型实例中各个域的值） |

紧接着的问题是：**这个映射关系是如何实现的**？

举一个例子：

~~~sql
CREATE TABLE `User` (`Name` text, `Age` integer);
INSERT INTO `User` (`Name`, `Age`) VALUES ("Tom", 18);
SELECT * FROM `User`;
~~~

假如使用 ORM 框架，可以这样写：

~~~go
type User struct {
    Name string
    Age int
}

orm.CreateTable(&User{})
orm.Save(&User{"Tom", 18})

var users []User
orm.Find(&users)
~~~

ORM 框架相当于**对象和数据库中间的一个桥梁**，借助 ORM 可以**避免写繁琐的 SQL 语句**，仅仅通过操作具体的对象，就能够完成对关系型数据库的操作。那**实现一个 ORM 框架需要考虑哪些问题**：

* `CreateTable` 方法需要从参数 `&User{}` 得到对应的结构体的名称 User 作为**表名**，成员变量 Name、Age 作为**列名**，同时还需要知道成员变量对应的**类型**。上述所有的信息，都是创建一张 Table 必须要的！
* `Save` 方法则需要知道每个成员变量的值。
* `Find` 方法仅从传入的空切片 `&[]User` 得到对应的结构体名也就是**表名** User，并从数据库中取到所有的记录，将其**转换**成 User 对象（即：数据库查询的结构转换成对象），添加到切片中。

另外，ORM 框架是通用的，也就是说可以将**任意合法的对象**转换成数据库中的表和记录。应用在工程中，可能会使用多个类型，比如 User、Account、Passenger、Password 等等，都是不一样的类型，对应的就会有很多张不同的表。

~~~go
type Account struct {
    Username string
    Password string
}
~~~

这就带来了一个很重要的问题：**如何根据任意类型的指针，得到其对应的结构体的信息**。这涉及到了 Go 语言的反射机制(reflect)，通过反射，可以获取到对象对应的结构体名称，成员变量、方法等信息，例如：

~~~go
typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type() // 为什么要如此使用？
fmt.Println(typ.Name)
for i :=0 i < typ.NumField(); i++ {
    field := typ.Field(i)
    fmt.Println(field.Name)
}
~~~

* `reflect.ValueOf` 获取指针对应的反射值。**一个指针变量的值，获取其反射值，其结果是什么含义**？
* `reflect.Indirect()` 获取指针指向的对象的反射值。**反射值的含义是什么**？
* `(reflect.Type).Name()` 依据类型信息返回类名；
* `(reflect.Type).Field(i)` 依据类型信息获取第 i 个成员变量。

~~~go
package main

import (
	"fmt"
	"reflect"
)

type Account struct {
	username string
	age      int8
}

func main() {
	ptr := &Account{} // 如果是指针变量
	typ := reflect.Indirect(reflect.ValueOf(ptr)).Type()
	fmt.Println(typ.Name())

	obj := Account{} // 如果是普通变量（非指针）
	value := reflect.ValueOf(obj)
	typ = value.Type()
	fmt.Println(typ.Name())
}
~~~

除了对象和表结构/记录的映射以外，设计 ORM 框架还需要关注什么问题呢？

1）MySQL，PostgreSQL，SQLite 等数据库的 **SQL 语句是有区别的**，ORM 框架如何在开发者不感知的情况下**适配多种数据库**？

2）如何对象的字段发生改变，数据库**表结构**能够**自动更新**，即是否支持**数据库自动迁移(migrate)**？

3）**数据库支持的功能**很多，例如事务(transaction)，ORM 框架能实现哪些？

4）...

> 一个数据库到底能够支持多少功能？这些功能都在哪里可以找到“线索”？关于数据库，在应用层面上能够使用的能力有哪些？数据库管理系统的能力是如何在一个后端应用程序上体现出来的？

**数据库的特性**非常多，简单的增删查改使用 ORM 替代 SQL 语句是没有问题的，但是也有很多特性难以用 ORM 替代，比如**复杂的多表关联查询**，ORM 也可能支持，但是基于性能的考虑，开发者自己写 SQL 语句很可能更高效。因此，设计实现一个 ORM 框架，就需要给功能特性排优先级了。

Go 语言中使用比较广泛 ORM 框架是 [gorm](https://github.com/jinzhu/gorm) 和 [xorm](https://github.com/go-xorm/xorm)。除了基础的功能，比如表的操作，记录的增删查改，gorm 还实现了关联关系(一对一、一对多等)，回调插件等；xorm 实现了**读写分离(支持配置多个数据库)**，**数据同步**，**导入导出**等。gorm 正在彻底重构 v1 版本，短期内看不到发布 v2 的可能。相比于 gorm-v1，xorm 在设计上更清晰。

接下来，我就参考 GeeORM 去实现一个 ORM 框架，支持的特性有：

* 表的创建、删除、迁移；
* 记录的 CRUD，查询条件的链式操作；
* 单一主键的设置；
* 钩子（在创建/更新/删除/查找之前或之后）；
* 事务等。

这一个个独立的特性组合在一起就是最终的 ORM 框架！

> 弄懂这些专业术语，能帮助理解整个系统的运行。

# 2 database/sql 基础

> SQLite is a C-language library that implements a small, fast, self-contained, high-reliability, full-featured, SQL database engine.
>
> 我作为一个初学者，最感兴趣的是：`full-featured` 这个形容词，也就是 SQLite 包含 RDBMS 关系型数据库管理系统的所有特征。那**一个 RDBMS 有哪些特征**？
>
> SQLite 可以**直接嵌入到代码中**，不需要像 MySQL、PostgreSQL 需要**启动独立的服务**才能使用。

Mac 系统上默认安装了 `SQLite`：

~~~shell
ant@MacBook-Pro ~ % sqlite3
SQLite version 3.32.3 2020-06-18 14:16:19
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> .help
.auth ON|OFF             Show authorizer callbacks
.backup ?DB? FILE        Backup DB (default "main") to FILE
.bail on|off             Stop after hitting an error.  Default OFF
.binary on|off           Turn binary output on or off.  Default OFF
.cd DIRECTORY            Change the working directory to DIRECTORY
.changes on|off          Show number of rows changed by SQL
.check GLOB              Fail if output since .testcase does not match
.clone NEWDB             Clone data into NEWDB from the existing database
.databases               List names and files of attached databases
.dbconfig ?op? ?val?     List or change sqlite3_db_config() options
.dbinfo ?DB?             Show status information about the database
.dump ?TABLE?            Render database content as SQL
.echo on|off             Turn command echo on or off
.eqp on|off|full|...     Enable or disable automatic EXPLAIN QUERY PLAN
.excel                   Display the output of next command in spreadsheet
.exit ?CODE?             Exit this program with return-code CODE
.expert                  EXPERIMENTAL. Suggest indexes for queries
.explain ?on|off|auto?   Change the EXPLAIN formatting mode.  Default: auto
.filectrl CMD ...        Run various sqlite3_file_control() operations
.fullschema ?--indent?   Show schema and the content of sqlite_stat tables
.headers on|off          Turn display of headers on or off
.help ?-all? ?PATTERN?   Show help text for PATTERN
.import FILE TABLE       Import data from FILE into TABLE
.imposter INDEX TABLE    Create imposter table TABLE on index INDEX
.indexes ?TABLE?         Show names of indexes
.limit ?LIMIT? ?VAL?     Display or change the value of an SQLITE_LIMIT
.lint OPTIONS            Report potential schema issues.
.log FILE|off            Turn logging on or off.  FILE can be stderr/stdout
.mode MODE ?TABLE?       Set output mode
.nullvalue STRING        Use STRING in place of NULL values
.once ?OPTIONS? ?FILE?   Output for the next SQL command only to FILE
.open ?OPTIONS? ?FILE?   Close existing database and reopen FILE
.output ?FILE?           Send output to FILE or stdout if FILE is omitted
.parameter CMD ...       Manage SQL parameter bindings
.print STRING...         Print literal STRING
.progress N              Invoke progress handler after every N opcodes
.prompt MAIN CONTINUE    Replace the standard prompts
.quit                    Exit this program
.read FILE               Read input from FILE
.recover                 Recover as much data as possible from corrupt db.
.restore ?DB? FILE       Restore content of DB (default "main") from FILE
.save FILE               Write in-memory database into FILE
.scanstats on|off        Turn sqlite3_stmt_scanstatus() metrics on or off
.schema ?PATTERN?        Show the CREATE statements matching PATTERN
.selftest ?OPTIONS?      Run tests defined in the SELFTEST table
.separator COL ?ROW?     Change the column and row separators
.session ?NAME? CMD ...  Create or control sessions
.sha3sum ...             Compute a SHA3 hash of database content
.shell CMD ARGS...       Run CMD ARGS... in a system shell
.show                    Show the current values for various settings
.stats ?on|off?          Show stats or turn stats on or off
.system CMD ARGS...      Run CMD ARGS... in a system shell
.tables ?TABLE?          List names of tables matching LIKE pattern TABLE
.testcase NAME           Begin redirecting output to 'testcase-out.txt'
.testctrl CMD ...        Run various sqlite3_test_control() operations
.timeout MS              Try opening locked tables for MS milliseconds
.timer on|off            Turn SQL timer on or off
.trace ?OPTIONS?         Output each SQL statement as it is run
.vfsinfo ?AUX?           Information about the top-level VFS
.vfslist                 List all available VFSes
.vfsname ?AUX?           Print the name of the VFS stack
.width NUM1 NUM2 ...     Set column widths for "column" mode
~~~

熟悉实验环境：看看在 SQLite 上都能够做什么？

~~~sqlite
sqlite> INSERT INTO User(Name, Age) VALUES ("Tom", 18),("Jack", 25);
sqlite> .head on
sqlite> SELECT * FROM User;
Name|Age
Tom|18
Jack|25
sqlite> SELECT * FROM User WHERE Age > 20;
Name|Age
Jack|25
sqlite> SELECT COUNT(*) FROM User;
COUNT(*)
2
~~~

有了实验的环境，再来看看 Go 语言提供的和 DB 相关的**标准库函数和类型**：

* database/sql 中的**函数**：

~~~shell
ant@MacBook-Pro ~ % go doc database/sql |grep "^func"
func Drivers() []string
func Register(name string, driver driver.Driver)
~~~

* database/sql 中的**数据类型**：

~~~shell
ant@MacBook-Pro ~ % go doc database/sql |grep "^type"|grep struct
type ColumnType struct{ ... }
type Conn struct{ ... }
type DB struct{ ... }
type DBStats struct{ ... }
type NamedArg struct{ ... }
type NullBool struct{ ... }
type NullFloat64 struct{ ... }
type NullInt32 struct{ ... }
type NullInt64 struct{ ... }
type NullString struct{ ... }
type NullTime struct{ ... }
type Out struct{ ... }
type Row struct{ ... }
type Rows struct{ ... }
type Stmt struct{ ... }
type Tx struct{ ... }
type TxOptions struct{ ... }
~~~

下面来看一个简单的实例：

~~~go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string
	Age  int8
}

func (user *User) String() string {
	return fmt.Sprintf("User: Name=%s, Age=%d\n", user.Name, user.Age)
}

func main() {
	db, err := sql.Open("sqlite3", "../gee.db") // driver 和 数据库名
	if err != nil {
		fmt.Println("err!")
		log.Fatal(err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	row, err := db.Query("SELECT * FROM User") // 查询数据库记录
	if err != nil {
		log.Fatal(err.Error())
	}

	var name string
	var age int8
	for row.Next() {
		if err = row.Scan(&name, &age); err == nil { // 遍历每条记录，获取记录值
			fmt.Printf("User: Name=%s, Age=%d\n", name, age)
		} else {
			fmt.Println(err.Error())
		}
	}
}
~~~

对上述程序使用到的 API 说明：

- 使用 `sql.Open()` 连接数据库，第一个参数是驱动名称，import 语句 `_ "github.com/mattn/go-sqlite3"` 包导入时会**注册 sqlite3 的驱动**，第二个参数是**数据库的名称**，对于 SQLite 来说，也就是**文件名**，不存在会新建。返回一个 `sql.DB` 实例的指针。
- `Exec()` 用于**执行 SQL 语句**，如果是查询语句，不会返回相关的记录。所以查询语句通常使用 `Query()` 和 `QueryRow()`，前者可以**返回多条记录**，后者**只返回一条记录**。
- `Exec()`、`Query()`、`QueryRow()` 接受1或多个入参，第一个入参是 SQL 语句，后面的入参是 SQL 语句中的占位符 `?` 对应的值，占位符一般用来防 SQL 注入。
- `QueryRow()` 的返回值类型是 `*sql.Row`，`row.Scan()` 接受1或多个指针作为参数，可以**获取对应列(column)的值**，在这个示例中，有 `Name` 和 `Age` 两列，因此传入字符串指针 `&name` 和 `&age` 即可获取到查询的结果。

开发一个框架/库并不容易，详细的日志能够**帮助我们快速地定位问题**。由此，接下来先写一个简单的 log 库：为什么不直接使用原生的 log 库呢？log 标准库没有日志分级，不打印文件和行号，这就意味着我们很难快速知道是哪个地方发生了错误。这个简易的 log 库具备以下特性：

- 支持日志分级（Info、Error、Disabled 三级）。
- 不同层级日志显示时使用不同的颜色区分。
- 显示打印日志代码对应的文件名和行号。

完成后的日志功能：

~~~go
package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log method，导出 Logger 上的方法
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disable
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
~~~

- 这一部分的实现非常简单，三个层级声明为三个常量，通过控制 `Output`，来控制日志是否打印。
- 如果设置为 ErrorLevel，infoLog 的输出会被定向到 `ioutil.Discard`，即不打印该日志。

如果使用**层级**思维来考虑 ORM 库的实现，最底层的应该是**直接与 RDBMS 的交互**，也就是 CRUD 操作的执行：

~~~go
package session

import (
	"database/sql"
	"strings"

	"github.com/go-examples-with-tests/database/v1/log"
)

type Session struct {
	db      *sql.DB         // 数据库实例，用于和数据库交互，执行 CRUD 操作
	sql     strings.Builder // SQL 语句
	sqlVars []interface{}   // SQL 语句中的 ? 占位符对应的参数
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec execs a SQL statement, and return sq.Result
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
    log.Info(s.sql.String(), s.sqlVars)
	// 调用的是 sql.DB 的 QueryRow 函数，仅返回一行结果
	return s.DB().QueryRow(s.sql.String(), s.sqlVars)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
    log.Info(s.sql.String(), s.sqlVars)
	// 调用的是 sql.DB 的 Query 函数，可返回多行结果
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
~~~

另外封装了 `sql.DB` 的方法：`Exec`、`QueryRow` 和 `Query`，其封装的目的在于：

* 统一打印，能够快速定位错误；
* 复用 Session，创建一个 Session 实例，可以多次和 RDBMS 交互。

Session 负责与数据库管理系统交互，那交互前的初始化工作交给 Engine 处理：

~~~go
package orm

import (
	"database/sql"

	"github.com/go-examples-with-tests/database/v1/log"
	"github.com/go-examples-with-tests/database/v1/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}
~~~

基本功能就是：

* 连接数据库，返回 `*sql.DB` 实例；
* 使用 Ping 方法，测试连接的可用性。

上述就是整个 ORM 框架的雏形了，接下来就是框架的使用：

~~~go
func main() {
	engine, err := orm.NewEngine("sqlite3", "../gee.db")
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	session := engine.NewSession()
	session.Raw("DROP TABLE IF EXISTS User;").Exec()

	session.Raw("CREATE TABLE User(Name text)").Exec()
	session.Raw("CREATE TABLE User(Name text)").Exec()

	result, _ := session.Raw("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
~~~

对应日志输出：

~~~shell
ant@MacBook-Pro v1 % go run base.go
[info ]2021/09/28 17:24:33 orm.go:26: Connect database success
[info ]2021/09/28 17:24:33 raw.go:38: DROP TABLE IF EXISTS User;  []
[info ]2021/09/28 17:24:33 raw.go:38: CREATE TABLE User(Name text)  []
[info ]2021/09/28 17:24:33 raw.go:38: CREATE TABLE User(Name text)  []
[error]2021/09/28 17:24:33 raw.go:40: table User already exists
[info ]2021/09/28 17:24:33 raw.go:38: INSERT INTO User(`Name`) VALUES (?), (?)  [Tom Sam]
Exec success, 2 affected
[info ]2021/09/28 17:24:33 orm.go:34: Close database success
~~~

# 对象表结构映射

对象表结构映射，解决的问题就是：**程序中的一个结构体指针（结构体）变量，转化为数据库中的一张表**。与之相关的，就是要获取到这个结构体指针（结构体）变量中各个字段的名称、类型和其他能够转化为表结构的约束信息。

SQL 语句中的**类型**和 Go 语言中的**类型**是不同的，例如Go 语言中的 `int`、`int8`、`int16` 等类型均对应 SQLite 中的 `integer` 类型。因此实现 ORM 映射的第一步，需要思考如何将 Go 语言的类型映射为数据库中的类型。

> 如果我是 RDBMS 的设计者，我也不会去**跟随**具体的编程语言。每一种编程语言的类型都不相同，而它们都需要和 RDBMS 交互，那是不是要为每一种编程语言在 RDBMS 中设计一种对应的类型呢？
>
> 与其如此，何不就在 RDBMS 中自成一体设计出一套自己的类型系统！

同时，不同数据库支持的数据类型也是有差异的，即使功能相同，在 **SQL 语句的表达**上也可能有差异。ORM 框架往往需要兼容多种数据库，因此我们需要将**差异**的这一部分提取出来，**每一种数据库分别实现，实现最大程度的复用和解耦**。这部分代码称之为 `dialect`。

> `dialect`：a particular form of a language which is peculiar to a specific region or social group.
>
> 方言

下面实现各种 RDBMS 的差异部分：

~~~go
package dialect

import (
	"fmt"
	"reflect"
)

var dialectsMap = map[string]Dialect{} // 进程全局保存注册的 name - Dialect

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                        // Go-type convert to RDMS-type
	TableExistSQLStmt(tableName string) (string, []interface{}) // 指定tablename是否存在的SQL语句
}

func RegisterDialect(name string, dialect Dialect) {
	_, ok := GetDialect(name)
	if ok {
		panic(fmt.Sprintf("dialect for %s just registe once", name))
	}
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
~~~

每一种 RDBMS 都对应有一种 Dialect，也就是差异的部分。比如对于 sqlite3 这种 RDBMS 对应的实现就是：

~~~go
package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

func init() {
	//FIXME 此处可能会导致空指针异常，var _ Dialect = (*sqlite3)(nil) 强制初始化 dialog.go
	RegisterDialect("sqlite3", &sqlite3{})
}

// DataTypeOf convert Go-type to RDMS-type
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool" // type of RDBMS
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice: // 使用实例？看看别人是怎么使用的
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (s *sqlite3) TableExistSQLStmt(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
~~~

对于 sqlite3 这种 dialect，实现了 Dialect 中的接口，也就是与其他关系型数据库管理系统的差异。不同数据库之间的差异远远不止这两个地方，随着 ORM 框架功能的增多，dialect 的实现也会逐渐丰富起来，同时框架的其他部分不会受到影响。

如下实现对 sqlite3.go 的测试：

~~~go
package dialect

import (
	"reflect"
	"testing"
)

func TestDataTypeOf(t *testing.T) {
	p := []struct {
		Values interface{}
		Type   string
	}{
		{"Tom", "text"},
		{123, "integer"},
		{1.23, "real"},
		{[]int{1, 2, 3, 4}, "blob"},
	}
	// 在执行 TestDataTypeOf 时，已经调用了sqlite3.go中的init()
	sqlDB, ok := GetDialect("sqlite3")
	if ok {
		for _, parameter := range p {
			if typ := sqlDB.DataTypeOf(reflect.ValueOf(parameter.Values)); typ != parameter.Type {
				t.Fatalf("Type of %v is %s, got:%s", parameter.Values, parameter.Type, typ)
			}
		}
	}
}
~~~

Dialect 实现了一些特定的 SQL 语句的转换，接下来我们将要实现 ORM 框架中最为核心的转换——对象(object)和表(table)的转换。给定一个**任意的对象**，转换为关系型数据库中的**表结构**。在数据库中创建一张表需要哪些**要素**呢？

- **表名**(table name) —— 结构体名(struct name)
- **字段名和字段类型** —— 成员变量和类型。
- **额外的约束条件**(例如非空、主键等) —— 成员变量的Tag（Go 语言通过 Tag 实现，Java、Python 等语言通过注解实现）

比如：Go 程序中定义的一个结构体类型

~~~go
type User struct {
    Name string `geeorm:"PRIMARY KEY"`
    Age  int
}
~~~

对应转化成数据库管理系统中的建表语句：

~~~sql
CREATE TABLE `User` (`Name` text PRIMARY KEY, `Age` integer);
~~~

也就是需要从一个结构体中解析上述要素：

~~~go
package schema

import (
	"go/ast"
	"reflect"

	"github.com/go-examples-with-tests/database/v2/dialect"
)

// 一张 Table 中，Column 相关的信息
type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}       // 值，一般是指针类型的值
	Name       string            // 类型名，指针类型的值中解析出类型名，作为表名
	Fields     []*Field          // 表相关的所有列信息
	FieldNames []string          // 表相关的所有列名（字段名）
	fieldMap   map[string]*Field // 列名（字段名） - 列信息
}

type ITableName interface {
	TableName() string
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 依据具体的 dialect.Dialect 作类型转换
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	var tableName string
	t, ok := dest.(ITableName) // 是否实现ITableName接口，可自定义表名
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i) // StructField 类型
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				// reflect.Indirect(reflect.New(p.Type)) --> 创建指针类型实例，并访问
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest)) // reflect.Value
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		// reflect.Value struct --> value
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
~~~

整个解析的过程使用的原理：Go reflect 机制。

~~~go
package schema

import (
	"testing"

	"github.com/go-examples-with-tests/database/v2/dialect"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"` // struct的TAG有固定的格式，写错则无效！
	Age  int
}

func TestSchema(t *testing.T) {
	dialect, _ := dialect.GetDialect("sqlite3")

	user := &User{}
	userSchema := Parse(user, dialect)
	if userSchema.Name != "User" && len(userSchema.Fields) != 2 {
		t.Fatal("schema parse User error")
	}
	if userSchema.fieldMap["Name"].Tag != "PRIMARY KEY" {
		t.Fatal("schema parse User error")
	}
}

func TestRecordValue(t *testing.T) {
	user := &User{
		Name: "Tom",
		Age:  18,
	}

	dialect, _ := dialect.GetDialect("sqlite3")

	schema := Parse(user, dialect)
	values := schema.RecordValues(user)

	name := values[0].(string)
	age := values[1].(int)
	if name != "Tom" && age != 18 {
		t.Fatal("record value is error")
	}
}

type Password struct {
	Len     int
	Content string
}

func (p *Password) TableName() string {
	return "test_password_name"
}

func TestSchemaPassword(t *testing.T) {
	dialect, _ := dialect.GetDialect("sqlite3")

	password := &Password{}
	passwordSchema := Parse(password, dialect)
	if passwordSchema.Name != password.TableName() && len(passwordSchema.Fields) != 2 {
		t.Fatal("schema parse Password error")
	}
}
~~~

Session 的核心功能是与数据库进行交互，因此，我们将**数据库表的增/删操作**实现在子包 session 中。在此之前，Session 的结构需要做一些调整。

~~~go
type Session struct {
	db      *sql.DB         // 数据库实例，用于和数据库交互，执行 CRUD 操作
	sql     strings.Builder // SQL 语句
	sqlVars []interface{}   // SQL 语句中的 ? 占位符对应的参数

	dialect  dialect.Dialect
	refTable *schema.Schema
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}
~~~

新增方法：

~~~go
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) { // 指针值
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}
~~~

实现数据库表的创建、删除和判断是否存在的功能。三个方法的实现逻辑是相似的，利用 `RefTable()` 返回的数据库表和字段的信息，拼接出 SQL 语句，调用原生 SQL 接口执行。

~~~go
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.refTable.Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQLStmt(s.refTable.Name)
	log.Infof("HasTable: %s, args:%v", sql, values)
	row := s.Raw(sql, values...).QueryRow()

	var tmp string
	_ = row.Scan(&tmp)
	log.Infof("Query:%s, Got:%s", s.refTable.Name, tmp)
	return tmp == s.refTable.Name
}
~~~

对应的测试用例：

~~~go
package session

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-examples-with-tests/database/v2/dialect"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
}

var (
	TestDB         *sql.DB
	TestDialect, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) { // 这个方法在所有测试用例之前执行
	fmt.Println("Main")
	TestDB, _ = sql.Open("sqlite3", "../../gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func TestSession(t *testing.T) {
	session := New(TestDB, TestDialect)
	session.Model(&User{})

	session.DropTable()
	session.CreateTable()
	if !session.HasTable() {
		t.Fatal("create table error!")
	}
}

func TestModel(t *testing.T) {
	session := New(TestDB, TestDialect)
	session.Model(&User{})
	table := session.refTable

	session.Model(&Session{})

	if table.Name != "User" || session.refTable.Name != "Session" {
		t.Fatal("failed to change model")
	}
}
~~~

Session 增加了对 dialect.Dialect 的依赖，调整 Engine：

~~~go
package orm

import (
	"database/sql"

	"github.com/go-examples-with-tests/database/v2/dialect"
	"github.com/go-examples-with-tests/database/v2/log"
	"github.com/go-examples-with-tests/database/v2/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("get dialect: %s error", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
~~~

至此，完成了对象表结构映射的目标：

1. 为隔离不同数据库管理系统的差异，创建了 Dialect，实现 sqlite3 的实例；
2. 使用 reflect 实现了从结构体类型中获取字段名、类型、Tag，并转化为 Table；
3. 实现创建 Table、删除 Table和判断 Table 是否存在的操作。

# 记录新增和查询

查询语句一般由多个**字句（Clause）**构成：

~~~sql
SELECT col1, col2, ...
	FROM table_name
	WHERE [condition]
	GROUP BY col1
	HAVING [condition];
~~~

也就是说，如果想一次构造出完整的 SQL 语句是比较困难的，由此将构造 SQL 语句这一部分独立处理啊，放在新创建的子包 clause 中实现。也就是说，**实现各种子句的生成规则**：

~~~go
package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)

	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	sql.WriteString("VALUES ")

	var bindStr string
	var vars []interface{}

    // 构造成这样的形式：VALUES (?, ?), (?, ?), (?, ?)  [Katyusha 31 Sam 32 Jason 33]
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}

	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	// values[1] 是什么类型？按照 []string 类型转换
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// []interface{}指明是 interface{} 的数组类型：[]interface{}
    // []interface{}{} 是 []interface{}类型的值
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{} 
}
~~~

紧接着的就是拼接各个字句，组成 SQL 语句：

~~~go
package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type Clause struct {
	sql     map[Type]string        // Type -- SQL
	sqlVars map[Type][]interface{} // Type -- Vars
}

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	//FIXME 根据 name 生成对应的 SQL 语句，此处一定要注意 vars...
	sql, vars := generators[name](vars...)

	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	// 依据 orders 构造完整的 SQL 语句
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
~~~

创建字句，以及构造完整的 SQL 语句的测试用例：

~~~go
package clause

import (
	"reflect"
	"testing"
)

func TestGenBindVars(t *testing.T) {
	result := genBindVars(3)
	if result != "?, ?, ?" {
		t.Fatalf("genBindVars error, got:(%s), want:(%s)", result, "?, ?, ?")
	}
}

func TestClasue(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name=?", "Tom")
	clause.Set(ORDERBY, "Age ASC")

	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	// SELECT * FROM User WHERE Name=? ORDER BY Age ASC LIMIT ? [Tom 3]
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name=? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestInsert(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"Name", "Age"})

	sql, vars := clause.Build(INSERT)
	t.Log(sql, vars)

	if sql != "INSERT INTO User (Name,Age)" {
		t.Fatal("failed to build SQL statement")
	}
}

func TestValues(t *testing.T) {
	var clause Clause

	clause.Set(VALUES, []interface{}{"Tom", "18"}, []interface{}{"Sam", 29})
	sql, vars := clause.Build(VALUES)
	// VALUES (?, ?), (?, ?) [Tom 18 Sam 29]
	t.Log(sql, vars)
}
~~~

可以看出 VALUES 的参数是为多组参数准备的，比如上述测试用例中的 `{"Tom", "18"}` 和 `{"Sam", 29}`。

实现 Insert 功能：

~~~go
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// INSERT INTO table_name(col1, col2, col3,...) VALUES (a1, a2, a3, ...), (b1, b2, b3, ...),...
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable() // 执行 Parse
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
        // 解析出对象中各个字段的值
		recordValues = append(recordValues, table.RecordValues(value)) 
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
~~~

最重要的一个步骤是：根据数据库中列的顺序，从对象中找到对应的值：

~~~go
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest)) // reflect.Value
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		// 顺序严格和 struct 定义中各个字段顺序一致
		log.Infof("field.Name:%s", field.Name)
		// reflect.Value struct --> value
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
~~~

在 ORM 角度来看，Insert 功能实际上就是将对象信息存储到 RDBMS 中；反之，对应的就是 Find 功能。

~~~go
~~~













# 链式操作与更新删除









# 实现钩子 Hooks









# 支持事务 Transaction













# 数据库迁移 Migrate