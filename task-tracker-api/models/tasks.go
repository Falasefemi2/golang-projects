// Package models
package models

import (
	"errors"
	"time"
)

const (
	maxTitleLength = 50
	maxDescLength  = 200

	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

var (
	ErrEmptyInput = errors.New("title and description cannot be empty")
	ErrTooLong    = errors.New("title or description too long")
	ErrNotFound   = errors.New("task not found")
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
