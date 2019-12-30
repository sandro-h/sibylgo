package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// MockSimpleJSONResponse sets up a test http server which responds with the passed JSON response
// for all requests.
func MockSimpleJSONResponse(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}
