package model

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Student struct {
	StudentID int     `json:"studentid"`
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

type StudentModel struct {
	DB *sql.DB
}

func (sm *StudentModel) List(page, pageSize int, filter, sort string) ([]*Student, error) {
	offset := (page - 1) * pageSize
	query := "SELECT studentid, name, age, gpa FROM student WHERE name ILIKE ?"

	if sort != "" {
		query += " ORDER BY " + sort
	} else {
		query += " ORDER BY studentid ASC" // default sorting
	}
	query += " LIMIT ? OFFSET ?"

	rows, err := sm.DB.Query(query, "%"+strings.ToLower(filter)+"%", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.StudentID, &student.Name, &student.Age, &student.GPA); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	return students, nil
}

func (sm *StudentModel) Get(id int) (*Student, error) {
	query := `
        SELECT studentid, name, age, gpa
        FROM student
        WHERE studentid = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	student := &Student{}
	err := sm.DB.QueryRowContext(ctx, query, id).Scan(&student.StudentID, &student.Name, &student.Age, &student.GPA)
	if err != nil { // nil => null
		if err == sql.ErrNoRows {
			return nil, errors.New("students not found")
		} else {
			// Some other error happened
			return nil, err
		}
	}

	return student, nil
}
func (sm *StudentModel) Insert(student *Student) error {

	query := `
		INSERT INTO student (studentid, name, age, gpa) 
		VALUES ($1, $2, $3, $4) 
		RETURNING studentid
	`
	args := []interface{}{student.StudentID, student.Name, student.Age, student.GPA}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sm.DB.QueryRowContext(ctx, query, args...).Scan(&student.StudentID)
}

func (sm *StudentModel) Update(student *Student) error {

	query := `
        UPDATE students
        SET name = $1, age = $2, gpa = $3
        WHERE studentid = $4
        RETURNING studentid
    `
	args := []interface{}{student.Name, student.Age, student.GPA, student.StudentID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sm.DB.QueryRowContext(ctx, query, args...).Scan(&student.StudentID)
}

func (sm *StudentModel) Delete(id int) error {

	query := `
        DELETE FROM student
        WHERE studentid = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := sm.DB.ExecContext(ctx, query, id)
	return err
}
func (sm *StudentModel) FetchStudentsByCourse(courseID int) ([]Student, error) {
	stmt := `SELECT s.studentid, s.name, s.age, s.gpa
			 FROM student s
			 JOIN student_course sc ON s.studentid = sc.studentid
			 JOIN course c ON sc.courseid = c.courseid
			 WHERE c.courseid = $1`

	rows, err := sm.DB.Query(stmt, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.StudentID, &student.Name, &student.Age, &student.GPA)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}
