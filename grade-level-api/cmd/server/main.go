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

	http.HandleFunc(
		"/admin/users",
		middleware.RoleAuth(handler.GetAllUsers, db.Admin),
	)

	http.HandleFunc(
		"/semesters",
		middleware.RoleAuth(handler.SemestersHandler, db.Admin, db.Lecturer, db.Student),
	)

	http.HandleFunc(
		"/semesters/",
		middleware.RoleAuth(handler.SemesterByIDHandler, db.Admin),
	)

	http.HandleFunc(
		"/courses",
		middleware.RoleAuth(handler.CreateCourse, db.Lecturer),
	)

	http.ListenAndServe(":8080", nil)
}
