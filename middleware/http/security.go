package httpmw

import (
	"net/http"
	"strconv"
	"time"
)

// SecurityConfig configures security hardening headers.
type SecurityConfig struct {
	// ContentSecurityPolicy sets the CSP header value.
	// Example: "default-src 'self'; script-src 'self' 'nonce-{nonce}'"
	ContentSecurityPolicy string
	// HSTSMaxAge sets Strict-Transport-Security max-age in seconds. Default: 1 year.
	HSTSMaxAge time.Duration
	// HSTSIncludeSubdomains enables includeSubDomains in HSTS.
	HSTSIncludeSubdomains bool
	// HSTSPreload enables the HSTS preload flag.
	HSTSPreload bool
	// FrameOptions sets X-Frame-Options. Default: DENY.
	FrameOptions string
	// ContentTypeOptions sets X-Content-Type-Options. Default: nosniff.
	ContentTypeOptions string
	// ReferrerPolicy sets the Referrer-Policy header. Default: strict-origin-when-cross-origin.
	ReferrerPolicy string
	// PermissionsPolicy sets the Permissions-Policy header.
	PermissionsPolicy string
	// DisableHSTS skips the Strict-Transport-Security header (useful in non-TLS dev environments).
	DisableHSTS bool
}

// WithSecurity returns a middleware that injects security-hardening HTTP headers
// on every response, following OWASP recommendations.
func WithSecurity(opts ...SecurityConfig) Middleware {
	cfg := SecurityConfig{
		HSTSMaxAge:            365 * 24 * time.Hour,
		HSTSIncludeSubdomains: true,
		FrameOptions:          "DENY",
		ContentTypeOptions:    "nosniff",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.ContentSecurityPolicy != "" {
			cfg.ContentSecurityPolicy = o.ContentSecurityPolicy
		}
		if o.HSTSMaxAge > 0 {
			cfg.HSTSMaxAge = o.HSTSMaxAge
		}
		if o.FrameOptions != "" {
			cfg.FrameOptions = o.FrameOptions
		}
		if o.ContentTypeOptions != "" {
			cfg.ContentTypeOptions = o.ContentTypeOptions
		}
		if o.ReferrerPolicy != "" {
			cfg.ReferrerPolicy = o.ReferrerPolicy
		}
		if o.PermissionsPolicy != "" {
			cfg.PermissionsPolicy = o.PermissionsPolicy
		}
		cfg.HSTSIncludeSubdomains = o.HSTSIncludeSubdomains
		cfg.HSTSPreload = o.HSTSPreload
		cfg.DisableHSTS = o.DisableHSTS
	}

	// Build static header values.
	hstsValue := buildHSTS(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()

			if !cfg.DisableHSTS && hstsValue != "" {
				h.Set("Strict-Transport-Security", hstsValue)
			}
			if cfg.FrameOptions != "" {
				h.Set("X-Frame-Options", cfg.FrameOptions)
			}
			if cfg.ContentTypeOptions != "" {
				h.Set("X-Content-Type-Options", cfg.ContentTypeOptions)
			}
			if cfg.ReferrerPolicy != "" {
				h.Set("Referrer-Policy", cfg.ReferrerPolicy)
			}
			if cfg.ContentSecurityPolicy != "" {
				h.Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
			}
			if cfg.PermissionsPolicy != "" {
				h.Set("Permissions-Policy", cfg.PermissionsPolicy)
			}
			// Always remove Server header to avoid fingerprinting.
			h.Del("Server")
			// Cross-Origin policies.
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-origin")
			h.Set("X-XSS-Protection", "0") // Modern browsers — CSP is preferred.

			next.ServeHTTP(w, r)
		})
	}
}

func buildHSTS(cfg SecurityConfig) string {
	if cfg.DisableHSTS {
		return ""
	}
	maxAge := int(cfg.HSTSMaxAge.Seconds())
	v := "max-age=" + strconv.Itoa(maxAge)
	if cfg.HSTSIncludeSubdomains {
		v += "; includeSubDomains"
	}
	if cfg.HSTSPreload {
		v += "; preload"
	}
	return v
}
