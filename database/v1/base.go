package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-examples-with-tests/database/v1/orm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string
	Age  int8
}

func (user *User) String() string {
	return fmt.Sprintf("User: Name=%s, Age=%d\n", user.Name, user.Age)
}

func main() {
	engine, err := orm.NewEngine("sqlite3", "../gee.db")
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	session := engine.NewSession()
	session.Raw("DROP TABLE IF EXISTS User;").Exec()

	session.Raw("CREATE TABLE User(Name text)").Exec()
	session.Raw("CREATE TABLE User(Name text)").Exec()

	result, _ := session.Raw("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}

func test() {
	db, err := sql.Open("sqlite3", "../gee.db")
	if err != nil {
		fmt.Println("err!")
		log.Fatal(err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	row, err := db.Query("SELECT * FROM User")
	if err != nil {
		log.Fatal(err.Error())
	}

	var name string
	var age int8
	for row.Next() {
		if err = row.Scan(&name, &age); err == nil {
			fmt.Printf("User: Name=%s, Age=%d\n", name, age)
		} else {
			fmt.Println(err.Error())
		}
	}

	drivers := sql.Drivers()
	for _, driver := range drivers {
		fmt.Println(driver)
	}
}
