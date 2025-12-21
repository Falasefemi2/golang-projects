package db

import (
	"database/sql"
	"errors"

	"github.com/falasefemi2/ask-tracker-api/auth"
	"github.com/falasefemi2/ask-tracker-api/models"
)

func CreateUser(email, password string) (*models.User, error) {
	if err := auth.ValidatePassword(password); err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	result, err := DB.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, hash)
	if err != nil {
		return nil, errors.New("email already exists or database error")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           int(id),
		Email:        email,
		PasswordHash: hash,
	}, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func VerifyUser(email, password string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := auth.VerifyPassword(password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
