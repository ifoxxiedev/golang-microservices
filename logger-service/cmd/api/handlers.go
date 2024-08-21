package main

import (
	"log-service/cmd/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.LogEntry.Insert(data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	})

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	app.writeJson(w, http.StatusOK, jsonResponse{
		Error:   false,
		Message: "logged",
	})
}
