package monitoring

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP Metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	HttpActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests being processed",
		},
		[]string{"method", "endpoint"},
	)

	// Custom Business Metrics
	UserRegistrations = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	SearchQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_queries_total",
			Help: "Total number of search queries",
		},
		[]string{"query_type"},
	)
)

// RegisterMetrics registers various Prometheus metrics used for monitoring
// the application. It logs the start and successful completion of the
// registration process. The metrics registered include:
//  - HttpRequestsTotal: Total number of HTTP requests
//  - HttpRequestDuration: Duration of HTTP requests
//  - HttpActiveRequests: Number of active HTTP requests
//  - HttpRequestSize: Size of HTTP requests
//  - HttpResponseSize: Size of HTTP responses
//  - UserRegistrations: Number of user registrations
//  - SearchQueries: Number of search queries

func RegisterMetrics() {
	LogInfo("Registering Prometheus metrics", nil)
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpActiveRequests)
	prometheus.MustRegister(UserRegistrations)
	prometheus.MustRegister(SearchQueries)
	LogInfo("Prometheus metrics registered successfully", nil)
}

// ExposeMetrics sets up an HTTP handler for Prometheus metrics at the endpoint "/api/metrics"
// and starts a server on port 9090 to expose these metrics. If the server fails to start,
// an error is logged. A log message is also generated when the server starts successfully.

func ExposeMetrics(router *mux.Router) {
	router.Handle("/api/metrics", promhttp.Handler())
	LogInfo("Starting Prometheus metrics server on port 9090", nil)
	go func() {
		if err := http.ListenAndServe(":9090", router); err != nil {
			fmt.Printf("Error starting Prometheus metrics server: %s", err)
		}
	}()
}

// IncrementHTTPRequest increments the HTTP request counter
func IncrementHTTPRequest(method, endpoint string) {
	HttpRequestsTotal.WithLabelValues(method, endpoint).Inc()
}

// IncrementActiveRequests increments the active HTTP requests counter
func IncrementActiveRequests(method, endpoint string) {
	HttpActiveRequests.WithLabelValues(method, endpoint).Inc()
}

// DecrementActiveRequests decrements the active HTTP requests counter
func DecrementActiveRequests(method, endpoint string) {
	HttpActiveRequests.WithLabelValues(method, endpoint).Dec()
}

// IncrementUserRegistrations increments the user registrations counter
func IncrementUserRegistrations() {
	UserRegistrations.Inc()
}

// IncrementSearchQueries increments the search queries counter
func IncrementSearchQueries(queryType string) {
	SearchQueries.WithLabelValues(queryType).Inc()
}
