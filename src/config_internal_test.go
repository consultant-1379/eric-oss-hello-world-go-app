// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	instance := getConfig()

	assert.NotNil(t, instance,
		"Instance should not be nill")

	assert.Equal(t, 8050, instance.LocalPort,
		"Port should be 8050, but got : "+strconv.Itoa(instance.LocalPort))

	assert.Equal(t, "http", instance.LocalProtocol,
		"Protocol should be http, but got : "+instance.LocalProtocol)

	assert.Equal(t, "certificate.pem", instance.CertFile,
		"Certificate file name should be `certificate.pem`, but got : "+instance.CertFile)

	assert.Equal(t, "key.pem", instance.KeyFile,
		"Keyfile file name should be `key.pem`, but got : "+instance.KeyFile)
}
