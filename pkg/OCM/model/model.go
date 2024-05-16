package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Courses       CourseModel
	Users         UserModel
	Assignments   AssignmentModel
	Verifications VerificationModel
	Roles         RoleModel
	Student       StudentModel
}
type StudentCourse struct {
	StudentID int `json:"studentid"`
	CourseID  int `json:"courseid"`
}
type CourseModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Courses: CourseModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Users: UserModel{
			DB: db,
		},
		Assignments: AssignmentModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Verifications: VerificationModel{
			DB: db,
		},
		Roles: RoleModel{
			DB: db,
		},
		Student: StudentModel{
			DB: db,
		},
	}
}
