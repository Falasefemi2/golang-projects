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
	//  dsn := "root:admin@tcp(localhost:3306)/gradingsystem"
	// DB, err = sql.Open("mysql", dsn)
	DB, err = sql.Open(
		"mysql",
		"root:admin@tcp(localhost:3306)/gradingsystem?parseTime=true")
	if err != nil {
		log.Fatal("Error opening database", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to db", err)
	}
	fmt.Println("Database connected successfully!")
}
