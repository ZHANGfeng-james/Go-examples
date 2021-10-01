package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-examples-with-tests/database/v3/dialect"
	"github.com/go-examples-with-tests/database/v3/log"
	"github.com/go-examples-with-tests/database/v3/session"
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

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	session := engine.NewSession()
	if err = session.Begin(); err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		log.Info("Transaction run...")
		if p := recover(); p != nil {
			session.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			log.Error(err.Error())
			_ = session.Rollback() // err is non-nil; don't change it
		} else {
			err = session.Commit() // err is nil; if Commit returns error update err
		}
	}()
	// 执行顺序：f(session) --> defer func(){}() 此时 err 变量已被 f(session) 赋值
	return f(session)
}

// difference get the difference of a - b
func difference(a, b []string) (diff []string) {
	mapD := make(map[string]bool)
	for _, v := range b {
		mapD[v] = true
	}

	for _, v := range a {
		if _, ok := mapD[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		// value interface{} --> new table with column changed
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		table := s.RefTable()
		// 虽然此处 table 的 column 改变了，但是 table_name 没有改变
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		log.Infof("origin table columns:%v", columns)

		addCols := difference(table.FieldNames, columns) // new - old = 在 new 中挑选 old 没有的
		delCols := difference(columns, table.FieldNames) // old - new = 在 old 中挑选 new 没有的
		log.Infof("added cols:%v, deleted cols:%s", addCols, delCols)

		for _, col := range addCols {
			field := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, field.Name, field.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ") // new columns
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))

		_, err = s.Exec()

		return
	})

	return err
}
