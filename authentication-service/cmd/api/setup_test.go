package main

import (
	"authentication/data"
	"os"
	"testing"
)

var testApp Config

// go test -v ./
func TestMain(m *testing.M) {
	repo := data.NewPostgresTestRepository(nil)
	testApp.Repo = repo

	os.Exit(m.Run())
}
