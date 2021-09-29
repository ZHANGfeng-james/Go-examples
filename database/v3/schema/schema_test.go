package schema

import (
	"testing"

	"github.com/go-examples-with-tests/database/v2/dialect"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"` // struct的TAG有固定的格式，写错则无效！
	Age  int
}

func TestSchema(t *testing.T) {
	dialect, _ := dialect.GetDialect("sqlite3")

	user := &User{}
	userSchema := Parse(user, dialect)
	if userSchema.Name != "User" && len(userSchema.Fields) != 2 {
		t.Fatal("schema parse User error")
	}
	if userSchema.fieldMap["Name"].Tag != "PRIMARY KEY" {
		t.Fatal("schema parse User error")
	}
}

func TestRecordValue(t *testing.T) {
	user := &User{
		Name: "Tom",
		Age:  18,
	}

	dialect, _ := dialect.GetDialect("sqlite3")

	schema := Parse(user, dialect)
	values := schema.RecordValues(user)

	name := values[0].(string)
	age := values[1].(int)
	if name != "Tom" && age != 18 {
		t.Fatal("record value is error")
	}
}

type Password struct {
	Len     int
	Content string
}

func (p *Password) TableName() string {
	return "test_password_name"
}

func TestSchemaPassword(t *testing.T) {
	dialect, _ := dialect.GetDialect("sqlite3")

	password := &Password{}
	passwordSchema := Parse(password, dialect)
	if passwordSchema.Name != password.TableName() && len(passwordSchema.Fields) != 2 {
		t.Fatal("schema parse Password error")
	}
}
