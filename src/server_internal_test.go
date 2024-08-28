// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

// because they use the server with a port.
// the noctx does not make sense in tests here
//
//nolint:paralleltest,noctx // the tests here are not supposed to run in parallel
package main

import (
	"context"
	"eric-oss-hello-world-go-app/src/internal/metric"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	log "gerrit-review.gic.ericsson.se/cloud-ran/src/golang-log-api/logapi"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const logOutputFileName = "testlogfile"

func TestHelloAndHealthResponseAreValid(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		status       int
		response     string
		routeHandler func(resp http.ResponseWriter, req *http.Request)
	}{
		{
			name: "`/hello` Hello Endpoint handler", endpoint: "/hello", status: 200, response: "Hello World!!",
			routeHandler: handleAPICall,
		},
		{
			name: "`/health` Health Endpoint handler", endpoint: "/health", status: 200, response: "Ok",
			routeHandler: checkServerHealth,
		},
	}

	for _, testParameters := range tests {
		// arrange
		request := httptest.NewRequest(http.MethodGet, testParameters.endpoint, nil)
		response := httptest.NewRecorder()

		// act
		testParameters.routeHandler(response, request)

		// assert
		res := response.Result()
		defer res.Body.Close()
		assert.Equal(t, testParameters.status, res.StatusCode,
			"Status code should be 200, but got : "+strconv.Itoa(res.StatusCode))

		data, _ := io.ReadAll(res.Body)
		assert.NotNil(t, data, "Data should not be nill")

		assert.Equal(t, testParameters.response, string(data),
			fmt.Sprintf("Should be returned `%v` but got : %v", testParameters.response, string(data)))
	}
}

func TestGetHelloAndHealthEndPointReturnValidResponse(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		status       int
		response     string
		routeHandler func(resp http.ResponseWriter, req *http.Request)
	}{
		{
			name: "`/hello` route testing", endpoint: "/hello", status: 200, response: "Hello World!!",
			routeHandler: handleAPICall,
		},
		{
			name: "`/health` route testing", endpoint: "/health", status: 200, response: "Ok",
			routeHandler: checkServerHealth,
		},
	}

	for _, testParameters := range tests {
		// arrange
		router := http.NewServeMux()
		router.HandleFunc(testParameters.endpoint, testParameters.routeHandler)
		svr := httptest.NewServer(router)
		defer svr.Close()

		// act
		res, err := http.Get(fmt.Sprintf("%s%s", svr.URL, testParameters.endpoint))

		// assert
		assert.Nil(t, err, fmt.Sprintf("Could not send GET request:  %v", err))
		defer res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK,
			fmt.Sprintf("Expected Status Ok; got %v", res.Status))

		data, _ := io.ReadAll(res.Body)
		assert.Equal(t, testParameters.response, string(data),
			fmt.Sprintf("Expected response is `%v`, but got %v", testParameters.response, string(data)))
	}
}

func TestMetricsSuccessRequestCount(t *testing.T) {
	// arrange
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))
	router.HandleFunc("/hello", handleAPICall)

	svr := httptest.NewServer(router)
	defer svr.Close()

	// act
	res, err := http.Get(fmt.Sprintf("%s/metrics", svr.URL))

	// assert
	assert.Nil(t, err, fmt.Sprintf("Could not send GET request:  %v", err))
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("Expected Status Ok; got %v", res.Status))

	data, _ := io.ReadAll(res.Body)
	initialCount := currentRequestCount(string(data))

	for i := 0; i < 10; i++ {
		helloResponse, helloError := http.Get(fmt.Sprintf("%s/hello", svr.URL))
		if helloError != nil {
			defer helloResponse.Body.Close()
		}
	}

	res2, err2 := http.Get(fmt.Sprintf("%s/metrics", svr.URL))

	assert.Nil(t, err2, fmt.Sprintf("Could not send GET request:  %v", err2))

	defer res2.Body.Close()

	assert.Equal(t, res2.StatusCode, http.StatusOK, fmt.Sprintf("Expected Status Ok; got %v", res2.Status))

	data2, _ := io.ReadAll(res2.Body)
	totalRequestCount := currentRequestCount(string(data2))

	assert.Equal(t, true, strings.Contains(string(data2),
		"hello_world_requests_total"), "Response should be contains `hello_world_requests_total`")

	assert.Equal(t, initialCount+10, totalRequestCount,
		fmt.Sprintf("Request count should be %v, but got : %v", initialCount+10, totalRequestCount))

	assert.Equal(t, true, strings.Contains(string(data2),
		"hello_world_requests_failed_total 0"), "Response should be contains `hello_world_requests_failed_total 0`")
}

