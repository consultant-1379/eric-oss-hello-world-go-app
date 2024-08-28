// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"context"
	"eric-oss-hello-world-go-app/src/internal/metric"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "gerrit-review.gic.ericsson.se/cloud-ran/src/golang-log-api/logapi"
	"gerrit-review.gic.ericsson.se/cloud-ran/src/golang-tlsconf/tlsconf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config = getConfig()

	// ServiceID The service id using in the log fields
	ServiceID = "rapp-eric-oss-hello-world-go-app"
	// Version The version id
	Version = "0.0.1"
	server  *http.Server
	// ExitSignal The ExitSignal
	ExitSignal chan os.Signal
)

func handleAPICall(resp http.ResponseWriter, req *http.Request) {
	log.Debug("Entering api handler...")

	log.Info("Request IP: %s", GetIPInfo(req))

	metric.RequestsTotal.Inc()

	err := HandleLogin(config.IamClientID, config.IamClientSecret, config.IamBaseURL)
	if err != nil {
		log.Error("Login Failed. %v", err.Error())
	} else {
		log.Debug("Login Success.")
	}

	fmt.Fprintf(resp, "Hello World!!")

	log.Debug("Leaving api handler...")
}

func checkServerHealth(resp http.ResponseWriter, req *http.Request) {
	// add some health checks here if required
	fmt.Fprintf(resp, "Ok")
}

func getExitSignalsChannel() chan os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	return channel
}

func startWebService() *http.Server {
	ctx, servercancel := context.WithCancel(context.Background())

	initLogger(ctx)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/hello", handleAPICall)
	mux.HandleFunc("/health", checkServerHealth)

	localPort := fmt.Sprintf(":%d", config.LocalPort)

	server = &http.Server{
		Addr:    localPort,
		Handler: mux,
	}

	go func() {
		if config.LocalProtocol == "https" {
			log.Error("%v", server.ListenAndServeTLS(config.CertFile, config.KeyFile))
		} else {
			log.Error("%v", server.ListenAndServe())
		}
		defer ctx.Done()
		defer servercancel()
		defer log.Wait()
	}()
	return server
}

func init() {
	log.Info("Hello World Sample App")
	ExitSignal = getExitSignalsChannel()

	metric.SetupMetric()
}

// initLogger function to configure logger from
// log control file on server StartUp
func initLogger(ctx context.Context) {
	// watch for logctrl file change
	if config.LogControlFile != "" {
		log.EnableLogCtrl(ctx, config.LogControlFile)
	}
	log.SetServiceID(ServiceID)
	log.SetVersion(Version)
	log.EnableInfoLog()
	logMsgBufferSize := 2048
	logOutput := "all" // application can pass the PREL-DR-LOG-007 value [stream, stdout, all]
	if config.LogEndpoint != "" {
		LogEndpoint := config.LogEndpoint
		// check if TLS configuration available
		if config.LogTLSKey != "" && config.LogTLSCert != "" && config.LogTLSCACert != "" &&
			config.logCaCertFilePath != "" && config.rAppLogCertFilePath != "" {
			certFiles := &tlsconf.CertFiles{
				Key:  config.rAppLogCertFilePath + config.LogTLSKey,
				File: config.rAppLogCertFilePath + config.LogTLSCert,
				CA:   config.logCaCertFilePath + config.LogTLSCACert,
			}
			logTLSConfig := tlsconf.NewTLSConfig("", certFiles, true)
			log.Init(ctx, logTLSConfig, LogEndpoint, logMsgBufferSize, logMsgBufferSize, nil, nil, logOutput, "http")
		} else {
			log.Init(ctx, nil, LogEndpoint, logMsgBufferSize, logMsgBufferSize, nil, nil, logOutput, "http")
		}
	}
	log.Info("Logging has been enabled successfully...")
}

func main() {
	go startWebService()
	log.Info("Server is ready to receive web requests")

	<-ExitSignal
	log.Info("Terminating Hello World")
}
