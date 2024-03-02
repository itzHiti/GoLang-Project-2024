package model

import (
	"context"
	"database/sql"
	"errors"
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
			return nil, errors.New("Courses not Found")
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
		INSERT INTO courses (course_id, title, description, course_duration) 
		VALUES ($1, $2, $3, $4) 
		RETURNING course_id
		`
	args := []interface{}{course.CourseId, course.Title, course.Description, course.CourseDuration}
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
