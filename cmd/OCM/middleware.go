package main

import (
	"OCM/pkg/OCM/model"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) authenticate(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]

		if r.URL.Path == "/token/refresh" {
			userRefresh, err := app.auth.ValidateRefreshToken(token)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			user, err := app.models.Users.GetByUsername(userRefresh.Username)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			actualCustomKey := app.auth.GenerateCustomKey(user.Username, user.TokenHash)
			if userRefresh.CustomKey != actualCustomKey {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			r = app.contextSetUser(r, user)

		} else {
			userAccess, err := app.auth.ValidateAccessToken(token)
			if err != nil {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			var user model.User
			user.ID = userAccess.UserId
			user.Username = userAccess.Username
			user.Email = userAccess.Email
			user.Activated = userAccess.Activated
			user.Role = userAccess.Role

			r = app.contextSetUser(r, &user)

		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requireRole(allowedRoles model.Roles, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		role := user.Role
		fmt.Println(role)
		fmt.Println(user)
		if !allowedRoles.Include(role) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}
