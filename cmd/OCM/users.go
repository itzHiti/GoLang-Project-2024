package main

import (
	"OCM/pkg/OCM/model"
	"OCM/pkg/OCM/validator"
	"errors"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		Username:  input.Username,
		Email:     input.Email,
		Activated: false,
		TokenHash: app.auth.GenerateRandomString(15),
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, model.ErrDuplicateUsername):
			v.AddError("username", "a user with this username already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)

	var output struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	output.Username = user.Username
	output.Email = user.Email

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": output}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
