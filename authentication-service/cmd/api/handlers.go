package main

import "net/http"

// Oww craaazy, this is a JSON payload
type requestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Oww craaazy, this is a JSON payload
type response struct {
	Token string `json:"token"`
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &payload)

	if err != nil {

	}

	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {

	}

	matches, err := app.Models.User.PasswordMatches(payload.Password)
	if err != nil {

	}

	// TODO - Generate JWT token
}
