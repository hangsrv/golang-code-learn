package main

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

// docker run -id --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
/*
CREATE DATABASE test;
USE test;
CREATE TABLE `user` (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(50) DEFAULT NULL,
	is_active BOOLEAN DEFAULT FALSE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

desc user;

ALTER TABLE user
ADD column1 BOOLEAN DEFAULT true,
ADD column2 VARCHAR(50) NULL,
ADD column3 VARCHAR(50) NOT NULL;

ALTER TABLE user
DROP column1,
DROP column2,
DROP column3;

ALTER TABLE USER
MODIFY  column1 BOOLEAN DEFAULT true;

*/

type User struct {
	ID       int
	Name     string
	IsActive bool
	Test     string
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 插入数据并获取自增ID
	id1 := insertData(db, "John Doe", true)
	fmt.Printf("Inserted data with ID: %d\n", id1)

	id2 := insertData(db, "Jane Smith", false)
	fmt.Printf("Inserted data with ID: %d\n", id2)

	id3 := insertData(db, "Michael Johnson", true)
	fmt.Printf("Inserted data with ID: %d\n", id3)

	fmt.Println("===============================")

	// 根据ID查询数据
	id := int(id2)
	result := queryDataByID(db, id)
	if result != nil {
		fmt.Printf("Data with ID %d: %+v\n", id, result)
	} else {
		fmt.Printf("No data found with ID %d\n", id)
	}

	fmt.Println("===============================")

	// 查询全部数据
	allData := queryAllData(db)
	fmt.Println("All Data:")
	for _, data := range allData {
		fmt.Printf("%+v\n", data)
	}

	fmt.Println("===============================")

	// 删除全部数据
	err = deleteAllData(db)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("All data deleted successfully.")
	}
}

// 插入数据并返回自增ID
func insertData(db *sql.DB, name string, isActive bool) int64 {
	query := "INSERT INTO user (name, is_active) VALUES (?, ?)"
	result, err := db.Exec(query, name, isActive)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return id
}

// 根据ID查询数据
func queryDataByID(db *sql.DB, id int) *User {
	query := "SELECT id, name, is_active, test FROM user WHERE id = ?"
	row := db.QueryRow(query, id)

	var data User
	err := row.Scan(&data.ID, &data.Name, &data.IsActive, &data.Test)
	if err != nil {
		log.Fatal(err)
		if err == sql.ErrNoRows {
			return nil
		}
		fmt.Println(err)
	}

	return &data
}

// 查询全部数据
func queryAllData(db *sql.DB) []User {
	query := "SELECT id, name, is_active FROM user"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var allData []User
	for rows.Next() {
		var data User
		err := rows.Scan(&data.ID, &data.Name, &data.IsActive)
		if err != nil {
			fmt.Println(err)
			continue
		}
		allData = append(allData, data)
	}

	return allData
}

// 删除全部数据
func deleteAllData(db *sql.DB) error {
	query := "DELETE FROM user"
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
