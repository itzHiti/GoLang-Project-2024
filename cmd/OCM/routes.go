package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// Perm
	r.HandleFunc("/courses", app.requireActivatedUser(app.listCoursesHandlerWithOutFilters)).Methods("GET")

	r.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	r.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

	//Authenticate new user
	r.HandleFunc("/login", app.createAuthTokenHandler).Methods("POST")

	return r
}
