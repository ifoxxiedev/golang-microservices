package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omityempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJson(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the sevice
	url := fmt.Sprintf("%s/authenticate", os.Getenv("AUTH_SERVICE_URL"))

	request, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, err, response.StatusCode)
		return
	}

	// create a variable we'll read response.Body into jsonResponse
	var jsonFromService jsonResponse

	// decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}
