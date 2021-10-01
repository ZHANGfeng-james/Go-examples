package orm

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-examples-with-tests/database/v2/log"
	"github.com/go-examples-with-tests/database/v3/session"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDb(t *testing.T) *Engine {
	t.Helper()

	engine, err := NewEngine("sqlite3", "../gee.db")
	if err != nil {
		t.Fatal("failed to connect:", err)
	}
	return engine
}

func TestORM(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
}

func TestSQLTransaction(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../gee.db")
	defer db.Close()

	// CREATE TABLE Account (ID integer ,Password text );
	_, _ = db.Exec("CREATE TABLE IF NOT EXISTS Account;")

	tx, _ := db.Begin()
	_, err1 := tx.Exec("INSERT INTO Account('ID', 'Password') VALUES (?, ?);", "1", "sdi")
	_, err2 := tx.Exec("INSERT INTO Account('ID', 'Password') VALUES (?, ?);", "2", "sdy")
	if err1 != nil || err2 != nil {
		_ = tx.Rollback()
		log.Info("Rollback", err1, err2)
	} else {
		_ = tx.Commit()
		log.Info("Commit")
	}
}

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func TestTransaction(t *testing.T) {
	engine, err := NewEngine("sqlite3", "../gee.db")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&Account{}).DropTable()
	_, err = engine.Transaction(func(s *session.Session) (interface{}, error) {
		// 此处的入参是来自 engine.Transaction 方法中
		_ = s.Model(&Account{}).CreateTable()
		_, err = s.Insert(&Account{ID: 1, Password: "123456"})
		// 此处故意返回一个 error 实例，以此触发 Rollback
		return nil, errors.New("ERROR")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}
