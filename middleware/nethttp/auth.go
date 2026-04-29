package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const authClaimsKey contextKey = "auth_claims"

type AuthAlgorithm string

const (
	AlgorithmHS256 AuthAlgorithm = "HS256"
	AlgorithmHS512 AuthAlgorithm = "HS512"
	AlgorithmRS256 AuthAlgorithm = "RS256"
	AlgorithmRS512 AuthAlgorithm = "RS512"
)

// Claims extends jwt.RegisteredClaims with common enterprise fields.
type Claims struct {
	jwt.RegisteredClaims
	UserID   string   `json:"uid,omitempty"`
	Email    string   `json:"email,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	TenantID string   `json:"tid,omitempty"`
	Scope    string   `json:"scope,omitempty"`
}

// AuthConfig configures JWT authentication.
type AuthConfig struct {
	Algorithm      AuthAlgorithm
	HMACSecret     []byte
	RSAPublicKey   *rsa.PublicKey
	Issuer         string
	Audience       []string
	Leeway         time.Duration
	SkipPaths      []string
	TokenExtractor func(r *http.Request) (string, error)
	OnUnauthorized func(w http.ResponseWriter, r *http.Request, err error)
}

type authErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

// WithAuth returns a JWT authentication middleware.
func WithAuth(cfg AuthConfig) Middleware {
	if cfg.Algorithm == "" {
		cfg.Algorithm = AlgorithmHS256
	}
	if cfg.Leeway == 0 {
		cfg.Leeway = 5 * time.Second
	}
	if cfg.TokenExtractor == nil {
		cfg.TokenExtractor = bearerTokenExtractor
	}
	if cfg.OnUnauthorized == nil {
		cfg.OnUnauthorized = defaultUnauthorized
	}

	skipSet := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipSet[p] = struct{}{}
	}

	parserOpts := []jwt.ParserOption{jwt.WithLeeway(cfg.Leeway)}
	if cfg.Issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(cfg.Issuer))
	}
	if len(cfg.Audience) > 0 {
		parserOpts = append(parserOpts, jwt.WithAudience(cfg.Audience[0]))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, skip := skipSet[r.URL.Path]; skip {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr, err := cfg.TokenExtractor(r)
			if err != nil {
				cfg.OnUnauthorized(w, r, err)
				return
			}
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtKeyFunc(t, cfg)
			}, parserOpts...)
			if err != nil || !token.Valid {
				if err == nil {
					err = errors.New("token is invalid")
				}
				cfg.OnUnauthorized(w, r, err)
				return
			}
			ctx := context.WithValue(r.Context(), authClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func jwtKeyFunc(t *jwt.Token, cfg AuthConfig) (interface{}, error) {
	switch cfg.Algorithm {
	case AlgorithmHS256, AlgorithmHS512:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return cfg.HMACSecret, nil
	case AlgorithmRS256, AlgorithmRS512:
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return cfg.RSAPublicKey, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", cfg.Algorithm)
	}
}

func bearerTokenExtractor(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid Authorization header format; expected 'Bearer <token>'")
	}
	if parts[1] == "" {
		return "", errors.New("empty bearer token")
	}
	return parts[1], nil
}

func defaultUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(authErrorResponse{
		Error:     err.Error(),
		RequestID: GetRequestID(r.Context()),
	})
}

// GetClaims retrieves the JWT claims from the context.
func GetClaims(ctx context.Context) *Claims {
	if c, ok := ctx.Value(authClaimsKey).(*Claims); ok {
		return c
	}
	return nil
}

// RequireRoles returns a middleware that enforces role-based access control.
// Must be placed after WithAuth in the chain.
func RequireRoles(roles ...string) Middleware {
	required := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		required[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			for _, role := range claims.Roles {
				if _, ok := required[role]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":      "forbidden: insufficient roles",
				"request_id": GetRequestID(r.Context()),
			})
		})
	}
}
