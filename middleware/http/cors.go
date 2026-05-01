package httpmw

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CORSConfig configures CORS behavior.
type CORSConfig struct {
	// AllowedOrigins is a list of origins allowed. Use ["*"] for any origin.
	AllowedOrigins []string
	// AllowedMethods lists allowed HTTP methods. Defaults to common safe methods.
	AllowedMethods []string
	// AllowedHeaders lists headers the client is allowed to send.
	AllowedHeaders []string
	// ExposedHeaders lists headers the browser is allowed to read.
	ExposedHeaders []string
	// AllowCredentials controls the Access-Control-Allow-Credentials header.
	AllowCredentials bool
	// MaxAge sets the preflight cache duration. Defaults to 12h.
	MaxAge time.Duration
	// OptionsPassthrough passes OPTIONS requests to the next handler after setting headers.
	OptionsPassthrough bool
}

type corsHandler struct {
	cfg            CORSConfig
	allowedOrigins map[string]struct{}
	allowAll       bool
	maxAgeStr      string
	allowedMethods string
	allowedHeaders string
	exposedHeaders string
}

// WithCORS returns a CORS middleware that handles preflight requests and sets
// appropriate Access-Control-* headers on all responses.
func WithCORS(cfg CORSConfig) Middleware {
	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = []string{
			http.MethodGet, http.MethodHead, http.MethodPost,
			http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
		}
	}
	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = []string{
			"Accept", "Authorization", "Content-Type",
			"X-Request-ID", "X-Correlation-ID",
		}
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 12 * time.Hour
	}

	h := &corsHandler{
		cfg:            cfg,
		allowedOrigins: make(map[string]struct{}),
		maxAgeStr:      strconv.Itoa(int(cfg.MaxAge.Seconds())),
		allowedMethods: strings.Join(cfg.AllowedMethods, ", "),
		allowedHeaders: strings.Join(cfg.AllowedHeaders, ", "),
		exposedHeaders: strings.Join(cfg.ExposedHeaders, ", "),
	}
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			h.allowAll = true
			break
		}
		h.allowedOrigins[o] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !h.isAllowedOrigin(origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Set standard CORS headers.
			h.setResponseHeaders(w, origin)

			// Handle preflight.
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Max-Age", h.maxAgeStr)
				w.WriteHeader(http.StatusNoContent)
				if !cfg.OptionsPassthrough {
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *corsHandler) isAllowedOrigin(origin string) bool {
	if h.allowAll {
		return true
	}
	_, ok := h.allowedOrigins[origin]
	return ok
}

func (h *corsHandler) setResponseHeaders(w http.ResponseWriter, origin string) {
	header := w.Header()
	if h.allowAll && !h.cfg.AllowCredentials {
		header.Set("Access-Control-Allow-Origin", "*")
	} else {
		header.Set("Access-Control-Allow-Origin", origin)
		header.Add("Vary", "Origin")
	}
	header.Set("Access-Control-Allow-Methods", h.allowedMethods)
	header.Set("Access-Control-Allow-Headers", h.allowedHeaders)
	if h.exposedHeaders != "" {
		header.Set("Access-Control-Expose-Headers", h.exposedHeaders)
	}
	if h.cfg.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
}
