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
	count := testutil.CollectAndCount(counter, "client_requests_total")
	if count != 1 {
		t.Errorf("expected 1, got %d\n", count)
	}
}
