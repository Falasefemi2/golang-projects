package models

type Grade struct {
	ID           int     `json:"id"`
	EnrollmentID int     `json:"enrollmentid"`
	Score        float64 `json:"score"`
}
