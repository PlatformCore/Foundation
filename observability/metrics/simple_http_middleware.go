package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type HTTPConfig struct {
	Prefix    string
	SkipPaths []string
}

func SimpleHTTPMiddleware(reg *SimpleRegistry, cfg HTTPConfig) func(http.Handler) http.Handler {
	if reg == nil {
		reg = NewSimpleRegistry()
	}
	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "http.server"
	}
	skip := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skip[p] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := skip[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()
			rw := &metricsResponseWriter{ResponseWriter: w, status: http.StatusOK}
			reg.Inc(prefix + ".in_flight")
			defer reg.Add(prefix+".in_flight", -1)
			next.ServeHTTP(rw, r)
			status := strconv.Itoa(rw.status)
			reg.Inc(prefix + ".requests_total")
			reg.Inc(prefix + ".status." + status)
			reg.Observe(prefix+".duration_ms", float64(time.Since(start).Milliseconds()))
			reg.Observe(prefix+".response_size_bytes", float64(rw.size))
		})
	}
}

type metricsResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *metricsResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *metricsResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}
