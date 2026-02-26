package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GET_Routes_V1(t *testing.T) {
	app := NewAPI(nil, nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Health check",
			path:        "/api/v1/health",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"ok\":true}\n",
		},
		{
			description: "Get non-existent record",
			path:        "/api/v1/records/15",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"record of id 15 does not exist\"}\n",
		},
		{
			description: "Get record with negative id",
			path:        "/api/v1/records/-1",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Invalid id (non-numeric)",
			path:        "/api/v1/records/abc",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
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
		body        string // needed for POST requests only; can be empty for GET requests
		path        string
		wantStatus  int
		PostResBody string
		GetResBody  string
	}{
		{
			description: "Post to negative id",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/-11",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid id; id must be a positive number\"}\n",
			GetResBody:  "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Post to invalid id (non-numeric)",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/abc",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid id; id must be a positive number\"}\n",
			GetResBody:  "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Post invalid json body",
			body:        "[{\"key1\":}]",
			path:        "/api/v1/records/18",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid input; could not parse json\"}\n",
			GetResBody:  "{\"error\":\"record of id 18 does not exist\"}\n",
		},
		{
			description: "Post new record",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			PostResBody: "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
			GetResBody:  "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			// the body has no key2 because we create a new router for each test case
			// so the record created in the first test case is not persisted in the second test case for v1 api
			PostResBody: "{\"id\":1,\"data\":{\"key1\":\"value2\",\"status\":\"ok\"}}\n",
			GetResBody:  "{\"id\":1,\"data\":{\"key1\":\"value2\",\"status\":\"ok\"}}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			app := NewAPI(nil, nil)
			router := app.SetupRouter()

			jsonBody := []byte(tc.body)
			body := bytes.NewBuffer(jsonBody)

			reqPost := httptest.NewRequest("POST", tc.path, body)
			reqPost.Header.Set("Content-Type", "application/json")
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.PostResBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			// and returns same exact body as the POST request
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.GetResBody, rrGet.Body.String())
		})
	}
}

func Test_Updates_V1(t *testing.T) {
	app := NewAPI(nil, nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		body        string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Post new record",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v1/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value2\",\"key2\":\"222\",\"status\":\"ok\"}}\n",
		},
		{
			description: "Delete field in existing record",
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

			reqPost := httptest.NewRequest("POST", tc.path, body)
			reqPost.Header.Set("Content-Type", "application/json")
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.wantBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			// and returns same exact body as the POST request
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.wantBody, rrGet.Body.String())
		})
	}
}

// Benchmark_POST_Routes_V1-12            1        1036914167 ns/op        98356320 B/op     580426 allocs/op
func Benchmark_POST_Routes_V1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	app := NewAPI(nil, nil)
	router := app.SetupRouter()

	var rrPost *httptest.ResponseRecorder

	for n := 0; n < 10000; n++ {
		bodyStr := "{" + "\"key" + strconv.Itoa(n) + "\":null," + "\"key" + strconv.Itoa(n+1) + "\":\"value" + strconv.Itoa(n+1) + "\"" + "}"
		jsonBody := []byte(bodyStr)
		body := bytes.NewBuffer(jsonBody)

		reqPost := httptest.NewRequest("POST", "/api/v1/records/33", body)
		reqPost.Header.Set("Content-Type", "application/json")
		rrPost = httptest.NewRecorder()
		router.ServeHTTP(rrPost, reqPost)
	}

	require.Equal(b, http.StatusOK, rrPost.Code)
	time.Sleep(1 * time.Second)

	reqGet := httptest.NewRequest("GET", "/api/v1/records/33", nil)
	rrGet := httptest.NewRecorder()
	router.ServeHTTP(rrGet, reqGet)

	require.Equal(b, http.StatusOK, rrGet.Code)
	require.Equal(b, "{\"id\":33,\"data\":{\"key10000\":\"value10000\"}}\n", rrGet.Body.String())
}