func TestLoggerLevelUpdate(t *testing.T) {
	// arrange
	log.Log().Log.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	defer file.Close()

	wrt := io.MultiWriter(os.Stdout, file)
	log.Log().Log.SetOutput(wrt)

	testServer := httptest.NewServer(http.HandlerFunc(handleAPICall))
	defer testServer.Close()

	helloResponse, helloError := http.Get(fmt.Sprintf("%s/hello", testServer.URL))
	if helloError != nil {
		defer helloResponse.Body.Close()
	}

	data, err := os.ReadFile(logOutputFileName)
	assert.Nil(t, err, fmt.Sprintf("error opening file: %v", err))

	assert.False(t, strings.Contains(string(data), "Leaving api handler..."), "This message should not be logged")

	// act
	log.Log().Log.SetLevel(logrus.DebugLevel)
	helloResponse2, helloError2 := http.Get(fmt.Sprintf("%s/hello", testServer.URL))
	if helloError2 != nil {
		defer helloResponse2.Body.Close()
	}

	// assert
	data2, err2 := os.ReadFile(logOutputFileName)
	assert.Nil(t, err2, fmt.Sprintf("error opening file: %v", err2))

	assert.True(t, strings.Contains(string(data2), "Leaving api handler..."), "This message should be logged")

	t.Cleanup(func() {
		e := os.Remove(logOutputFileName)
		assert.Nil(t, err2, fmt.Sprintf("error deleting test log file: %v", e))
	})
}

func TestSetLogger(t *testing.T) {
	// arrange
	ctx, servercancel := context.WithCancel(context.Background())
	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	defer file.Close()

	testfile, openFileError := os.OpenFile("test.file", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, openFileError, fmt.Sprintf("error creating file: %v", openFileError))
	defer testfile.Close()

	wrt := io.MultiWriter(os.Stdout, file)
	log.Log().Log.SetOutput(wrt)

	config = &Config{
		LocalPort:      getOsEnvInt("LOCAL_PORT", 8050),
		LocalProtocol:  getOsEnvString("LOCAL_PROTOCOL", "http"),
		CertFile:       getOsEnvString("CERT_FILE", "certificate.pem"),
		KeyFile:        getOsEnvString("KEY_FILE", "key.pem"),
		LogControlFile: getOsEnvString("LOG_CTRL_FILE", "test.file"),
		LogEndpoint:    getOsEnvString("LOG_ENDPOINT", ""),
		LogTLSKey:      getOsEnvString("APP_LOG_TLS_KEY", ""),
		LogTLSCert:     getOsEnvString("APP_LOG_TLS_CERT", ""),
		LogTLSCACert:   getOsEnvString("LOG_TLS_CA_CERT", ""),
	}

	initLogger(ctx)

	// assert
	data, err := os.ReadFile(logOutputFileName)
	assert.Nil(t, err, fmt.Sprintf("error opening file: %v", err))
	assert.NotNil(t, data, "Content should not be nill")

	assert.Equal(t, true, strings.Contains(string(data), "Logging has been enabled successfully..."))

	t.Cleanup(func() {
		e := os.Remove(logOutputFileName)
		assert.Nil(t, err, fmt.Sprintf("error deleting log test file: %v", e))

		err := os.Remove("test.file")
		assert.Nil(t, err, fmt.Sprintf("error deleting test file: %v", err))
	})
	servercancel()
}

