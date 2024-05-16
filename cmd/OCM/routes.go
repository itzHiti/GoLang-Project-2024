package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// Healthcheck
	r.HandleFunc("/api/healthcheck", app.healthcheckHandler).Methods("GET")

	// Courses
	r.HandleFunc("/courses", app.requireActivatedUser(app.listCoursesHandlerWithOutFilters)).Methods("GET")
	r.HandleFunc("/courses/{id}", app.requireRole([]string{"admin"}, app.updateCourseHandler)).Methods("PUT")
	r.HandleFunc("/courses", app.createCourseHandler).Methods("POST")
	r.HandleFunc("/courses/{id}", app.requireRole([]string{"admin"}, app.deleteCourseHandler)).Methods("DELETE")
	// Courses - filter/pagination/sort
	r.HandleFunc("/coursess", app.listCoursesHandler).Methods("GET")

	// Combined
	r.HandleFunc("/courses/{id}/assignments", app.listAssignmentsByCourse).Methods("GET")
	r.HandleFunc("/courses/{id:[0-9]+}/students", app.listStudentsByCourse).Methods("GET")

	// Assignments
	r.HandleFunc("/assignments", app.requireActivatedUser(app.listAssignmnetsWithoutFilters)).Methods("GET")
	r.HandleFunc("/assignments", app.AssignmentsById).Methods("POST")
	r.HandleFunc("/assignments/{id}", app.AssignmentUpdate).Methods("PUT")
	r.HandleFunc("/assignments/{id}", app.AssigmentDelete).Methods("DELETE")
	// Assignments - filter/pagination/sort
	r.HandleFunc("/assignmentss", app.listAssignmentsHandler).Methods("GET")
	// Student
	r.HandleFunc("/students/{id}", app.getStudentHandler).Methods("GET")
	r.HandleFunc("/students", app.createStudentHandler).Methods("POST")
	// Student - filter/pagination/sort
	r.HandleFunc("/studentss", app.listStudentsHandler).Methods("GET")
	// user auth
	r.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	//Authenticate new user
	r.HandleFunc("/login", app.createAuthTokenHandler).Methods("POST")

	return r
}
