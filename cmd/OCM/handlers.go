package main

import (
	"OCM/pkg/OCM/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}
func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Handle home page
}

func (app *application) createCourseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		CourseDuration string `json:"courseDuration"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	course := &model.Course{
		Title:          input.Title,
		Description:    input.Description,
		CourseDuration: input.CourseDuration,
	}

	err = app.models.Courses.Insert(course)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, course)
}

func (app *application) getCourseHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	course, err := app.models.Courses.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, course)
}

func (app *application) updateCourseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	course, err := app.models.Courses.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title          *string `json:"title"`
		Description    *string `json:"description"`
		CourseDuration *string `json:"courseDuration"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		course.Title = *input.Title
	}

	if input.Description != nil {
		course.Description = *input.Description
	}

	if input.CourseDuration != nil {
		course.CourseDuration = *input.CourseDuration
	}

	err = app.models.Courses.Update(course)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, course)
}

func (app *application) deleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	err = app.models.Courses.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user page
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) listCoursesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // def value
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10 // def value
	}

	courses, err := app.models.Courses.List(page, pageSize, filter, sort)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, courses)
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var u model.UserModel
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := app.models.Users.Register(&u); err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}
	if err := app.models.Users.ActivateUser(token); err != nil {
		http.Error(w, "Activation failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
