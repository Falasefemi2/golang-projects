package models

type Enrollment struct {
	ID         int `json:"id"`
	StudentID  int `json:"studentID"`
	CourseID   int `json:"courseID"`
	SemesterID int `json:"semesterID"`
}
