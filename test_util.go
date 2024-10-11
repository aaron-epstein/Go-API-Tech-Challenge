package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, r *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func applyVars(src *string, tctx TestContext) {
	for k, v := range tctx.Vars {
		*src = strings.ReplaceAll(*src, "{"+k+"}", v)
	}
}

func executeTests(tctx TestContext, tests []UnitTest) error {
	for i, test := range tests {
		fmt.Printf("TEST %d: %v", i+1, test)
		tctx.Test = &test
		err := executeTest(tctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeTest(tctx TestContext) error {

	t := tctx.T
	test := tctx.Test
	applyVars(&test.Url, tctx)
	applyVars(&test.Body, tctx)

	var req *http.Request
	if test.Body != "" {
		body := bytes.NewReader([]byte(test.Body))
		req, _ = http.NewRequest(test.Method, test.Url, body)
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(test.Method, test.Url, nil)
	}
	res := executeRequest(req, tctx.R)

	require.Equal(t, test.Status, res.Code)
	if test.ResponseFn != nil {
		err := test.ResponseFn(tctx, res)
		if err != nil {
			fmt.Println("Error in response function:", err)
			require.Nil(t, err)
			return err
		}
	}
	return nil
}

type UnitTest struct {
	Method     string
	Url        string
	Body       string
	Status     int
	ResponseFn func(TestContext, *httptest.ResponseRecorder) error
}

func (t UnitTest) String() string {
	return fmt.Sprintf(`{
  Method: %v
  Url: %v
  Body: %v
  Status: %d
}
`, t.Method, t.Url, t.Body, t.Status)
}

type TestContext struct {
	T    *testing.T
	R    *chi.Mux
	Vars map[string]string
	Test *UnitTest
}

func NewTestContext(t *testing.T, r *chi.Mux) TestContext {
	tctx := TestContext{
		T: t,
		R: r,
	}
	tctx.Vars = make(map[string]string)
	return tctx
}
