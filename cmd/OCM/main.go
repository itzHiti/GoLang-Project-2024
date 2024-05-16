package main

import (
	"OCM/pkg/OCM/auth"
	"OCM/pkg/OCM/mailer"
	"OCM/pkg/OCM/model"
	"context"
	"database/sql"
	"flag"
	"fmt"
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

func main() {
	err := godotenv.Load()
	var cfg config
	fmt.Println(os.Getenv("SMTP_PASSWORD"))
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234Asdf@localhost:5432/go?sslmode=disable", "PostgreSQL DSN")
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

	handler := corsMiddleware(app.authenticate(app.routes()))

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      handler,
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware: %s %s", r.Method, r.RequestURI)

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
