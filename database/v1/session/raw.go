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
	log.Info(s.sql.String, s.sqlVars)
	// 调用的是 sql.DB 的 QueryRow 函数，仅返回一行结果
	return s.DB().QueryRow(s.sql.String(), s.sqlVars)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String, s.sqlVars)
	// 调用的是 sql.DB 的 Query 函数，可返回多行结果
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars); err != nil {
		log.Error(err)
	}
	return
}
