package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_Routes_exist(t *testing.T) {
	testApp = Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router) // force casting to chi.Router

	routes := []string{"/authenticate"}

	for _, r := range routes {
		routeExists(t, chiRoutes, r)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	walkFn := func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route {
			found = true
		}
		return nil
	}

	_ = chi.Walk(routes, walkFn)

	if !found {
		t.Errorf("route \"%s \" did not found in registered routes ", route)
	}
}
