package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Hidden driver for sql import
)

var once sync.Once

// Database provides a type
type Database struct {
	database *sql.DB
}

// variavel Global
var instance *Database

// Connect provides a singleton pattern for router
func Connect() *Database {

	once.Do(func() {
		db, err := sql.Open("mysql",
			"root:@tcp(127.0.0.1:3306)/messaging")
		if err != nil {
			log.Fatal(err)
		}

		instance = &Database{database: db}
	})

	return instance
}

// InsertMessage function
func (r *Database) InsertMessage(sender string, msg string, conversation string) {
	var (
		participants string
	)
	var participantList string

	rows, err := r.database.Query("SELECT participants FROM conversations WHERE id = ?", conversation)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&participants)
		if err != nil {
			log.Fatal(err)
		}
		participantList = participants
	}

	thing := strings.Split(participantList, "|")

	found := Find(thing, sender)
	if !found {
		fmt.Println("Value not found in slice")
		return
	}
}

// Test function purely for testing
func (r *Database) Test() {
	var (
		num int
	)
	rows, err := r.database.Query("select num from test")
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
}

// Find takes a slice and looks for an element in it
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
