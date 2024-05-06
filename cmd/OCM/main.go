package main

import (
	"OCM/pkg/OCM/model"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env string
	db struct {
	dsn string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
	}
	limiter struct {
	enabled bool
	rps float64
	burst int
	}
	smtp struct {
	host string
	port int
	username string
	password string
	sender string
	}
	cors struct {
	trustedOrigins []string
	}
	jwt struct {
	secret string // Add a new field to store the JWT signing secret.
	}
}
	

type application struct {
	config config
	models model.Models
}

func (app *application) run() {
	r := mux.NewRouter()

	r.HandleFunc("/", app.HomeHandler)     // home page
	r.HandleFunc("/user", app.UserHandler) // user page

	// Course Singleton
	// r.HandleFunc("/courses/{id}", app.getCourseHandler).Methods("GET")       // [LEGACY] Get a specific course
	r.HandleFunc("/courses", app.listCoursesHandler).Methods("GET")
	r.HandleFunc("/courses", app.createCourseHandler).Methods("POST")        // Create a new course
	r.HandleFunc("/courses/{id}", app.updateCourseHandler).Methods("PUT")    // Update a specific course
	r.HandleFunc("/courses/{id}", app.deleteCourseHandler).Methods("DELETE") // Delete a specific course
	r.HandleFunc("/register", app.registerUserHandler).Methods("POST")
	r.HandleFunc("/activate", app.activateUserHandler).Methods("GET") // Using JWT tokens

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234Asdf@localhost:5432/go?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()

}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
