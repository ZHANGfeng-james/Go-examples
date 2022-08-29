package testing_test

import (
	"database/sql"
	"strings"
)

// someOutput return some string value
func someOutput(input string) string {
	return strings.ToUpper(input)
}

func recordStats(db *sql.DB, userID, productID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
		return err
	}
	if _, err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", userID, productID); err != nil {
		return err
	}
	return nil
}
