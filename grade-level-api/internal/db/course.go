// Package db
package db

import (
	"database/sql"
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

func UpdateCourse(id int, name string, level, lecturerID int) (*models.Course, error) {
	if name == "" {
		return nil, errors.New("course name cannot be empty")
	}
	if level <= 0 {
		return nil, errors.New("course level must be a positive integer")
	}
	result, err := DB.Exec(
		`UPDATE course SET name = ?, level = ?, lecturer_id = ? WHERE id = ?`, name, level, lecturerID, id,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("no course found with the given ID")
	}
	return &models.Course{
		ID:         id,
		Name:       name,
		Level:      level,
		LecturerID: lecturerID,
	}, nil
}

func ListCourses() ([]*models.Course, error) {
	rows, err := DB.Query(`SELECT id, name, level, lecturer_id FROM course`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []*models.Course

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Level, &course.LecturerID); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return courses, nil
}

func FindCourseByID(id int) (*models.Course, error) {
	var course models.Course
	err := DB.QueryRow(
		`SELECT id, name, level, lecturer_id FROM course WHERE id = ?`, id,
	).Scan(&course.ID, &course.Name, &course.Level, &course.LecturerID)
	if err == sql.ErrNoRows {
		return nil, errors.New("no course found with the given ID")
	}
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func DeleteCourse(id int) error {
	if id <= 0 {
		return errors.New("invalid course ID")
	}
	result, err := DB.Exec(`DELETE FROM course WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no course found with the given ID")
	}
	return nil
}

func FindCoursesByLecturerID(lecturerID int) ([]*models.Course, error) {
	if lecturerID <= 0 {
		return nil, errors.New("invalid lecturer ID")
	}

	rows, err := DB.Query(
		`SELECT id, name, level, lecturer_id FROM course WHERE lecturer_id = ?`,
		lecturerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course

	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Level, &c.LecturerID); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func FindCoursesByLevel(level int) ([]*models.Course, error) {
	if level <= 0 {
		return nil, errors.New("invalid course level")
	}

	rows, err := DB.Query(
		`SELECT id, name, level, lecturer_id FROM course WHERE level = ?`,
		level,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course

	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Level, &c.LecturerID); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func FindCoursesByLecturerAndLevel(lecturerID, level int) ([]*models.Course, error) {
	if lecturerID <= 0 || level <= 0 {
		return nil, errors.New("invalid lecturer ID or level")
	}

	rows, err := DB.Query(
		`SELECT id, name, level, lecturer_id 
		 FROM course 
		 WHERE lecturer_id = ? AND level = ?`,
		lecturerID, level,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course

	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Level, &c.LecturerID); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
