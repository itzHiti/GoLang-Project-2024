package main

import (
	"OCM/pkg/OCM/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
		CourseDuration string `json:"courseduration"`
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
		CourseDuration *string `json:"courseduration"`
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
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error (Probably course was not created)")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user page
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

func (app *application) listCoursesHandlerWithOutFilters(w http.ResponseWriter, r *http.Request) {
	courses, err := app.models.Courses.AllList()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, courses)
}

func (app *application) listAssignmentsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // default value
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10 // default value
	}

	validSortColumns := map[string]string{
		"title_asc":  "title ASC",
		"title_desc": "title DESC",
		"id_asc":     "id ASC",
		"id_desc":    "id DESC",
	}

	sortColumn, ok := validSortColumns[sort]
	if !ok {
		sortColumn = "id ASC" // default sorting
	}

	assignments, err := app.models.Assignments.List(page, pageSize, filter, sortColumn)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err))
		return
	}

	app.respondWithJSON(w, http.StatusOK, assignments)
}

func (app *application) listAssignmnetsWithoutFilters(w http.ResponseWriter, r *http.Request) {
	assignments, err := app.models.Assignments.AllAssignments()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, assignments)
}

func (app *application) listStudentsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // default value
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10 // default value
	}

	validSortColumns := map[string]string{
		"name_asc":  "name ASC",
		"name_desc": "name DESC",
		"age_asc":   "age ASC",
		"age_desc":  "age DESC",
		"gpa_asc":   "gpa ASC",
		"gpa_desc":  "gpa DESC",
	}

	sortColumn, ok := validSortColumns[sort]
	if !ok {
		sortColumn = "studentid ASC" // default sorting
	}

	students, err := app.models.Student.List(page, pageSize, filter, sortColumn)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err))
		return
	}

	app.respondWithJSON(w, http.StatusOK, students)
}

func (app *application) AssignmentsById(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CourseId    int    `json:"courseid"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	assignment := &model.Assignment{
		Title:       input.Title,
		Description: input.Description,
		CourseId:    input.CourseId,
	}

	err = app.models.Assignments.InsertAssignment(assignment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, assignment)
}

func (app *application) listAssignmentsByCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	assign, err := app.models.Assignments.FetchAssignmentsByCourse(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, assign)
}

func (app *application) AssignmentUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Assignment ID")
		return
	}

	assignment, err := app.models.Assignments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		CourseId    int     `json:"courseid"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		assignment.Title = *input.Title
	}

	if input.Description != nil {
		assignment.Description = *input.Description
	}

	if input.CourseId != 0 {
		assignment.CourseId = input.CourseId
	}

	err = app.models.Assignments.Update(assignment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, assignment)
}

func (app *application) AssigmentDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	err = app.models.Assignments.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error (Probably course was not created)")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) listStudentsByCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	students, err := app.models.Student.FetchStudentsByCourse(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, students)
}

func (app *application) createStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		GPA  float64 `json:"gpa"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	student := &model.Student{
		Name: input.Name,
		Age:  input.Age,
		GPA:  input.GPA,
	}

	err = app.models.Student.Insert(student)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, student)
}

func (app *application) getStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	student, err := app.models.Student.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, student)
}
func (app *application) updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	var input struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		GPA  float64 `json:"gpa"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name == "" || input.Age <= 0 || input.GPA < 0.0 || input.GPA > 4.0 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid student data")
		return
	}

	student := &model.Student{
		StudentID: id,
		Name:      input.Name,
		Age:       input.Age,
		GPA:       input.GPA,
	}

	err = app.models.Student.Update(student)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err))
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Student updated successfully"})
}

func (app *application) deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	err = app.models.Student.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.respondWithError(w, http.StatusNotFound, "Student not found")
		} else {
			app.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err))
		}
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}
