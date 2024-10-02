package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaron-epstein/Go-API-Tech-Challenge/internal"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, r *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func testCourses(t *testing.T, r *chi.Mux) {

	var req *http.Request
	var response *httptest.ResponseRecorder
	var payload *bytes.Reader

	req, _ = http.NewRequest("GET", "/api/course", nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/course/1", nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)

	payload = bytes.NewReader([]byte(`{
		"name": "Test User"
	}`))
	req, _ = http.NewRequest("POST", "/api/course", payload)
	req.Header.Add("Content-Type", "application/json")
	response = executeRequest(req, r)
	require.Equal(t, http.StatusCreated, response.Code)
	var course internal.Course
	json.NewDecoder(bytes.NewReader(response.Body.Bytes())).Decode(&course)
	id := fmt.Sprintf("%d", course.ID)

	payload = bytes.NewReader([]byte(`{
		"name": "Test User Modified"
	}`))
	req, _ = http.NewRequest("PUT", "/api/course/"+id, payload)
	req.Header.Add("Content-Type", "application/json")
	response = executeRequest(req, r)
	require.Equal(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/course/"+id, nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)
}

func testPersons(t *testing.T, r *chi.Mux) {

	var req *http.Request
	var response *httptest.ResponseRecorder
	var payload *bytes.Reader

	name := "Test User"

	req, _ = http.NewRequest("GET", "/api/person", nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/person/Bill Gates", nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)

	payload = bytes.NewReader([]byte(`{
		"first_name": "Test",
		"last_name": "User",
		"type": "student",
		"age": 24,
		"courses": [1, 2]
	}`))
	req, _ = http.NewRequest("POST", "/api/person", payload)
	req.Header.Add("Content-Type", "application/json")
	response = executeRequest(req, r)
	require.Equal(t, http.StatusCreated, response.Code)

	payload = bytes.NewReader([]byte(`{
		"first_name": "Test",
		"last_name": "User",
		"type": "professor",
		"age": 25,
		"courses": [2, 3]
	}`))
	req, _ = http.NewRequest("PUT", "/api/person/"+name, payload)
	req.Header.Add("Content-Type", "application/json")
	response = executeRequest(req, r)
	require.Equal(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/person/"+name, nil)
	response = executeRequest(req, r)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestMain(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = internal.InitDB()
	if err != nil {
		// panic(err)
		log.Fatal("Error connecting to DB")
	}

	r := internal.InitServer()

	testCourses(t, r)
	testPersons(t, r)

}
