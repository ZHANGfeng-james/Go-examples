package orm

import (
	"database/sql"

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
