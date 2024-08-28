// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main_test

import (
	ipUtils "eric-oss-hello-world-go-app/src"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpAddressWithOnlyRemoteAddr(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = "192.0.0.1:8080"

	// act
	ipAddress := ipUtils.GetIPInfo(request)

	// assert
	assert.NotNil(t, ipAddress, "IP address should not be nill")
	assert.Contains(t, ipAddress, "RemoteAddr: '192.0.0.1:8080'")
}

func TestIpAddressWithOnlyXForwordedForHeaderSet(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = ""
	request.Header.Set("X-Forwarded-For", "192.0.0.1,192.0.0.2")

	// act
	ipAddress := ipUtils.GetIPInfo(request)

	// assert
	assert.NotNil(t, ipAddress, "IP address should not be nill")
	assert.Contains(t, ipAddress, "X-Forwarded-For: '192.0.0.1,192.0.0.2'")
}

func TestIpAddressWithBothRemoteAddrAndXForwordedForHeader(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = "192.79.12.10:9090"
	request.Header.Set("X-Forwarded-For", "192.0.0.1,192.0.0.2")

	// act
	ipAddress := ipUtils.GetIPInfo(request)

	// assert
	assert.Contains(t, ipAddress, "X-Forwarded-For: '192.0.0.1,192.0.0.2', RemoteAddr: '192.79.12.10:9090'")
}

func TestIpAddressWithoutXForwordedForHeaderAndRemoteAddr(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.Header.Set("X-Forwarded-For", "")
	request.RemoteAddr = ""

	// act
	ipAddress := ipUtils.GetIPInfo(request)

	// assert
	assert.Equal(t, "X-Forwarded-For: '', RemoteAddr: ''", ipAddress)
}
