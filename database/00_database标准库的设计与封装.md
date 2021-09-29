对于 `database/sql` 标准库来说，

主要的函数：

~~~shell
ant@MacBook-Pro ~ % go doc database/sql |grep "^func"
func Drivers() []string
func Register(name string, driver driver.Driver)
~~~

主要的结构体类型：

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

其中最重要的是：DB、Stmt、Tx 和 Conn 这 4 个类型。