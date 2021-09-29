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
