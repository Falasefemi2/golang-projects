// Package handler
package handler

import (
	"net/http"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/middleware"
	"github.com/falasefemi2/gradesystem/utils"
)

func SemesterByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetSemesterByID(w, r)
	case http.MethodPut:
		UpdateSemester(w, r)
	case http.MethodDelete:
		DeleteSemester(w, r)
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func SemestersHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	switch r.Method {
	case http.MethodPost:
		if user.Role != string(db.Admin) {
			utils.WriteError(w, http.StatusForbidden, "admin access required")
			return
		}
		CreateSemester(w, r)

	case http.MethodGet:
		GetAllSemesters(w, r)

	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
