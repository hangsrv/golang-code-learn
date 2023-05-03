package main

import (
	"database/sql"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

// docker run -id --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
/*
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name` varchar(45) DEFAULT '',
    `age` int(11) NOT NULL DEFAULT '0',
    `sex` tinyint(3) NOT NULL DEFAULT '0',
    `phone` varchar(45) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
*/
func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/hang")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("connect success~~")
}
