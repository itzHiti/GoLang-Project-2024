package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	_ "os"
	"strings"
	"time"
)

type Assignment struct {
	AssignmentId int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CourseId     int    `json:"courseid"`
}

type AssignmentModel struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CourseID    int    `json:"courseid"`
	DB          *sql.DB
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
}

func (am *AssignmentModel) List(page, pageSize int, filter, sort string) ([]*Assignment, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, title, description, courseid FROM assignmentmodel WHERE title ILIKE $1 ORDER BY " + sort + " LIMIT $2 OFFSET $3"

	log.Printf("Executing query: %s with params: filter=%s, pageSize=%d, offset=%d", query, filter, pageSize, offset)

	rows, err := am.DB.Query(query, "%"+strings.ToLower(filter)+"%", pageSize, offset)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var assignments []*Assignment
	for rows.Next() {
		var assignment Assignment
		if err := rows.Scan(&assignment.AssignmentId, &assignment.Title, &assignment.Description, &assignment.CourseId); err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func (am *AssignmentModel) AllAssignments() ([]*Assignment, error) {
	var assignments []*Assignment
	baseQuery := `SELECT id, title, description, courseid FROM assignmentmodel`
	rows, err := am.DB.Query(baseQuery)
	if err != nil {
		return nil, err // Properly return the error if the query execution fails
	}
	defer rows.Close()

	for rows.Next() {
		var assignment Assignment
		// Scanning each row into a Course struct
		if err := rows.Scan(&assignment.AssignmentId, &assignment.Title, &assignment.Description, &assignment.CourseId); err != nil {
			return nil, err // Return an error if any occurs during row scanning
		}
		assignments = append(assignments, &assignment) // Append each course to the slice
	}

	// Check for any error that occurred during the iteration over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil // Return the slice of courses
}

func (am *AssignmentModel) InsertAssignment(assignment *Assignment) error {
	query := `
		INSERT INTO assignmentmodel (title, description, courseid) 
		VALUES ($1, $2, $3) 
		RETURNING courseid
		`
	args := []interface{}{assignment.Title, assignment.Description, assignment.CourseId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&assignment.CourseId)
}

func (am *AssignmentModel) Update(assignment *Assignment) error {
	query := `
        UPDATE assignmentmodel
        SET title = $1, description = $2, id = $3
        WHERE id = $4
        RETURNING id
        `
	args := []interface{}{assignment.Title, assignment.Description, assignment.AssignmentId, assignment.CourseId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&assignment.AssignmentId)
}

func (am *AssignmentModel) FetchAssignmentsByCourse(courseId int) ([]Assignment, error) {
	query := `
    SELECT 
        a.id, a.title, a.description, a.courseid
    FROM 
        assignmentmodel a
    JOIN 
        course c ON a.courseid = c.courseid
    WHERE 
        c.courseid = $1
    `

	rows, err := am.DB.Query(query, courseId)
	if err != nil {
		am.ErrorLog.Printf("Error fetching assignments for course ID %d: %v", courseId, err)
		return nil, err
	}
	defer rows.Close()

	var assignments []Assignment

	for rows.Next() {
		var assignment Assignment
		if err := rows.Scan(
			&assignment.AssignmentId,
			&assignment.Title,
			&assignment.Description,
			&assignment.CourseId,
		); err != nil {
			am.ErrorLog.Printf("Error scanning assignment: %v", err)
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	if err = rows.Err(); err != nil {
		am.ErrorLog.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return assignments, nil
}

func (am *AssignmentModel) Get(id int) (*Assignment, error) {
	// Query the course from the database.
	query := `
        SELECT id, title, description, courseid
        FROM assignmentmodel
        WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	assignment := &Assignment{}
	err := am.DB.QueryRowContext(ctx, query, id).Scan(&assignment.CourseId, &assignment.Title, &assignment.Description, &assignment.AssignmentId)
	if err != nil { // nil => null
		if err == sql.ErrNoRows {
			// The course was not found
			return nil, errors.New("courses not found")
		} else {
			// Some other error happened
			return nil, err
		}
	}

	return assignment, nil
}
func (am *AssignmentModel) Delete(id int) error {
	// Delete a specific course from the database.
	query := `
        DELETE FROM assignmentmodel
        WHERE id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := am.DB.ExecContext(ctx, query, id)
	return err
}