func Test_GET_Routes_V2(t *testing.T) {
	app := NewAPI(nil, nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Health check",
			path:        "/api/v2/health",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"ok\":true}\n",
		},
		{
			description: "Get non-existent record",
			path:        "/api/v2/records/15",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"record of id 15 does not exist\"}\n",
		},
		{
			description: "Get record with negative id",
			path:        "/api/v2/records/-1",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Invalid id (non-numeric)",
			path:        "/api/v2/records/abc",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tc.wantStatus, rr.Code)
			require.Equal(t, tc.wantBody, rr.Body.String())
		})
	}
}

func Test_POST_Routes_V2(t *testing.T) {
	tests := []struct {
		description string
		body        string // needed for POST requests only; can be empty for GET requests
		path        string
		wantStatus  int
		PostResBody string
		GetResBody  string
	}{
		{
			description: "Post to negative id",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v2/records/-11",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid id; id must be a positive number\"}\n",
			GetResBody:  "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Post to invalid id (non-numeric)",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v2/records/abc",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid id; id must be a positive number\"}\n",
			GetResBody:  "{\"error\":\"invalid id; id must be a positive number\"}\n",
		},
		{
			description: "Post invalid json body",
			body:        "[{\"key1\":}]",
			path:        "/api/v2/records/18",
			wantStatus:  http.StatusBadRequest,
			PostResBody: "{\"error\":\"invalid input; could not parse json\"}\n",
			GetResBody:  "{\"error\":\"record of id 18 does not exist\"}\n",
		},
		{
			description: "Post new record",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v2/records/1",
			wantStatus:  http.StatusOK,
			PostResBody: "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
			GetResBody:  "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v2/records/1",
			wantStatus:  http.StatusOK,
			// the body has no key2 because we create a new router for each test case
			// so the record created in the first test case is not persisted in the second test case for v1 api
			PostResBody: "{\"id\":1,\"data\":{\"key1\":\"value2\",\"status\":\"ok\"}}\n",
			GetResBody:  "{\"id\":1,\"data\":{\"key1\":\"value2\",\"status\":\"ok\"}}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			app := NewAPI(nil, nil)
			router := app.SetupRouter()

			jsonBody := []byte(tc.body)
			body := bytes.NewBuffer(jsonBody)

			reqPost := httptest.NewRequest("POST", tc.path, body)
			reqPost.Header.Set("Content-Type", "application/json")
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.PostResBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			// and returns same exact body as the POST request
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.GetResBody, rrGet.Body.String())
		})
	}
}

func Test_Updates_V2(t *testing.T) {
	app := NewAPI(nil, nil)
	router := app.SetupRouter()

	tests := []struct {
		description string
		body        string
		path        string
		wantStatus  int
		wantBody    string
	}{
		{
			description: "Post new record",
			body:        "{\"key1\":\"value1\",\"key2\":\"222\"}",
			path:        "/api/v2/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value1\",\"key2\":\"222\"}}\n",
		},
		{
			description: "Update existing record",
			body:        "{\"key1\":\"value2\",\"status\":\"ok\"}",
			path:        "/api/v2/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key1\":\"value2\",\"key2\":\"222\",\"status\":\"ok\"}}\n",
		},
		{
			description: "Delete field in existing record",
			body:        "{\"key1\":null,\"status\":null}",
			path:        "/api/v2/records/1",
			wantStatus:  http.StatusOK,
			wantBody:    "{\"id\":1,\"data\":{\"key2\":\"222\"}}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			jsonBody := []byte(tc.body)
			body := bytes.NewBuffer(jsonBody)

			reqPost := httptest.NewRequest("POST", tc.path, body)
			reqPost.Header.Set("Content-Type", "application/json")
			rrPost := httptest.NewRecorder()
			router.ServeHTTP(rrPost, reqPost)

			require.Equal(t, tc.wantStatus, rrPost.Code)
			require.Equal(t, tc.wantBody, rrPost.Body.String())

			// every POST request is followed by a GET request to ensure that the record was actually created
			// and returns same exact body as the POST request
			reqGet := httptest.NewRequest("GET", tc.path, nil)
			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			require.Equal(t, tc.wantStatus, rrGet.Code)
			require.Equal(t, tc.wantBody, rrGet.Body.String())
		})
	}
}
