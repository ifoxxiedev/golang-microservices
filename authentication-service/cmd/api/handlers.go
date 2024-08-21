package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

	// log authentification
	err = app.logRequest("authentification", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logger in user %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, responsePayload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := fmt.Sprintf("%s/log", os.Getenv("LOGGER_SERVICE_URL"))

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	return nil
}
