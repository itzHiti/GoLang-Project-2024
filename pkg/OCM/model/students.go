package model

import (
	"context"
	_ "database/sql"
	"errors"
	"time"
)

type Student struct {
	StudentID int     `json:"student_id"`
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	GPA       float64 `json:"gpa"`
}

var students = []Student{
	{
		StudentID: 1,
		Name:      "John Doe",
		Age:       20,
		GPA:       3.8,
	},
	{
		StudentID: 2,
		Name:      "Jane Smith",
		Age:       22,
		GPA:       3.9,
	},
	{
		StudentID: 3,
		Name:      "Alice Johnson",
		Age:       21,
		GPA:       3.7,
	},
}

func GetStudents() []Student {
	return students
}

func (sm *UserModel) Get(id int) (*Student, error) {
	for _, s := range students {
		if s.StudentID == id {
			return &s, nil
		}
	}
	return nil, errors.New("Student not found")
}

func (sm *UserModel) Insert(student *Student) error {

	query := `
		INSERT INTO students (student_id, name, age, gpa) 
		VALUES ($1, $2, $3, $4) 
		RETURNING student_id
	`
	args := []interface{}{student.StudentID, student.Name, student.Age, student.GPA}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sm.DB.QueryRowContext(ctx, query, args...).Scan(&student.StudentID)
}

func (sm *UserModel) Update(student *Student) error {

	query := `
        UPDATE students
        SET name = $1, age = $2, gpa = $3
        WHERE student_id = $4
        RETURNING student_id
    `
	args := []interface{}{student.Name, student.Age, student.GPA, student.StudentID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sm.DB.QueryRowContext(ctx, query, args...).Scan(&student.StudentID)
}

func (sm *UserModel) Delete(id int) error {

	query := `
        DELETE FROM students
        WHERE student_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := sm.DB.ExecContext(ctx, query, id)
	return err
}
