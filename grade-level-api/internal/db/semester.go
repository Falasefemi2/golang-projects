// Package db
package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/falasefemi2/gradesystem/internal/models"
)

type Semester string

const (
	FirstSemster   Semester = "firstsemster"
	SecondSemester Semester = "secondsemester"
	ThirdSemester  Semester = "thirdsemester"
)

func CreateSemester(name Semester, startDate, endDate time.Time) (*models.Semester, error) {
	if name == "" {
		return nil, errors.New("semester name cannot be empty")
	}
	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	result, err := DB.Exec("INSERT INTO semester (name, start_date, end_date) VALUES (?, ?, ?)", string(name), startDate, endDate)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Semester{
		ID:        int(id),
		Name:      string(name),
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func GetAllSemesters() ([]models.Semester, error) {
	rows, err := DB.Query(
		`SELECT id, name, start_date, end_date FROM semester ORDER BY start_date`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var semesters []models.Semester

	for rows.Next() {
		var s models.Semester
		if err := rows.Scan(&s.ID, &s.Name, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		semesters = append(semesters, s)
	}

	return semesters, nil
}

func GetSemesterByID(id int) (*models.Semester, error) {
	var s models.Semester

	err := DB.QueryRow(
		`SELECT id, name, start_date, end_date FROM semester WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.Name, &s.StartDate, &s.EndDate)

	if err == sql.ErrNoRows {
		return nil, errors.New("semester not found")
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func UpdateSemester(id int, name Semester, startDate, endDate time.Time) (*models.Semester, error) {
	if id <= 0 {
		return nil, errors.New("invalid semester id")
	}
	if name == "" {
		return nil, errors.New("semester name cannot be empty")
	}
	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	result, err := DB.Exec(
		`UPDATE semester 
		 SET name = ?, start_date = ?, end_date = ?
		 WHERE id = ?`,
		string(name),
		startDate,
		endDate,
		id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("semester not found")
	}

	return &models.Semester{
		ID:        id,
		Name:      string(name),
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func DeleteSemester(id int) error {
	if id <= 0 {
		return errors.New("invalid semester id")
	}

	result, err := DB.Exec(`DELETE FROM semester WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("semester not found")
	}

	return nil
}
