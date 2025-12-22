package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/falasefemi2/ask-tracker-api/models"
)

func CreateTask(userID int, title, description string) (*models.Task, error) {
	now := time.Now()
	task := &models.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := "INSERT INTO tasks (user_id, title, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	res, err := DB.Exec(query, userID, title, description, now, now)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return nil, err
	}
	task.ID = int(id)
	return task, nil
}

func UpdateTask(id, userID int, title, description, status string) (*models.Task, error) {
	if title == "" && description == "" && status == "" {
		return nil, errors.New("title, description, or status must be provided")
	}
	setClauses := []string{}
	args := []interface{}{}
	if title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, title)
	}
	if description != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, description)
	}
	if status != "" {
		setClauses = append(setClauses, "status = ?")
		args = append(args, status)
	}
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)
	query := `
		UPDATE tasks
		SET ` + strings.Join(setClauses, ", ") + `
		WHERE id = ? AND user_id = ?
	`
	args = append(args, id, userID)
	result, err := DB.Exec(query, args...)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("task not found or not authorized")
	}
	var task models.Task
	var createdAtStr string
	var updatedAtStr string

	err = DB.QueryRow(`
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`, id, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found after update")
		}
		return nil, err
	}

	task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	task.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func DeleteTask(id, userID int) error {
	query := `
DELETE FROM tasks
where id = ? AND user_id = ?
	`

	result, err := DB.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	}

	fmt.Printf("%d row(s) deleted\n", rowsAffected)
	return nil
}

func GetUserTasks(userID int) ([]models.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var createdAtStr string
		var updatedAtStr string

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		task.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
