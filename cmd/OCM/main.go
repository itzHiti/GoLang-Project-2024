package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)           // home page
	r.HandleFunc("/courses", CoursesHandler) // course page
	r.HandleFunc("/user", UserHandler)       // user page

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Handle home page
}

func CoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Handle courses page
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user page
}