func TestLogStreamingOverHttpProtocol(t *testing.T) {
	testAssert := assert.New(t)
	counter := 0
	ctx, servercancel := context.WithCancel(context.Background())
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	serverURL, _ := url.Parse(srv.URL)
	host, port, _ := net.SplitHostPort(serverURL.Host)
	defer srv.Close()

	config = &Config{
		LocalPort:      getOsEnvInt("LOCAL_PORT", 8050),
		LocalProtocol:  getOsEnvString("LOCAL_PROTOCOL", "http"),
		CertFile:       getOsEnvString("CERT_FILE", ""),
		KeyFile:        getOsEnvString("KEY_FILE", ""),
		LogControlFile: getOsEnvString("LOG_CTRL_FILE", ""),
		LogEndpoint:    getOsEnvString("LOG_ENDPOINT", host+":"+port),
		LogTLSKey:      getOsEnvString("APP_LOG_TLS_KEY", ""),
		LogTLSCert:     getOsEnvString("APP_LOG_TLS_CERT", ""),
		LogTLSCACert:   getOsEnvString("LOG_TLS_CA_CERT", ""),
	}

	initLogger(ctx)

	for i := 0; i < 10; i++ {
		log.Log().Info("REGULAR" + strconv.Itoa(i))
		log.Log().InfoSecAuth("Security" + strconv.Itoa(i))
	}

	time.Sleep(100 * time.Millisecond)
	t.Log("No. of log statements should match with no of logstream http request made")
	testAssert.Equal(21, counter)
	servercancel()
}

func TestMain(t *testing.T) {
	// update log level for main testing
	log.Log().Log.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	defer file.Close()

	wrt := io.MultiWriter(os.Stdout, file)
	log.Log().Log.SetOutput(wrt)

	// to stop application from waiting to receive interupt signal
	close(ExitSignal)
	// call main to start testing
	main()

	data, err := os.ReadFile(logOutputFileName)
	assert.Nil(t, err, fmt.Sprintf("error opening file: %v", err))

	assert.True(t, strings.Contains(string(data), "Server is ready to receive web requests"),
		"This message should be logged")

	// reset log level to debug
	log.Log().Log.SetLevel(logrus.DebugLevel)
}

func TestLogStreamingOverHttpsProtocol(t *testing.T) {
	ctx, servercancel := context.WithCancel(context.Background())

	testAssert := assert.New(t)
	counter := 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	serverURL, _ := url.Parse(srv.URL)
	host, port, _ := net.SplitHostPort(serverURL.Host)

	srv.Client()

	defer srv.Close()
	config.LogEndpoint = host + ":" + port

	initLogger(ctx)

	log.Log().Info("REGULAR" + strconv.Itoa(1))
	log.Log().InfoSecAuth("Security" + strconv.Itoa(1))

	t.Log("No. of log statements should match with no of logstream http request made")
	testAssert.Equal(0, counter)
	servercancel()
	t.Cleanup(
		func() {
			config.LogEndpoint = ""
			config.LogTLSKey = ""
			config.LogTLSCert = ""
			config.LogTLSCACert = ""
		})
}

func TestExitSignalChannel(t *testing.T) {
	// act
	channel := getExitSignalsChannel()

	// assert
	assert.NotNil(t, channel, "Channel should not be nill")
	assert.Equal(t, 1, cap(channel), "Capacity should be 1")
	assert.Equal(t, "chan", reflect.ValueOf(channel).Kind().String(), "Kind should be 'chan'")
}

