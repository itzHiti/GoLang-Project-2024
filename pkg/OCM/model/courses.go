package model

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Course struct {
	CourseId       int    `json:"courseid"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CourseDuration string `json:"courseduration"`
}

var courses = []Course{
	{
		CourseId:       1,
		Title:          "Programming Principles I",
		Description:    "Starter programming courses",
		CourseDuration: "3 month",
	}, {
		CourseId:       2,
		Title:          "Programming Principles II",
		Description:    "Advanced programming courses",
		CourseDuration: "4 month",
	}, {
		CourseId:       3,
		Title:          "Algorithms and Data Structures",
		Description:    "This course for student who accepted their fate",
		CourseDuration: "3 month",
	}, {
		CourseId:       4,
		Title:          "Computer Networks",
		Description:    "This course will provide good topic for understanding world web",
		CourseDuration: "5 month",
	},
}

func GetCourses() []Course {
	return courses
}

func (cm *CourseModel) Get(id int) (*Course, error) {
	// Query the course from the database.
	query := `
        SELECT courseid, title, description, courseduration
        FROM course
        WHERE courseid = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	course := &Course{}
	err := cm.DB.QueryRowContext(ctx, query, id).Scan(&course.CourseId, &course.Title, &course.Description, &course.CourseDuration)
	if err != nil { // nil => null
		if err == sql.ErrNoRows {
			// The course was not found
			return nil, errors.New("courses not found")
		} else {
			// Some other error happened
			return nil, err
		}
	}

	return course, nil
}

func (cm *CourseModel) Insert(course *Course) error {
	// Insert a new course into the database.
	query := `
		INSERT INTO course (title, description, courseduration) 
		VALUES ($1, $2, $3) 
		RETURNING courseid
		`
	args := []interface{}{course.Title, course.Description, course.CourseDuration}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return cm.DB.QueryRowContext(ctx, query, args...).Scan(&course.CourseId)
}

func (cm *CourseModel) Update(course *Course) error {
	// Update a specific course in the database.
	query := `
        UPDATE course
        SET title = $1, description = $2, courseduration = $3
        WHERE courseid = $4
        RETURNING courseid
        `
	args := []interface{}{course.Title, course.Description, course.CourseDuration, course.CourseId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return cm.DB.QueryRowContext(ctx, query, args...).Scan(&course.CourseId)
}

func (cm *CourseModel) Delete(id int) error {
	// Delete a specific course from the database.
	query := `
        DELETE FROM course
        WHERE courseid = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.ExecContext(ctx, query, id)
	return err
}

func (cm *CourseModel) List(page, pageSize int, filter, sort string) ([]*Course, error) {
	var courses []*Course

	baseQuery := `SELECT courseid, title, description, courseduration FROM course`
	whereClauses, args := []string{}, []interface{}{}

	// Фильтрация
	if filter != "" {
		whereClauses = append(whereClauses, "title ILIKE $1")
		args = append(args, "%"+filter+"%")
	}

	// Добавляем WHERE только если есть условия фильтрации
	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Сортировка
	orderBy := " ORDER BY course_id ASC" // default sort by course_id in ascending order
	if sort != "" {
		switch sort {
		case "title_asc":
			orderBy = " ORDER BY title ASC"
		case "title_desc":
			orderBy = " ORDER BY title DESC"
		case "duration_asc":
			orderBy = " ORDER BY course_duration ASC"
		case "duration_desc":
			orderBy = " ORDER BY course_duration DESC"
		}
	}

	// Пагинация
	pagination := " LIMIT $2 OFFSET $3"
	args = append(args, pageSize, (page-1)*pageSize)

	finalQuery := baseQuery + orderBy + pagination

	rows, err := cm.DB.Query(finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.CourseId, &course.Title, &course.Description, &course.CourseDuration); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (cm *CourseModel) AllList() ([]*Course, error) {
	var courses []*Course
	baseQuery := `SELECT courseid, title, description, courseduration FROM course`
	rows, err := cm.DB.Query(baseQuery)
	if err != nil {
		return nil, err // Properly return the error if the query execution fails
	}
	defer rows.Close() // Ensure we close the rows to free up resources

	for rows.Next() {
		var course Course
		// Scanning each row into a Course struct
		if err := rows.Scan(&course.CourseId, &course.Title, &course.Description, &course.CourseDuration); err != nil {
			return nil, err // Return an error if any occurs during row scanning
		}
		courses = append(courses, &course) // Append each course to the slice
	}

	// Check for any error that occurred during the iteration over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil // Return the slice of courses
}
