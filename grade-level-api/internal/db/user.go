package db

import (
	"database/sql"
	"errors"

	"github.com/falasefemi2/gradesystem/internal/auth"
	"github.com/falasefemi2/gradesystem/internal/models"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type Role string

const (
	Student  Role = "student"
	Lecturer Role = "lecturer"
	Admin    Role = "admin"
)

func CreateUser(db DBExecutor, name, email, password string, role Role) (*models.User, error) {
	if err := auth.ValidatePassword(password); err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	result, err := db.Exec("INSERT INTO user (name, email, password, role) VALUES (?, ?, ?, ?)", name, email, hash, string(role))
	if err != nil {
		return nil, errors.New("email already exists or database error")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       int(id),
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     string(role),
	}, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(
		"SELECT id, name, email, password, role FROM user WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(id int) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(
		"SELECT id, name, email, password, role FROM user WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(
		"SELECT id, name, email, role FROM user",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersByRole(role Role) ([]models.User, error) {
	rows, err := DB.Query(
		"SELECT id, name, email, role FROM user WHERE role = ?",
		string(role),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func VerifyUser(email, password string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := auth.VerifyPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
