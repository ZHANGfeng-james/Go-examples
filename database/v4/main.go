package main

import (
	"context"
	"log"

	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/go-examples-with-tests/database/v4/store"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	log.Println("hello, world!")

	options := pkg.NewMySQLOptions()
	initOptions(options)
	log.Printf("MariaDB options:%s", options)

	// connection to MariaDB
	dbFactory, err := store.GetMySQLFactoryOr(options)
	if err != nil {
		log.Fatal(err)
	}

	// CRUD，拿着 userStore 就可以和 MariaDB 交互
	userStore := dbFactory.Users()
	user, err := userStore.Get(context.TODO(), "admin", metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(user)
	}

	// close DB connection
	dbFactory.Close()
}

func initOptions(options *pkg.MySQLOptions) {
	if options == nil {
		return
	}
	options.Username = "iam"
	options.Password = "iam59!z$"
	options.Database = "iam"
	options.LogLevel = 4
}
