package model

import (
	"context"
	"database/sql"
	"log"
	_ "os"
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
