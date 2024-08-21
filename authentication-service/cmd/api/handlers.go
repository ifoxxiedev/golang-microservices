package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logger in user %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, responsePayload)
}
