package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/falasefemi2/ask-tracker-api/auth"
	"github.com/falasefemi2/ask-tracker-api/db"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		http.Error(w, "Could not retrieve user from context", http.StatusInternalServerError)
		return
	}
	user, err := db.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	task, err := db.CreateTask(user.ID, req.Title, req.Description)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Task ID must be a number", http.StatusBadRequest)
		return
	}
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		http.Error(w, "Could not retrieve user from context", http.StatusInternalServerError)
		return
	}
	user, err := db.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	task, err := db.UpdateTask(
		taskID,
		user.ID,
		req.Title,
		req.Description,
		req.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
