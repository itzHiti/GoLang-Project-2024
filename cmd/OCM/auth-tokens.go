package main

import (
	"OCM/pkg/OCM/model"
	"OCM/pkg/OCM/validator"
	"errors"
	"net/http"
	"time"
)

func (app *application) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
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
	v := validator.New()
	model.ValidatePasswordPlaintext(v, input.Password)
	model.ValidateEmailOrUsername(v, input.Username, input.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	v1 := validator.New()
	model.ValidateUsername(v1, input.Username)
	usernameValid := v1.Valid()
	v2 := validator.New()
	model.ValidateEmail(v2, input.Email)
	emailValid := v2.Valid()
	if !usernameValid && !emailValid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user *model.User
	if emailValid {
		user, err = app.models.Users.GetByEmail(input.Email)
	} else {
		user, err = app.models.Users.GetByUsername(input.Username)
	}
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	role, err := app.models.Users.GetRole(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	accessToken, err := app.auth.GenerateAccessToken(user, role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	refreshToken, err := app.auth.GenerateRefreshToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	verification := &model.Verification{
		UserID:    user.ID,
		PlainText: accessToken,
		Expiry:    time.Now().Add(time.Hour * 24),
	}

	verification, err = app.models.Verifications.New(verification.UserID, time.Hour*24)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"access_token": accessToken, "refresh_token": refreshToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
}
