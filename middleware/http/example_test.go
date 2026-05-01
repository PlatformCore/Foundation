package httpmw_test

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/PlatformCore/libpackage/middleware/http"
)

// ExampleNew demonstrates assembling a full enterprise middleware chain.
func ExampleNew() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Build the middleware chain (outermost first).
	chain := middleware.New(
		// 1. Inject unique request ID into context and response header.
		middleware.WithRequestID(),

		// 2. Structured request/response logging.
		middleware.WithLogger(logger, middleware.LoggerOptions{
			SlowRequestThreshold: 500 * time.Millisecond,
			SkipPaths:            []string{"/healthz", "/readyz"},
		}),

		// 3. Recover from panics and return a 500 JSON response.
		middleware.WithRecovery(middleware.RecoveryOptions{
			Logger:     logger,
			PrintStack: true,
		}),

		// 4. Set security headers (HSTS, CSP, X-Frame-Options, etc.).
		middleware.WithSecurity(middleware.SecurityConfig{
			ContentSecurityPolicy: "default-src 'self'",
			HSTSIncludeSubdomains: true,
		}),

		// 5. Handle CORS preflight and set Access-Control-* headers.
		middleware.WithCORS(middleware.CORSConfig{
			AllowedOrigins:   []string{"https://app.example.com"},
			AllowCredentials: true,
		}),

		// 6. Gzip response compression.
		middleware.WithCompress(),

		// 7. OpenTelemetry distributed tracing.
		middleware.WithTracing(middleware.TracingConfig{
			ServiceName: "my-service",
			SkipPaths:   []string{"/healthz", "/readyz"},
		}),

		// 8. Prometheus metrics collection.
		middleware.WithMetrics(middleware.MetricsConfig{
			SkipPaths: []string{"/healthz", "/readyz", "/metrics"},
		}),

		// 9. JWT authentication.
		middleware.WithAuth(middleware.AuthConfig{
			Algorithm:  middleware.AlgorithmHS256,
			HMACSecret: []byte("super-secret-key"),
			Issuer:     "my-auth-service",
			SkipPaths:  []string{"/healthz", "/readyz", "/api/v1/login"},
		}),

		// 10. Token-bucket rate limiting per user.
		middleware.WithRateLimit(middleware.RateLimitConfig{
			RequestsPerSecond: 100,
			Burst:             200,
			Strategy:          middleware.StrategyUser,
		}),

		// 11. Circuit breaker — opens after 5 consecutive 5xx responses.
		middleware.WithCircuitBreaker(middleware.CircuitBreakerConfig{
			Name:        "downstream-api",
			MaxFailures: 5,
			OpenTimeout: 10 * time.Second,
			OnStateChange: func(name string, from, to middleware.CBState) {
				logger.Warn("circuit breaker state change",
					zap.String("breaker", name),
					zap.String("from", from.String()),
					zap.String("to", to.String()),
				)
			},
		}),

		// 12. Cancel request context after 30 seconds.
		middleware.WithTimeout(middleware.TimeoutConfig{
			Timeout: 30 * time.Second,
		}),

		// 13. Cache GET responses for 1 minute.
		middleware.WithCache(middleware.CacheConfig{
			TTL: time.Minute,
		}),
	)

	mux := http.NewServeMux()

	// Register health check endpoints (outside middleware chain).
	middleware.WithHealth(middleware.HealthConfig{
		Checks: map[string]middleware.HealthCheck{
			"database": func(ctx interface{ Done() <-chan struct{} }) error {
				// ping your DB here
				return nil
			},
		},
	}, mux)

	// Apply RBAC after auth.
	adminChain := chain.Append(middleware.RequireRoles("admin", "superuser"))

	mux.Handle("/api/v1/users", chain.ThenFunc(listUsersHandler))
	mux.Handle("/api/v1/admin", adminChain.ThenFunc(adminHandler))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server starting on :8080")
	_ = srv.ListenAndServe()
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	requestID := middleware.GetRequestID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"request_id":%q,"user_id":%q}`, requestID, claims.UserID)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"admin access granted"}`))
}
