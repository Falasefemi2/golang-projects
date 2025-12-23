package handler

import (
	"encoding/json"
	"net/http"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/middleware"
	"github.com/falasefemi2/gradesystem/utils"
)

type CreateCourseRequest struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Level <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required and must be valid")
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user.Role != string(db.Lecturer) {
		utils.WriteError(w, http.StatusForbidden, "lecturer access required")
		return
	}

	course, err := db.CreateCourse(
		req.Name,
		req.Level,
		user.ID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, course)
}
