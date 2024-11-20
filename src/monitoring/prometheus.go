package monitoring

import (
	"net/http"

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

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
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

	HttpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6), // 100B, 1KB, 10KB...
		},
		[]string{"method", "endpoint"},
	)

	HttpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6),
		},
		[]string{"method", "endpoint"},
	)

	// Database Metrics
	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	DBFailedQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_failed_queries_total",
			Help: "Total number of failed database queries",
		},
		[]string{"query_type"},
	)

	// Cache Metrics
	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_name"},
	)

	CacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_name"},
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
//  - DBQueryDuration: Duration of database queries
//  - DBFailedQueries: Number of failed database queries
//  - CacheHits: Number of cache hits
//  - CacheMisses: Number of cache misses
//  - UserRegistrations: Number of user registrations
//  - SearchQueries: Number of search queries
func RegisterMetrics() {
	LogInfo("Registering Prometheus metrics", nil)
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpActiveRequests)
	prometheus.MustRegister(HttpRequestSize)
	prometheus.MustRegister(HttpResponseSize)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(DBFailedQueries)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(UserRegistrations)
	prometheus.MustRegister(SearchQueries)
	LogInfo("Prometheus metrics registered successfully", nil)
}

// ExposeMetrics sets up an HTTP handler for Prometheus metrics at the endpoint "/api/metrics"
// and starts a server on port 9090 to expose these metrics. If the server fails to start,
// an error is logged. A log message is also generated when the server starts successfully.
func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			LogError(err, "Failed to start Prometheus metrics server", nil)
		}
	}()
	LogInfo("Prometheus metrics server started on port 9090", nil)
}

// IncrementHTTPRequest increments the HTTP request counter
func IncrementHTTPRequest(method, endpoint string) {
	HttpRequestsTotal.WithLabelValues(method, endpoint).Inc()
}

// ObserveHTTPRequestDuration observes the duration of an HTTP request
func ObserveHTTPRequestDuration(method, endpoint string, duration float64) {
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// IncrementActiveRequests increments the active HTTP requests counter
func IncrementActiveRequests(method, endpoint string) {
	HttpActiveRequests.WithLabelValues(method, endpoint).Inc()
}

// DecrementActiveRequests decrements the active HTTP requests counter
func DecrementActiveRequests(method, endpoint string) {
	HttpActiveRequests.WithLabelValues(method, endpoint).Dec()
}

// ObserveRequestSize observes the size of an HTTP request
func ObserveRequestSize(method, endpoint string, size float64) {
	HttpRequestSize.WithLabelValues(method, endpoint).Observe(size)
}

// ObserveResponseSize observes the size of an HTTP response
func ObserveResponseSize(method, endpoint string, size float64) {
	HttpResponseSize.WithLabelValues(method, endpoint).Observe(size)
}

// IncrementDBQuery increments the database query counter
func IncrementDBQuery(queryType string) {
	DBFailedQueries.WithLabelValues(queryType).Inc()
}

// ObserveDBQueryDuration observes the duration of a database query
func ObserveDBQueryDuration(queryType string, duration float64) {
	DBQueryDuration.WithLabelValues(queryType).Observe(duration)
}

// IncrementCacheHit increments the cache hit counter
func IncrementCacheHit(cacheName string) {
	CacheHits.WithLabelValues(cacheName).Inc()
}

// IncrementCacheMiss increments the cache miss counter
func IncrementCacheMiss(cacheName string) {
	CacheMisses.WithLabelValues(cacheName).Inc()
}

// IncrementUserRegistrations increments the user registrations counter
func IncrementUserRegistrations() {
	UserRegistrations.Inc()
}

// IncrementSearchQueries increments the search queries counter
func IncrementSearchQueries(queryType string) {
	SearchQueries.WithLabelValues(queryType).Inc()
}