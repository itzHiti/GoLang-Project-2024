package main

import (
	"OCM/pkg/OCM/auth"
	"OCM/pkg/OCM/mailer"
	"OCM/pkg/OCM/model"
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	_ "github.com/joho/godotenv"
)

const version = "1.5.2"

type config struct {
	port string
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *log.Logger
	models model.Models
	mailer mailer.Mailer
	auth   auth.AuthService
}

//func (app *application) run() {
//	r := mux.NewRouter()
//
//	r.HandleFunc("/", app.HomeHandler)     // home page
//	r.HandleFunc("/user", app.UserHandler) // user page
//
//	// Course Singleton
//	// r.HandleFunc("/courses/{id}", app.getCourseHandler).Methods("GET")       // [LEGACY] Get a specific course
//
//	r.HandleFunc("/courses", app.requireActivatedUser(app.listCoursesHandlerWithOutFilters)).Methods("GET")
//
//	r.HandleFunc("/courses", app.createCourseHandler).Methods("POST")        // Create a new course
//	r.HandleFunc("/courses/{id}", app.updateCourseHandler).Methods("PUT")    // Update a specific course
//	r.HandleFunc("/courses/{id}", app.deleteCourseHandler).Methods("DELETE") // Delete a specific course
//
//	r.HandleFunc("/users", app.registerUserHandler).Methods("POST")
//	r.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
//	//Authenticate new user
//	r.HandleFunc("/login", app.createAuthTokenHandler).Methods("POST")
//
//	log.Println("Starting server on ", app.config.port)
//	err := http.ListenAndServe(app.config.port, r)
//	log.Fatal(err)
//}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:qwe@localhost:5432/go?sslmode=disable", "PostgreSQL DSN")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "SFw6DlXYh4B4SM75hwf6cqvzgF30e5SKPSYt0hVXHCBMnOM8lRmI4EQm5hIqdRfIL4kG4VANPMQqQjHImXwbNg==", "JWT secret")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      app.authenticate(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s. Version: %s", cfg.env, srv.Addr, version)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func autoMigrate(db *sql.DB) error {
	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance("file:///OCM/pkg/OCM/migrations", "go", migrationDriver)

	if err != nil {
		return err
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
