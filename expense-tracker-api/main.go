package main

import (
	"log"
	"net/http"

	"expense-tracker/handlers"
	"expense-tracker/internal/db"
	"expense-tracker/middleware"
)

func main() {
	db.Init()
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", handlers.Signup)
	mux.HandleFunc("/login", handlers.Login)

	mux.Handle("/expenses",
		middleware.Auth(http.HandlerFunc(handlers.CreateExpense)),
	)

	mux.Handle("/expenses/list",
		middleware.Auth(http.HandlerFunc(handlers.ListExpenses)),
	)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
