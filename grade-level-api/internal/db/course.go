// Package db
package db

import (
	"errors"

	"github.com/falasefemi2/gradesystem/internal/models"
)

func CreateCourse(name string, level, lecturerID int) (*models.Course, error) {
	if name == "" {
		return nil, errors.New("course name  cannot be empty")
	}

	if level <= 0 {
		return nil, errors.New("course level must be a positive integer")
	}

	result, err := DB.Exec(
		`INSERT INTO course (name, level, lecturer_id) VALUES (?, ?, ?)`,
		name, level, lecturerID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Course{
		ID:         int(id),
		Name:       name,
		Level:      level,
		LecturerID: lecturerID,
	}, nil
}
