// Package handler
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/middleware"
	"github.com/falasefemi2/gradesystem/utils"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	role := db.Role(req.Role)
	user, err := db.CreateUser(
		db.DB,
		req.Name, req.Email, req.Password, role)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = ""
	utils.WriteJSON(w, http.StatusCreated, user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := db.VerifyUser(req.Email, req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := middleware.GenerateJWT(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	user, err := db.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
