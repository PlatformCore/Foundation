package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"go.uber.org/zap"
)

// RecoveryOptions configures the recovery middleware.
type RecoveryOptions struct {
	// Logger is used to log panic details. If nil, uses fmt.Println.
	Logger *zap.Logger
	// StackSize limits the stack trace buffer size in bytes. Default: 64KB.
	StackSize int
	// PrintStack controls whether the full stack trace is logged.
	PrintStack bool
	// OnPanic is an optional hook called after a panic is recovered.
	// It receives the recovered value and the stack trace.
	OnPanic func(recovered interface{}, stack []byte, r *http.Request)
	// ErrorResponse is an optional custom JSON response body on panic.
	// If nil, a default JSON error is written.
	ErrorResponse func(recovered interface{}) interface{}
}

// panicResponse is the default JSON error body.
type panicResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

// WithRecovery returns a middleware that recovers from panics, logs them,
// and writes a 500 Internal Server Error JSON response.
func WithRecovery(opts ...RecoveryOptions) Middleware {
	cfg := RecoveryOptions{
		StackSize:  64 << 10, // 64 KB
		PrintStack: true,
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.Logger != nil {
			cfg.Logger = o.Logger
		}
		if o.StackSize > 0 {
			cfg.StackSize = o.StackSize
		}
		cfg.PrintStack = o.PrintStack
		cfg.OnPanic = o.OnPanic
		cfg.ErrorResponse = o.ErrorResponse
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Capture stack trace.
					buf := make([]byte, cfg.StackSize)
					n := runtime.Stack(buf, false)
					stack := buf[:n]

					if cfg.PrintStack {
						if cfg.Logger != nil {
							cfg.Logger.Error("panic recovered",
								zap.Any("panic", rec),
								zap.String("request_id", GetRequestID(r.Context())),
								zap.String("method", r.Method),
								zap.String("path", r.URL.Path),
								zap.ByteString("stack", stack),
							)
						} else {
							fmt.Printf("[RECOVERY] panic: %v\n%s\n", rec, debug.Stack())
						}
					}

					if cfg.OnPanic != nil {
						cfg.OnPanic(rec, stack, r)
					}

					// Write error response only if headers haven't been sent.
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.Header().Set("X-Content-Type-Options", "nosniff")
					w.WriteHeader(http.StatusInternalServerError)

					var body interface{}
					if cfg.ErrorResponse != nil {
						body = cfg.ErrorResponse(rec)
					} else {
						body = panicResponse{
							Error:     "internal server error",
							RequestID: GetRequestID(r.Context()),
						}
					}
					_ = json.NewEncoder(w).Encode(body)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
