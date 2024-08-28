package main_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	api "eric-oss-hello-world-go-app/src"

	"github.com/stretchr/testify/assert"
)

var formData = getFormData()

func TestHandleFormRequestWithResponse(t *testing.T) {
	// same response should be returned if server is sending response
	t.Parallel()
	testResponse := []byte(`test`)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(testResponse)
	}))
	defer server.Close()
	resp, _ := api.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Equal(t, resp, testResponse)
}

func TestHandleFormRequestWithNoResponse(t *testing.T) {
	// if no response returned, the len of byte returned is 0
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	defer server.Close()
	resp, _ := api.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Equal(t, len(resp), 0)
}

func TestHandleFormRequestWithIncorrectStatusCode(t *testing.T) {
	// when receiving 403 response, error message should contain 403
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	// configure the test certificate file
	api.CaCertPath = "test.crt"
	f, _ := os.Create("test.crt")
	defer f.Close()

	_, err := api.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Contains(t, err.Error(), "403")

	// clear certificate file
	os.Remove("test.crt")
}

func getFormData() url.Values {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", "ClientID")
	formData.Set("client_secret", "ClientSecret")
	formData.Set("tenant_id", "master")
	return formData
}
