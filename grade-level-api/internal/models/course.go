package models

type Course struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	LecturerID int    `json:"LecturerID"`
}