func TestLoggingIpAddressWithValidAndInvalidIp(t *testing.T) {
	remoteAddrStr := "191.0.0.1:9090"
	ValidResponse := "RemoteAddr: '191.0.0.1:9090'"

	tests := []struct {
		name             string
		remoteAddrString string
		status           int
		response         string
	}{
		{name: "Valid IP Address", remoteAddrString: remoteAddrStr, status: 200, response: ValidResponse},
		{name: "Invalid IP Address", remoteAddrString: "", status: 200, response: "RemoteAddr: ''"},
	}

	for _, testCase := range tests {
		// arrange
		logOutputFile, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			t.Fatalf("error creating file: %v", err)
		}
		defer logOutputFile.Close()

		wrt := io.MultiWriter(os.Stdout, logOutputFile)
		log.Log().Log.SetOutput(wrt)

		request := httptest.NewRequest(http.MethodGet, "/hello", nil)
		request.RemoteAddr = testCase.remoteAddrString
		response := httptest.NewRecorder()

		// act
		handleAPICall(response, request)

		// assert
		data, err := os.ReadFile(logOutputFileName)
		if err != nil {
			t.Fatalf("error opening file: %v", err)
		}

		assert.Contains(t, string(data), testCase.response)

		if err := os.Truncate(logOutputFileName, 0); err != nil {
			t.Fatalf("Failed to remove the file data: %v", err)
		}
	}

	t.Cleanup(func() {
		e := os.Remove(logOutputFileName)
		if e != nil {
			t.Fatalf("error deleting test file: %v", e)
		}
	})
}

func TestStartWebService(t *testing.T) {
	retries := 3
	var res *http.Response
	var err error

	// act
	srv := startWebService()

	// assert
	assert.NotNil(t, srv, "Should not be nill")

	for retries > 0 {
		client := http.Client{}
		res, err = client.Get("http://localhost:8050/hello")
		if err != nil {
			log.Error("Request failed %v: ", err)
			retries--
		} else {
			defer res.Body.Close()
			break
		}
	}

	assert.NotNil(t, res, "Response should not be nill")
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("Expected Status Ok; got %v", res.Status))

	content, _ := io.ReadAll(res.Body)
	assert.NotNil(t, content, "Content should not be nill")
	assert.Equal(t, "Hello World!!", string(content), "Should be returned `Hello World!!` but got : "+string(content))

	ctxShutDown, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	srv.Shutdown(ctxShutDown)
}

func TestSampleAppArtifactsFileContents(t *testing.T) {
	errorMessage := "%v is missing in the artifact contents file"

	content, err := os.ReadFile("../zip-artifact-contents.txt")
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	contents := [9]string{
		"charts", "src", "vendor", "Dockerfile-template", "go.mod", "go.sum", "README.md",
		"version", "csar",
	}

	for _, contentName := range contents {
		assert.Contains(t, string(content), contentName, fmt.Sprintf(errorMessage, contentName))
	}
}

func TestErrorOnBadCertPath(t *testing.T) {
	ctx, servercancel := context.WithCancel(context.Background())
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	serverURL, _ := url.Parse(srv.URL)
	host, port, _ := net.SplitHostPort(serverURL.Host)
	defer srv.Close()

	// Providing test values so initLogger will attempt to run with TLS
	config.LogEndpoint = host + ":" + port
	config.LogTLSKey = "test.key"
	config.LogTLSCert = "test.cert"
	config.LogTLSCACert = "test.CACert"
	config.rAppLogCertFilePath = "/etc/tls/log/"
	config.logCaCertFilePath = "/etc/tls-ca/log/"
	TestCertPath := config.rAppLogCertFilePath + config.LogTLSCert

	// Assert that the Error returns with a Cert path matching our test Path
	defer func() {
		if err := recover(); err != nil {
			str := fmt.Sprintf("%v", err)
			assert.Contains(t, str, TestCertPath)
			servercancel()
		}
	}()
	initLogger(ctx)
}

func currentRequestCount(data string) int {
	start := strings.LastIndex(data, "hello_world_requests_total")
	start += len("hello_world_requests_total")
	end := start + 1
	for end < len(data) && (data[end] >= '0' && data[end] <= '9') {
		end++
	}
	intVal, _ := strconv.Atoi(strings.TrimSpace(data[start:end]))
	return intVal
}
