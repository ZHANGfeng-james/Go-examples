package session

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-examples-with-tests/database/v1/log"
	"github.com/go-examples-with-tests/database/v2/dialect"
	"github.com/go-examples-with-tests/database/v2/schema"
)

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
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String, s.sqlVars)
	// 调用的是 sql.DB 的 Query 函数，可返回多行结果
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

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
