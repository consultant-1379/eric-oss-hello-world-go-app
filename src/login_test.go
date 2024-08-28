package main_test

import (
	"encoding/json"
	api "eric-oss-hello-world-go-app/src"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleLoginWithInvalidURL(t *testing.T) {
	// when invalid baseURL is passed, we are expecting an error
	t.Parallel()
	err := api.HandleLogin("testID", "testSecret", "")
	assert.Contains(t, err.Error(), "Request Failed")
}

func TestHandleLoginWithJSON(t *testing.T) {
	// when server returns token, there is no error
	t.Parallel()
	testResponse := api.Token{AccessToken: "testToken"}
	jsonResp, _ := json.Marshal(testResponse)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(jsonResp)
	}))
	defer server.Close()

	err := api.HandleLogin("testID", "testSecret", server.URL)
	assert.Nil(t, err, "err should be nill")
}

func TestHandleLoginWithoutJSON(t *testing.T) {
	// JSON Unmarshal error when server does not return JSON
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	err := api.HandleLogin("testID", "testSecret", server.URL)
	assert.Contains(t, err.Error(), "JSON Unmarshal Failed")
}
