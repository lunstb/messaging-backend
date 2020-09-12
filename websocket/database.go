package websocket

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Hidden driver for sql import
)

// Test function purely for testing
func Test() {
	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/messaging")
	if err != nil {
		log.Fatal(err)
	}

	var (
		num int
	)
	rows, err := db.Query("select num from test")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&num)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(num)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
