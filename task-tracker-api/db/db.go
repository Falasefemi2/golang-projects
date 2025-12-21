// Package db

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	dsn := "root:admin@tcp(localhost:3306)/tasktracker"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	fmt.Println("Database connected successfully!")
}
