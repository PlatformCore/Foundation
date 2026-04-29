package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// CompressConfig configures response compression.
type CompressConfig struct {
	// Level is the gzip compression level. Default: gzip.DefaultCompression.
	Level int
	// MinSize is the minimum response body size in bytes to compress. Default: 1024.
	MinSize int
	// ContentTypes is the list of Content-Type values to compress.
	// If empty, a default set of compressible types is used.
	ContentTypes []string
}

var defaultCompressibleTypes = []string{
	"text/html", "text/css", "text/plain", "text/xml", "text/javascript",
	"application/json", "application/javascript", "application/xml",
	"application/rss+xml", "application/atom+xml", "image/svg+xml",
}

// gzipPool reuses gzip.Writer instances.
var gzipPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
		return w
	},
}

// WithCompress returns a middleware that gzip-compresses responses for clients
// that send Accept-Encoding: gzip. Only compressible content types and responses
// above the MinSize threshold are compressed.
func WithCompress(opts ...CompressConfig) Middleware {
	cfg := CompressConfig{
		Level:        gzip.DefaultCompression,
		MinSize:      1024,
		ContentTypes: defaultCompressibleTypes,
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.Level != 0 {
			cfg.Level = o.Level
		}
		if o.MinSize > 0 {
			cfg.MinSize = o.MinSize
		}
		if len(o.ContentTypes) > 0 {
			cfg.ContentTypes = o.ContentTypes
		}
	}

	ctSet := make(map[string]struct{}, len(cfg.ContentTypes))
	for _, ct := range cfg.ContentTypes {
		ctSet[ct] = struct{}{}
	}

	pool := sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(io.Discard, cfg.Level)
			return w
		},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !acceptsGzip(r) {
				next.ServeHTTP(w, r)
				return
			}

			gw := &gzipWriter{
				ResponseWriter: w,
				pool:           &pool,
				minSize:        cfg.MinSize,
				ctSet:          ctSet,
			}
			defer gw.close()

			w.Header().Del("Content-Length") // Will be recomputed.
			next.ServeHTTP(gw, r)
		})
	}
}

type gzipWriter struct {
	http.ResponseWriter
	pool    *sync.Pool
	gz      *gzip.Writer
	minSize int
	ctSet   map[string]struct{}
	buf     []byte
	status  int
	ready   bool
}

func (gw *gzipWriter) WriteHeader(code int) {
	gw.status = code
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	if !gw.ready {
		gw.buf = append(gw.buf, b...)
		if len(gw.buf) >= gw.minSize && gw.shouldCompress() {
			gw.startGzip()
			return gw.gz.Write(gw.buf)
		}
		return len(b), nil
	}
	if gw.gz != nil {
		return gw.gz.Write(b)
	}
	return gw.ResponseWriter.Write(b)
}

func (gw *gzipWriter) shouldCompress() bool {
	ct := gw.ResponseWriter.Header().Get("Content-Type")
	for k := range gw.ctSet {
		if strings.HasPrefix(ct, k) {
			return true
		}
	}
	return false
}

func (gw *gzipWriter) startGzip() {
	gw.ready = true
	gz := gw.pool.Get().(*gzip.Writer)
	gz.Reset(gw.ResponseWriter)
	gw.gz = gz
	gw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	gw.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
	gw.ResponseWriter.Header().Del("Content-Length")
	if gw.status != 0 {
		gw.ResponseWriter.WriteHeader(gw.status)
	}
}

func (gw *gzipWriter) close() {
	if !gw.ready {
		// Buffer didn't reach MinSize or wrong content type — write plain.
		if gw.status != 0 {
			gw.ResponseWriter.WriteHeader(gw.status)
		}
		if len(gw.buf) > 0 {
			_, _ = gw.ResponseWriter.Write(gw.buf)
		}
		return
	}
	if gw.gz != nil {
		_ = gw.gz.Close()
		gw.pool.Put(gw.gz)
	}
}

func acceptsGzip(r *http.Request) bool {
	ae := r.Header.Get("Accept-Encoding")
	return strings.Contains(ae, "gzip")
}
