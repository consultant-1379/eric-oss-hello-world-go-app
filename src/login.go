package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Token retrieved here
type Token struct {
	AccessToken string `json:"accessToken"`
}

const loginPath = "/auth/realms/master/protocol/openid-connect/token"

// HandleLogin Create an instance of the request body
func HandleLogin(clientID, clientSecret, baseURL string) error {
	loginURL := baseURL + path.Join(loginPath)
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("tenant_id", "master")

	respBody, err := HandleFormRequest(loginURL, formData, http.Header{})
	if err != nil {
		return err
	}
	var token Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return fmt.Errorf("JSON Unmarshal Failed with following error: %w", err)
	}

	return nil
}
