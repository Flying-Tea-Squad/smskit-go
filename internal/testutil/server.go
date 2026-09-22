package testutil

import (
	"net/http"
	"net/http/httptest"
)

// Server returns a local HTTP server and its request recorder.
func Server(handler http.Handler) (*httptest.Server, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	return httptest.NewServer(handler), recorder
}
