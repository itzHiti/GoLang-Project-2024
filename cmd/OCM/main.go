package main

import (
	"OCM/pkg/OCM/model"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func (app *application) run() {
	r := mux.NewRouter()

	r.HandleFunc("/", app.HomeHandler)           // home page
	r.HandleFunc("/courses", app.CoursesHandler) // course page
	r.HandleFunc("/user", app.UserHandler)       // user page

	// Course Singleton
	r.HandleFunc("/courses/{id}", app.getCourseHandler).Methods("GET")       // Get a specific course
	r.HandleFunc("/courses", app.createCourseHandler).Methods("POST")        // Create a new course
	r.HandleFunc("/courses/{id}", app.updateCourseHandler).Methods("PUT")    // Update a specific course
	r.HandleFunc("/courses/{id}", app.deleteCourseHandler).Methods("DELETE") // Delete a specific course

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
