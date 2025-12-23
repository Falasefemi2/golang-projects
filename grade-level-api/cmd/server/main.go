// Package main
package main

import (
	"net/http"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/handler"
	"github.com/falasefemi2/gradesystem/internal/middleware"
)

func main() {
	db.Init()

	http.HandleFunc("/signup", handler.SignUp)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/admin/users", middleware.RoleAuth(handler.GetAllUsers, db.Admin))

	http.ListenAndServe(":8080", nil)
}
