package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) error {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		return app.errorJSON(w, err, http.StatusBadRequest)
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		return app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		return app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	return app.writeJSON(w, http.StatusAccepted, payload)
}
