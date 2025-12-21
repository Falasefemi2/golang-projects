package main

import (
	"net/http"

	"github.com/falasefemi2/ask-tracker-api/db"
	"github.com/falasefemi2/ask-tracker-api/handlers"
)

func main() {
	db.Init()

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.ListenAndServe(":8080", nil)
}
