package model

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	Courses     CourseModel
	Users       UserModel
	Assignments AssignmentModel
}

type CourseModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type UserModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type AssignmentModel struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CourseID    int    `json:"course_id"`
	DB          *sql.DB
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
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
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Assignments: AssignmentModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}

func (u *UserModel) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)
	return nil
}

func (u *UserModel) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}
