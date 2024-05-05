package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// Healthcheck
	r.HandleFunc("/api/healthcheck", app.healthcheckHandler).Methods("GET")

	// Perm Courses
	r.HandleFunc("/courses", app.requireActivatedUser(app.listCoursesHandlerWithOutFilters)).Methods("GET")
	r.HandleFunc("/courses/{id}", app.requireRole([]string{"admin"}, app.updateCourseHandler)).Methods("PUT")
	r.HandleFunc("/courses", app.createCourseHandler).Methods("POST")
	r.HandleFunc("/courses/{id}", app.requireRole([]string{"admin"}, app.deleteCourseHandler)).Methods("DELETE")

	r.HandleFunc("/courses/{id}/assignments", app.listAssignmentsByCourse).Methods("GET")

	r.HandleFunc("/assignments", app.requireActivatedUser(app.listAssignmnets)).Methods("GET")
	r.HandleFunc("/assignments", app.AssignmentsById).Methods("POST")

	r.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	//Authenticate new user
	r.HandleFunc("/login", app.createAuthTokenHandler).Methods("POST")

	return r
}
