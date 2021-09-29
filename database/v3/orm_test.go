package orm

import (
	"testing"

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
