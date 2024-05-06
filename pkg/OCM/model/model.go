package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
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

type Password struct {
	plaintext *string // открытый текст пароля
	hash      *string // хэш пароля
}

type UserModel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	DB        *sql.DB
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
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
