package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/middleware"
	"github.com/falasefemi2/gradesystem/utils"
)

func CoursesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	switch r.Method {

	// ---------------- CREATE COURSE (LECTURER ONLY) ----------------
	case http.MethodPost:
		if user.Role != string(db.Lecturer) {
			utils.WriteError(w, http.StatusForbidden, "lecturer access required")
			return
		}

		var req CreateCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Name == "" || req.Level <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "name and level are required")
			return
		}

		course, err := db.CreateCourse(req.Name, req.Level, user.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusCreated, course)

	// ---------------- LIST COURSES ----------------
	case http.MethodGet:
		levelParam := r.URL.Query().Get("level")

		// Lecturer → only their courses
		if user.Role == string(db.Lecturer) {
			courses, err := db.FindCoursesByLecturerID(user.ID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			utils.WriteJSON(w, http.StatusOK, courses)
			return
		}

		// Student/Admin → filter by level if provided
		if levelParam != "" {
			level, err := strconv.Atoi(levelParam)
			if err != nil || level <= 0 {
				utils.WriteError(w, http.StatusBadRequest, "invalid level")
				return
			}

			courses, err := db.FindCoursesByLevel(level)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			utils.WriteJSON(w, http.StatusOK, courses)
			return
		}

		// Admin → all courses
		if user.Role == string(db.Admin) {
			courses, err := db.ListCourses()
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			utils.WriteJSON(w, http.StatusOK, courses)
			return
		}

		utils.WriteError(w, http.StatusForbidden, "access denied")

	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func CourseByIDHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/courses/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	switch r.Method {

	case http.MethodGet:
		course, err := db.FindCourseByID(id)
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, course)

	case http.MethodPut:
		if user.Role != string(db.Lecturer) {
			utils.WriteError(w, http.StatusForbidden, "lecturer access required")
			return
		}

		var req CreateCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		course, err := db.UpdateCourse(id, req.Name, req.Level, user.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, course)

	case http.MethodDelete:
		if user.Role != string(db.Admin) {
			utils.WriteError(w, http.StatusForbidden, "admin access required")
			return
		}

		if err := db.DeleteCourse(id); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "course deleted"})

	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
