package handlers

import (
	"encoding/json"
	"net/http"

	"expense-tracker/auth"
	"expense-tracker/internal/db"
	"expense-tracker/utils"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var c credentials
	json.NewDecoder(r.Body).Decode(&c)

	hash, _ := auth.HashPassword(c.Password)
	_, err := db.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		c.Email, hash,
	)
	if err != nil {
		http.Error(w, "user exits", http.StatusUnauthorized)
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var c credentials
	json.NewDecoder(r.Body).Decode(&c)

	row := db.DB.QueryRow(
		"SELECT id, password_hash FROM users WHERE email = ?", c.Email,
	)

	var id int
	var hash string
	if row.Scan(&id, &hash) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if auth.CheckPassword(hash, c.Password) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(id)
	utils.JSON(w, http.StatusOK, map[string]string{"token": token})
}
