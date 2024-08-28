// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"os"
	"strconv"
	"strings"
)

// Config todo de-export
type Config struct {
	LocalPort           int
	LocalProtocol       string
	CertFile            string
	KeyFile             string
	LogControlFile      string
	LogEndpoint         string
	LogTLSKey           string
	LogTLSCert          string
	LogTLSCACert        string
	IamClientID         string
	IamClientSecret     string
	IamBaseURL          string
	CaCertFileName      string
	CaMountPath         string
	logCaCertFilePath   string
	rAppLogCertFilePath string
}

var instance *Config

const localPort = 8050

func getConfig() *Config {
	if instance == nil {
		instance = &Config{
			LocalPort:           getOsEnvInt("LOCAL_PORT", localPort),
			LocalProtocol:       getOsEnvString("LOCAL_PROTOCOL", "http"),
			CertFile:            getOsEnvString("CERT_FILE", "certificate.pem"),
			KeyFile:             getOsEnvString("KEY_FILE", "key.pem"),
			LogControlFile:      getOsEnvString("LOG_CTRL_FILE", ""),
			LogEndpoint:         getOsEnvString("LOG_ENDPOINT", ""),
			LogTLSKey:           getOsEnvString("APP_LOG_TLS_KEY", ""),
			LogTLSCert:          getOsEnvString("APP_LOG_TLS_CERT", ""),
			LogTLSCACert:        getOsEnvString("LOG_TLS_CA_CERT", ""),
			IamClientID:         getOsEnvString("IAM_CLIENT_ID", ""),
			IamClientSecret:     getOsEnvString("IAM_CLIENT_SECRET", ""),
			IamBaseURL:          getOsEnvString("IAM_BASE_URL", ""),
			CaCertFileName:      getOsEnvString("CA_CERT_FILENAME", ""),
			CaMountPath:         getOsEnvString("CA_CERT_MOUNT_PATH", ""),
			logCaCertFilePath:   getOsEnvString("LOG_CA_CERT_FILE_PATH", ""),
			rAppLogCertFilePath: getOsEnvString("APP_LOG_CERT_FILE_PATH", ""),
		}
	}

	return instance
}

func getOsEnvInt(envName string, defaultValue int) (result int) {
	envValue := strings.TrimSpace(os.Getenv(envName))
	result, err := strconv.Atoi(envValue)
	if err != nil {
		result = defaultValue
	}

	return
}

func getOsEnvString(envName string, defaultValue string) (result string) {
	result = strings.TrimSpace(os.Getenv(envName))

	if result == "" {
		result = defaultValue
	}

	return
}
