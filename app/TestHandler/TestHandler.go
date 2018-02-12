package TestHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galactic-filament/go-home/app/RouteHandler"
	"github.com/galactic-filament/go-home/app/Util"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// TestHandler - context for running tests
type TestHandler struct {
	Db *sqlx.DB
	T  *testing.T
}

func (th TestHandler) testRequest(method string, dest string, body io.Reader, status int) *httptest.ResponseRecorder {
	// fetching the request router
	r := RouteHandler.GetHandler(th.Db)

	// generating a request to test it
	req, err := http.NewRequest(method, dest, body)
	assert.Nil(th.T, err, fmt.Sprintf("Could not create new %s %s request", method, dest))

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// checking for 500 errors
	if w.Code == http.StatusInternalServerError {
		var errResponse Util.ErrorResponse
		err = json.NewDecoder(w.Body).Decode(&errResponse)
		assert.Nil(th.T, err, "Could not decode response body")
		assert.NotNil(th.T, nil, fmt.Sprintf("Response code was 500: %s", errResponse.Error))
		return w
	}

	// asserting that it worked properly
	assert.Equal(th.T, status, w.Code, fmt.Sprintf("Response code was not %d", http.StatusOK))
	return w
}

// TestGetRequest - generates a test request and runs it
func (th TestHandler) TestGetRequest(dest string, status int) *httptest.ResponseRecorder {
	return th.testRequest("GET", dest, nil, status)
}

// TestGetJSONRequest - generates a test json GET request and runs it
func (th TestHandler) TestGetJSONRequest(dest string, status int) *httptest.ResponseRecorder {
	w := th.testRequest("GET", dest, nil, status)
	assert.Equal(th.T, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

// TestDeleteJSONRequest - generates a test json DELETE request and runs it
func (th TestHandler) TestDeleteJSONRequest(dest string, status int) *httptest.ResponseRecorder {
	w := th.testRequest("DELETE", dest, nil, status)
	assert.Equal(th.T, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

// TestPostJSONRequest - generates a test json POST request and runs it
func (th TestHandler) TestPostJSONRequest(dest string, payload io.Reader, status int) *httptest.ResponseRecorder {
	w := th.testRequest("POST", dest, payload, status)
	assert.Equal(th.T, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

// TestPutJSONRequest - generates a test json PUT request and runs it
func (th TestHandler) TestPutJSONRequest(dest string, payload io.Reader, status int) *httptest.ResponseRecorder {
	w := th.testRequest("PUT", dest, payload, status)
	assert.Equal(th.T, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}
