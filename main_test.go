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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func handleCourse() func(TestContext, *httptest.ResponseRecorder) error {
	return handleCourseFn(func(tctx TestContext, course internal.Course) error { return nil })
}

func handleCourseFn(fn func(TestContext, internal.Course) error) func(TestContext, *httptest.ResponseRecorder) error {
	return func(tctx TestContext, res *httptest.ResponseRecorder) error {
		fmt.Println("Start course", res.Body.String())
		var course internal.Course
		err := json.NewDecoder(bytes.NewReader(res.Body.Bytes())).Decode(&course)
		fmt.Println("Course:", course)
		if err != nil {
			require.Nil(tctx.T, err)
			return err
		}
		if fn != nil {
			err = fn(tctx, course)
			if err != nil {
				require.Nil(tctx.T, err)
				return err
			}
		}

		return nil
	}
}

func handleCourses() func(TestContext, *httptest.ResponseRecorder) error {
	return handleCoursesFn(func(tctx TestContext, courses []internal.Course) error { return nil })
}

func handleCoursesFn(fn func(TestContext, []internal.Course) error) func(TestContext, *httptest.ResponseRecorder) error {
	return func(tctx TestContext, res *httptest.ResponseRecorder) error {
		var courses []internal.Course
		err := json.NewDecoder(bytes.NewReader(res.Body.Bytes())).Decode(&courses)
		if err != nil {
			require.Nil(tctx.T, err)
			return err
		}
		if fn != nil {
			err = fn(tctx, courses)
			if err != nil {
				require.Nil(tctx.T, err)
				return err
			}
		}

		return nil
	}
}

func handlePerson() func(TestContext, *httptest.ResponseRecorder) error {
	return handlePersonFn(func(tctx TestContext, course internal.PersonResponse) error { return nil })
}

func handlePersonFn(fn func(TestContext, internal.PersonResponse) error) func(TestContext, *httptest.ResponseRecorder) error {
	return func(tctx TestContext, res *httptest.ResponseRecorder) error {
		var person internal.PersonResponse
		err := json.NewDecoder(bytes.NewReader(res.Body.Bytes())).Decode(&person)
		if err != nil {
			require.Nil(tctx.T, err)
			return err
		}
		if fn != nil {
			err = fn(tctx, person)
			if err != nil {
				require.Nil(tctx.T, err)
				return err
			}
		}

		return nil
	}
}

func handlePersons() func(TestContext, *httptest.ResponseRecorder) error {
	return handlePersonsFn(func(tctx TestContext, persons []internal.PersonResponse) error { return nil })
}

func handlePersonsFn(fn func(TestContext, []internal.PersonResponse) error) func(TestContext, *httptest.ResponseRecorder) error {
	return func(tctx TestContext, res *httptest.ResponseRecorder) error {
		var persons []internal.PersonResponse
		err := json.NewDecoder(bytes.NewReader(res.Body.Bytes())).Decode(&persons)
		if err != nil {
			require.Nil(tctx.T, err)
			return err
		}
		if fn != nil {
			err = fn(tctx, persons)
			if err != nil {
				require.Nil(tctx.T, err)
				return err
			}
		}

		return nil
	}
}

func testCourses(tctx TestContext) {

	tests := []UnitTest{
		{Method: "GET", Url: "/api/course", Status: http.StatusOK, ResponseFn: handleCourses()},
		{Method: "GET", Url: "/api/course/1", Status: http.StatusOK, ResponseFn: handleCourse()},
		{Method: "POST", Url: "/api/course", Status: http.StatusCreated, Body: `
    {
      "name": "Test User"
    }`, ResponseFn: handleCourseFn(func(tctx TestContext, course internal.Course) error {
			id := fmt.Sprintf("%d", course.ID)
			tctx.Vars["id"] = id
			return nil
		})},
		{Method: "PUT", Url: "/api/course/{id}", Status: http.StatusAccepted, Body: `
    {
      "name": "Test User Modified"
    }`, ResponseFn: handleCourse()},
		{Method: "DELETE", Url: "/api/course/{id}", Status: http.StatusOK},
	}

	executeTests(tctx, tests)
}

func testPersons(tctx TestContext) {

	tctx.Vars["name"] = "Test User"

	tests := []UnitTest{
		{Method: "GET", Url: "/api/person", Status: http.StatusOK, ResponseFn: handlePersons()},
		{Method: "GET", Url: "/api/person/Bill Gates", Status: http.StatusOK, ResponseFn: handlePerson()},
		{Method: "POST", Url: "/api/person", Status: http.StatusCreated, Body: `
    {
      "first_name": "Test",
      "last_name": "User",
      "type": "student",
      "age": 24,
      "courses": [1, 2]
		}`, ResponseFn: handlePerson()},
		{Method: "PUT", Url: "/api/person/{name}", Status: http.StatusAccepted, Body: `
    {
      "first_name": "Test",
      "last_name": "User",
      "type": "professor",
      "age": 25,
      "courses": [2, 3]
    }`, ResponseFn: handlePerson()},
		{Method: "DELETE", Url: "/api/person/{name}", Status: http.StatusOK},
	}

	executeTests(tctx, tests)
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

	tctx := NewTestContext(t, r)
	testCourses(tctx)
	testPersons(tctx)

}
