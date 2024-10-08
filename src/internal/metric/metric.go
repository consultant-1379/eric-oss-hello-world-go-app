// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	servicePrefix = "hello_world"
)

var (
	// Registry prometheus registry TODO: should it be exported
	Registry = prometheus.NewRegistry()
	// RequestsTotal total number of API requests
	RequestsTotal prometheus.Counter
	// RequestsFailedTotal total number of API request failures
	RequestsFailedTotal prometheus.Counter
	// HelloWorldHTTPRequestsTotal total number of HTTP responses by status codes
	HelloWorldHTTPRequestsTotal *prometheus.CounterVec
)

func createMetrics() {
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "requests_total",
			Help:      "Total number of API requests",
		})
	RequestsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "requests_failed_total",
			Help:      "Total number of API requests failures",
		})
	HelloWorldHTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "hello_world_http_requests_total",
			Help:      "Total number of HTTP responses by status codes",
		},
		[]string{"code"})
}

func registerMetrics() {
	Registry.Register(RequestsTotal)
	Registry.Register(RequestsFailedTotal)
	Registry.Register(HelloWorldHTTPRequestsTotal)
}

// SetupMetric setups the metric
func SetupMetric() {
	createMetrics()
	registerMetrics()
}
