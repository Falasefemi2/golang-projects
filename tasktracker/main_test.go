package main

import (
	"strings"
	"testing"
)

func TestAddTask(t *testing.T) {
	tests := []struct {
		name   string
		title  string
		desc   string
		expErr bool
	}{
		{
			"valid title and description",
			"Laundary",
			"Doing a lot of washing clothes",
			false,
		},
		{
			"empty title",
			"",
			"some description",
			true,
		},
		{
			"empty description",
			"some title",
			"",
			true,
		},
		{
			"title too long",
			strings.Repeat("a", 51),
			"valid description",
			true,
		},
		{
			"description too long",
			"valid title",
			strings.Repeat("a", 500),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskList := &TaskList{[]Task{}}
			err := taskList.AddTask(tt.title, tt.desc)

			if (err != nil) != tt.expErr {
				t.Fatalf("AddTask() error: %v, WantErr: %v", err, tt.expErr)
			}
			if !tt.expErr {
				if len(taskList.Tasks) != 1 {
					t.Fatalf("expected one task got %d", len(taskList.Tasks))
				}
				task := taskList.Tasks[0]
				if task.Status != StatusTodo {
					t.Errorf("expected status todo, got %s", task.Status)
				}
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		initialTask []Task
		deletedID   int
		wantErr     bool
		wantLen     int
	}{
		{
			"delete a valid ID",
			[]Task{
				{ID: 1, Title: "valid title", Description: "valid description"},
			},
			1,
			false,
			0,
		},
		{
			"delete invalid ID",
			[]Task{
				{ID: 1, Title: "valid title", Description: "valid description"},
			},
			44,
			true,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskList := &TaskList{
				tt.initialTask,
			}
			err := taskList.DeleteTask(tt.deletedID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() Error: %v want error: %v", err, tt.wantErr)
			}

			if len(taskList.Tasks) != tt.wantLen {
				t.Fatalf(
					"expected %d tasks after delete, got %d",
					tt.wantLen,
					len(taskList.Tasks),
				)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		title     string
		desc      string
		initial   []Task
		wantErr   bool
		wantTitle string
		wantDesc  string
	}{
		{
			name:  "update existing task",
			id:    1,
			title: "New title",
			desc:  "New description",
			initial: []Task{
				{ID: 1, Title: "Old title", Description: "Old description"},
			},
			wantErr:   false,
			wantTitle: "New title",
			wantDesc:  "New description",
		},
		{
			name:  "update non-existing task",
			id:    99,
			title: "New title",
			desc:  "New description",
			initial: []Task{
				{ID: 1, Title: "Old title", Description: "Old description"},
			},
			wantErr:   true,
			wantTitle: "Old title",
			wantDesc:  "Old description",
		},
		{
			name:  "empty title",
			id:    1,
			title: "",
			desc:  "New description",
			initial: []Task{
				{ID: 1, Title: "Old title", Description: "Old description"},
			},
			wantErr:   true,
			wantTitle: "Old title",
			wantDesc:  "Old description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := TaskList{
				Tasks: tt.initial,
			}

			err := tl.UpdateTask(tt.id, tt.title, tt.desc)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			task := tl.Tasks[0]
			if task.Title != tt.wantTitle {
				t.Errorf("title mismatch: got %q, want %q", task.Title, tt.wantTitle)
			}
			if task.Description != tt.wantDesc {
				t.Errorf("description mismatch: got %q, want %q", task.Description, tt.wantDesc)
			}
		})
	}
}

func TestMarkStatus(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		status     string
		initial    []Task
		wantErr    bool
		wantStatus string
	}{
		{
			name:   "change status from todo to in_progress",
			id:     1,
			status: StatusInProgress,
			initial: []Task{
				{
					ID:          1,
					Title:       "test title",
					Description: "test description",
					Status:      StatusTodo,
				},
			},
			wantErr:    false,
			wantStatus: StatusInProgress,
		},
		{
			name:   "invalid status",
			id:     1,
			status: "donee",
			initial: []Task{
				{
					ID:          1,
					Title:       "test title",
					Description: "test description",
					Status:      StatusTodo,
				},
			},
			wantErr:    true,
			wantStatus: StatusTodo,
		},
		{
			name:   "non-existing task id",
			id:     99,
			status: StatusDone,
			initial: []Task{
				{
					ID:          1,
					Title:       "test title",
					Description: "test description",
					Status:      StatusTodo,
				},
			},
			wantErr:    true,
			wantStatus: StatusTodo,
		},
		{
			name:   "status unchanged on error",
			id:     1,
			status: "",
			initial: []Task{
				{
					ID:          1,
					Title:       "test title",
					Description: "test description",
					Status:      StatusTodo,
				},
			},
			wantErr:    true,
			wantStatus: StatusTodo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := TaskList{
				Tasks: tt.initial,
			}
			err := tl.MarkStatus(tt.id, tt.status)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tl.Tasks[0].Status != tt.wantStatus {
				t.Errorf(
					"status mismatch: got %q, want %q",
					tl.Tasks[0].Status,
					tt.wantStatus,
				)
			}
		})
	}
}
