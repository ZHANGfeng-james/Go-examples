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

func TestMain(m *testing.M) {
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

func TestExec(t *testing.T) {
	session := New(TestDB, TestDialect)

	_, _ = session.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = session.Raw("CREATE TABLE User(Name text);").Exec()

	result, _ := session.Raw("INSERT INTO User('Name') VALUES (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got:", count)
	}
}

type Person struct {
	Name string
	Age  int8
}

func TestInsert(t *testing.T) {
	session := New(TestDB, TestDialect)

	session.Model(&Person{})
	if !session.HasTable() {
		session.CreateTable()
	} else {
		session.DropTable()
		session.CreateTable()
	}

	if session.HasTable() {
		p1 := &Person{
			Name: "Katyusha",
			Age:  31,
		}
		p2 := &Person{
			Name: "Sam",
			Age:  32,
		}
		p3 := &Person{
			Name: "Jason",
			Age:  33,
		}

		count, err := session.Insert(p1, p2, p3)
		if err != nil || count != 3 {
			t.Fatal("insert action error")
		}
	}
}
