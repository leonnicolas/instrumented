package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCounter(t *testing.T) {
	r := prometheus.NewRegistry()

	c := NewInstrumentedClient(r)

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	server := httptest.NewServer(m)
	defer server.Close()

	_, err := c.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	t.Run("counter", func(t *testing.T) {
		count := testutil.CollectAndCount(counter, "client_requests_total")
		if count != 1 {
			t.Errorf("expected 1, got %d\n", count)
		}
	})
}
func TestLint(t *testing.T) {
	for _, c := range []prometheus.Collector{counter, inFlightGauge, dnsLatencyVec, tlsLatencyVec, histVec} {
		t.Run("lint", func(t *testing.T) {
			p, err := testutil.CollectAndLint(c)
			if err != nil {
				t.Errorf("linting error: %v\n", err)
			}
			if p != nil {
				t.Errorf("linting problems: %v\n", p)
			}

		})
	}
}
