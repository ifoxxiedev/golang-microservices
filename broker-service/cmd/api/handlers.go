package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"os"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct {
	Name string `json:"level"`
	Data string `json:"message"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omityempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
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
	case "log":
		app.logItemViaRpc(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	// create some json we'll send to mail microservice
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// call the sevice
	url := fmt.Sprintf("%s/send", os.Getenv("MAIL_SERVICE_URL"))

	request, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + m.To

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	// create some json we'll send to auth microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	// call the sevice
	url := fmt.Sprintf("%s/log", os.Getenv("LOGGER_SERVICE_URL"))

	request, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJson(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	app.writeJson(w, http.StatusAccepted, payload)
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

	request.Header.Set("Content-Type", "application/json")

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

func (app *Config) logItemViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(
		l.Name,
		l.Data,
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged via RabbitMQ"

	app.writeJson(w, http.StatusAccepted, payload)
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRpc(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	fmt.Println("LOG VIA RPC", l)

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string

	// Trick (Call the same struct and method)
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	return emitter.Push(string(j), "log.INFO")
}
