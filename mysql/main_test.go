package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// insert into user(is_active) values(false)
func Test_selectNull(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u := queryDataByID(db, 3)
	fmt.Println(u)
}
