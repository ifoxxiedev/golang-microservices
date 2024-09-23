package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(req http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(*req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func Test_Authenticate_Handler(t *testing.T) {
	jsonToReturn := `
{
  "error": false,
  "message": "some message"
}
`

	client := NewTestClient(func(req http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header:     make(http.Header),
		}
	})

	testApp.Client = client

	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "password",
	}

	body, _ := json.Marshal(postBody)
	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.authenticate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusAccepted)
	}
}
