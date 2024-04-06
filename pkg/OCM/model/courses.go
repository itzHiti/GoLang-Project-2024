package model

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Course struct {
	CourseId       int    `json:"course_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CourseDuration string `json:"courseDuration"`
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
        SELECT course_id, title, description, course_duration
        FROM courses
        WHERE course_id = $1
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
		INSERT INTO courses (title, description, course_duration) 
		VALUES ($1, $2, $3) 
		RETURNING course_id
		`
	args := []interface{}{course.Title, course.Description, course.CourseDuration}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return cm.DB.QueryRowContext(ctx, query, args...).Scan(&course.CourseId)
}

func (cm *CourseModel) Update(course *Course) error {
	// Update a specific course in the database.
	query := `
        UPDATE courses
        SET title = $1, description = $2, course_duration = $3
        WHERE course_id = $4
        RETURNING course_id
        `
	args := []interface{}{course.Title, course.Description, course.CourseDuration, course.CourseId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return cm.DB.QueryRowContext(ctx, query, args...).Scan(&course.CourseId)
}

func (cm *CourseModel) Delete(id int) error {
	// Delete a specific course from the database.
	query := `
        DELETE FROM courses
        WHERE course_id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.ExecContext(ctx, query, id)
	return err
}

func (cm *CourseModel) List(page, pageSize int, filter, sort string) ([]*Course, error) {
	var courses []*Course

	baseQuery := `SELECT course_id, title, description, course_duration FROM courses`
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