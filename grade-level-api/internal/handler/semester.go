// Package handler
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/utils"
)

type CreateSemesterRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func CreateSemester(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateSemesterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid start_date format (YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid end_date format (YYYY-MM-DD)")
		return
	}

	semester, err := db.CreateSemester(
		db.Semester(req.Name),
		startDate,
		endDate,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, semester)
}

func GetAllSemesters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	semesters, err := db.ListSemesters()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, semesters)
}

func GetSemesterByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "missing semester id")
		return
	}

	semesterID, err := strconv.Atoi(idStr)
	if err != nil || semesterID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid semester id")
		return
	}

	semester, err := db.FindSemesterByID(semesterID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, semester)
}

func UpdateSemester(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		utils.WriteError(w, http.StatusBadRequest, "invalid URL path")
		return
	}

	semesterID, err := strconv.Atoi(parts[2])
	if err != nil || semesterID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid semester id")
		return
	}

	var req CreateSemesterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid start_date format (YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid end_date format (YYYY-MM-DD)")
		return
	}

	semester, err := db.UpdateSemester(
		semesterID,
		db.Semester(req.Name),
		startDate,
		endDate,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, semester)
}

func DeleteSemester(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		utils.WriteError(w, http.StatusBadRequest, "invalid URL path")
		return
	}
	semesterID, err := strconv.Atoi(parts[2])
	if err != nil || semesterID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid semester id")
		return
	}
	semester := db.DeleteSemester(semesterID)
	if semester != nil {
		utils.WriteError(w, http.StatusNotFound, semester.Error())
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
