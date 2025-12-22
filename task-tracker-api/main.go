package main

import (
	"net/http"

	"github.com/falasefemi2/ask-tracker-api/auth"
	"github.com/falasefemi2/ask-tracker-api/db"
	"github.com/falasefemi2/ask-tracker-api/handlers"
)

func main() {
	db.Init()

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.Handle(
		"/tasks/",
		auth.JWTMiddleware(http.HandlerFunc(handlers.TasksHandler)),
	)

	http.ListenAndServe(":8080", nil)
}
