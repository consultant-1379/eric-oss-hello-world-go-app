package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	log "gerrit-review.gic.ericsson.se/cloud-ran/src/golang-log-api/logapi"
)

// CaCertPath for loading cert
var CaCertPath = getCertPath()

// HandleFormRequest for Client Credential Flow Login
func HandleFormRequest(endpoint string, formData url.Values, headers http.Header) ([]byte, error) {
	// Create a new TLS config with the server's CA cert
	tlsConfig := newTLSConfig()

	// Create an HTTP client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create a new http.Request object
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Create http.Request object failed: %w", err)
	}

	// Set the headers on the request
	if headers != nil {
		req.Header = headers
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request to the specified endpoint
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request Failed with following error: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading response body failed: %w", err)
	}

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &httpError{
			statusCode: resp.StatusCode,
			statusText: resp.Status,
			body:       respBody,
		}
	}

	// If the response body is empty, return nil
	if len(respBody) == 0 {
		return nil, nil
	}

	// Return the response body
	return respBody, nil
}

func newTLSConfig() *tls.Config {
	// Load the root CA certificate
	caCert, err := os.ReadFile(CaCertPath)
	if err != nil {
		log.Error("Failed to read root CA certificate: %v", err)
		return nil
	}

	// Create a new CertPool and add the root CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS config with the client certificates and the server's CA cert
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{},
		RootCAs:            caCertPool,
	}

	return tlsConfig
}

type httpError struct {
	statusCode int
	statusText string
	body       []byte
}

func (e *httpError) Error() string {
	return e.statusText
}

// combines CaMountPath and CaCertFileName as a full path
func getCertPath() (certFilePath string) {
	certFilePath = path.Join(config.CaMountPath, config.CaCertFileName)
	return certFilePath
}
