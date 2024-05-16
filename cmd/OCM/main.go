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

	"github.com/joho/godotenv"

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
	err := godotenv.Load()
	var cfg config

	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234Asdf@localhost:5432/go?sslmode=disable", "PostgreSQL DSN")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "OCM <no-reply@OCM.net>", "SMTP sender")

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
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      app.authenticate(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("database migrations applied")

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
	path := "C:/Users/akaza/OneDrive/Рабочий стол/visual/OCM/pkg/OCM/migrations/"

	migrator, err := migrate.NewWithDatabaseInstance("file://"+path, "go", migrationDriver)

	if err != nil {
		return err
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
