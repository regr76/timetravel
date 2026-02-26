package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GET_Routes_V1(t *testing.T) {
	app := NewAPI(nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		method      string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Health check",
			method:      "GET",
			path:        "/api/v1/health",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"ok\":true}\n",
		},
		{
			description: "Get non-existent record",
			method:      "GET",
			path:        "/api/v1/records/15",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"record of id 15 does not exist\"}\n",
		},
		{
			description: "Get record with negative id",
			method:      "GET",
			path:        "/api/v1/records/-1",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Invalid id (non-numeric)",
			method:      "GET",
			path:        "/api/v1/records/abc",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tc.wantStatus, rr.Code)
			require.Equal(t, tc.wantBody, rr.Body.String())
		})
	}
}

func Test_POST_Routes_V1(t *testing.T) {
	tests := []struct {
		description string
		method      string
		body        string // needed for POST requests only; can be empty for GET requests
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Post new record",
			method:      "POST",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			method:      "POST",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			// the body has no key2 because we create a new router for each test case
			// so the record created in the first test case is not persisted in the second test case for v1 api
			wantBody: "{\"id\":1,\"data\":{\"key1\":\"value2\",\"status\":\"ok\"}}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			app := NewAPI(nil)
			router := app.SetupRouter()

			jsonBody := []byte(tc.body)
			body := bytes.NewBuffer(jsonBody)

			reqPost := httptest.NewRequest(tc.method, tc.path, body)
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.wantBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.wantBody, rrGet.Body.String())
		})
	}
}

func Test_Updates_V1(t *testing.T) {
	app := NewAPI(nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		method      string
		body        string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Post new record",
			method:      "POST",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			method:      "POST",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value2\",\"key2\":\"222\",\"status\":\"ok\"}}\n",
		},
		{
			description: "Delete field in existing record",
			method:      "POST",
			body:        "{\"key1\":null,\"status\":null}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key2\":\"222\"}}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			jsonBody := []byte(tc.body)
			body := bytes.NewBuffer(jsonBody)

			reqPost := httptest.NewRequest(tc.method, tc.path, body)
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.wantBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.wantBody, rrGet.Body.String())
		})
	}
}
