package client

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewInstrumentedClient returns a new instrumented http Client.
func NewInstrumentedClient(r prometheus.Registerer) *http.Client {
	return InstrumentClient(&http.Client{}, r)
}

var counter *prometheus.CounterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "client_requests_total",
		Help: "A counter for requests from the wrapped client.",
	},
	[]string{"code", "method"},
)

var inFlightGauge prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "client_in_flight_requests",
	Help: "A gauge of in-flight requests for the wrapped client.",
})

// dnsLatencyVec uses custom buckets based on expected dns durations.
// It has an instance label "event", which is set in the
// DNSStart and DNSDonehook functions defined in the
// InstrumentTrace struct below.
var dnsLatencyVec *prometheus.HistogramVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "dns_duration_seconds",
		Help:    "Trace dns latency histogram.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"event"},
)

// tlsLatencyVec uses custom buckets based on expected tls durations.
// It has an instance label "event", which is set in the
// TLSHandshakeStart and TLSHandshakeDone hook functions defined in the
// InstrumentTrace struct below.
var tlsLatencyVec *prometheus.HistogramVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "tls_duration_seconds",
		Help:    "Trace tls latency histogram.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"event"},
)

// histVec has no labels, making it a zero-dimensional ObserverVec.
var histVec *prometheus.HistogramVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "client_request_duration_seconds",
		Help:    "A histogram of request latencies.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{},
)

// InstrumentClient instruments the given http Client.
func InstrumentClient(c *http.Client, r prometheus.Registerer) *http.Client {
	// Define functions for the available httptrace.ClientTrace hook
	// functions that we want to instrument.
	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_start").Observe(t)
		},
		DNSDone: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_done").Observe(t)
		},
		TLSHandshakeStart: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_start").Observe(t)
		},
		TLSHandshakeDone: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_done").Observe(t)
		},
	}

	r.MustRegister(counter, tlsLatencyVec, dnsLatencyVec, histVec, inFlightGauge)

	// Wrap the default RoundTripper with middleware.
	c.Transport = promhttp.InstrumentRoundTripperInFlight(inFlightGauge,
		promhttp.InstrumentRoundTripperCounter(counter,
			promhttp.InstrumentRoundTripperTrace(trace,
				promhttp.InstrumentRoundTripperDuration(histVec, http.DefaultTransport),
			),
		),
	)

	return c
}
