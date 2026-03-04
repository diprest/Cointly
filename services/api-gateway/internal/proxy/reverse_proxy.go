package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sony/gobreaker"
)

func NewReverseProxy(target string) http.Handler {
	url, err := url.Parse(target)
	if err != nil {
		slog.Error("Failed to parse proxy target", "target", target, "error", err)
		return http.NotFoundHandler()
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = url.Host
	}

	cbSettings := gobreaker.Settings{
		Name:        "CB-" + target,
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Info("Circuit Breaker State Change", "name", name, "from", from, "to", to)
		},
	}
	cb := gobreaker.NewCircuitBreaker(cbSettings)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := cb.Execute(func() (interface{}, error) {
			proxy.ServeHTTP(w, r)
			return nil, nil
		})

		if err != nil {
			slog.Error("Circuit Breaker blocked request", "target", target, "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
	})
}
