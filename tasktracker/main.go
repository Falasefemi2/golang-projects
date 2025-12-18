// Project main
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	taskFile       = "tasks.json"
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
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func printUsage() {
	fmt.Println(`
		Usage:
  task-cli add <title>                    Add a new task
  task-cli update <id> <new-title>        Update a task
  task-cli delete <id>                    Delete a task
  task-cli mark-in-progress <id>          Mark task as in progress
  task-cli mark-done <id>                 Mark task as done
  task-cli list                           List all tasks
  task-cli list <status>                  List tasks by status (todo, in_progress, done)
	`)
}

func (tl *TaskList) getNextID() int {
	maxID := 0
	for _, task := range tl.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func (tl *TaskList) AddTask(title, description string) error {
	if title == "" || description == "" {
		return ErrEmptyInput
	}
	if len(title) > maxTitleLength || len(description) > maxDescLength {
		return ErrTooLong
	}
	now := time.Now().UTC()
	task := Task{
		ID:          tl.getNextID(),
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tl.Tasks = append(tl.Tasks, task)
	return nil
}

func (tl *TaskList) DeleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (tl *TaskList) UpdateTask(id int, title, description string) error {
	if title == "" || description == "" {
		return ErrEmptyInput
	}
	if len(title) > maxTitleLength || len(description) > maxDescLength {
		return ErrTooLong
	}
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Title = title
			tl.Tasks[i].Description = description
			tl.Tasks[i].UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return ErrNotFound
}

func isValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (tl *TaskList) MarkStatus(id int, status string) error {
	if !isValidStatus(status) {
		return errors.New("status not available")
	}
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Status = status
			tl.Tasks[i].UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return ErrNotFound
}

func (t *TaskList) FilterTasks(status string) []Task {
	if status == "" {
		return t.Tasks
	}

	var result []Task
	for _, task := range t.Tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	return result
}

func (t *TaskList) List(status string) {
	tasks := t.FilterTasks(status)

	if len(tasks) == 0 {
		if status == "" {
			fmt.Println("No tasks found")
		} else {
			fmt.Printf("No tasks with status %q\n", status)
		}
		return
	}

	fmt.Println(strings.Repeat("=", 80))
	for _, task := range tasks {
		fmt.Printf(
			"ID: %d | Status: %s | Title: %s | Description: %s\n",
			task.ID,
			task.Status,
			task.Title,
			task.Description,
		)
	}
	fmt.Println(strings.Repeat("=", 80))
}

func loadTasks() (*TaskList, error) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{Tasks: []Task{}}, nil
		}
		return nil, err
	}

	var taskList TaskList
	if err := json.Unmarshal(data, &taskList); err != nil {
		return nil, err
	}
	return &taskList, nil
}

func saveTasks(taskList *TaskList) error {
	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(taskFile, data, 0o644)
}

func withTaskList(fn func(*TaskList) error) error {
	taskList, err := loadTasks()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	if err := fn(taskList); err != nil {
		return err
	}

	if err := saveTasks(taskList); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {

	case "add":
		if len(os.Args) < 4 {
			fmt.Println("add requires <title> <description>")
			return
		}

		title := os.Args[2]
		description := os.Args[3]

		err := withTaskList(func(tl *TaskList) error {
			return tl.AddTask(title, description)
		})
		if err != nil {
			fmt.Println("error:", err)
		}

	case "update":
		if len(os.Args) < 5 {
			fmt.Println("update requires <id> <title> <description>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}

		title := os.Args[3]
		description := os.Args[4]

		err = withTaskList(func(tl *TaskList) error {
			return tl.UpdateTask(id, title, description)
		})
		if err != nil {
			fmt.Println("error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("delete requires <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}

		err = withTaskList(func(tl *TaskList) error {
			return tl.DeleteTask(id)
		})
		if err != nil {
			fmt.Println("error:", err)
		}

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("mark-in-progress requires <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}

		err = withTaskList(func(tl *TaskList) error {
			return tl.MarkStatus(id, StatusInProgress)
		})
		if err != nil {
			fmt.Println("error:", err)
		}

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("mark-done requires <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}

		err = withTaskList(func(tl *TaskList) error {
			return tl.MarkStatus(id, StatusDone)
		})
		if err != nil {
			fmt.Println("error:", err)
		}

	case "list":
		var status string
		if len(os.Args) > 2 {
			status = os.Args[2]
		}

		taskList, err := loadTasks()
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		taskList.List(status)

	default:
		printUsage()
	}
}
