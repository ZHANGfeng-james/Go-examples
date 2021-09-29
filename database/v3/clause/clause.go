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

// 每一个 Clause 实例，就对应的是一个 SQL 语句
type Clause struct {
	sql     map[Type]string        // Type -- SQL
	sqlVars map[Type][]interface{} // Type -- Vars
}

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	// 根据 name 生成对应的 SQL 语句，此处一定要注意 vars...
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
