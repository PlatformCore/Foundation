package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Validator is the interface that request body validators must implement.
type Validator interface {
	Validate() error
}

// ValidationError represents a structured validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Message)
}

// ValidationErrors is a slice of ValidationError that implements error.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	msgs := make([]string, len(ve))
	for i, e := range ve {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// ValidatorConfig configures the request validation middleware.
type ValidatorConfig struct {
	// MaxBodyBytes is the maximum allowed request body size. Default: 1MB.
	MaxBodyBytes int64
	// OnValidationError is called when validation fails.
	OnValidationError func(w http.ResponseWriter, r *http.Request, errs ValidationErrors)
	// OnDecodeError is called when JSON decoding fails.
	OnDecodeError func(w http.ResponseWriter, r *http.Request, err error)
}

type validatedBodyKeyType struct{ typ string }

// WithBodyValidator returns a middleware that decodes a JSON request body
// into type T, calls Validate() on it, and stores the result in the context.
//
// Example:
//
//	chain.Append(middleware.WithBodyValidator[CreateUserRequest](cfg))
//	// In handler: body, _ := middleware.GetValidatedBody[CreateUserRequest](r.Context())
func WithBodyValidator[T Validator](cfg ValidatorConfig) Middleware {
	if cfg.MaxBodyBytes == 0 {
		cfg.MaxBodyBytes = 1 << 20
	}
	if cfg.OnValidationError == nil {
		cfg.OnValidationError = defaultValidationError
	}
	if cfg.OnDecodeError == nil {
		cfg.OnDecodeError = defaultDecodeError
	}

	var zero T
	bodyKey := validatedBodyKeyType{typ: fmt.Sprintf("%T", zero)}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil || r.Body == http.NoBody {
				next.ServeHTTP(w, r)
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxBodyBytes)
			defer r.Body.Close()

			body, err := io.ReadAll(r.Body)
			if err != nil {
				cfg.OnDecodeError(w, r, err)
				return
			}

			var target T
			if err := json.Unmarshal(body, &target); err != nil {
				cfg.OnDecodeError(w, r, err)
				return
			}

			if err := target.Validate(); err != nil {
				var ve ValidationErrors
				switch e := err.(type) {
				case ValidationErrors:
					ve = e
				case ValidationError:
					ve = ValidationErrors{e}
				default:
					ve = ValidationErrors{{Field: "body", Message: err.Error()}}
				}
				cfg.OnValidationError(w, r, ve)
				return
			}

			ctx := context.WithValue(r.Context(), bodyKey, target)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetValidatedBody retrieves the validated request body from the context.
// T must match the type used with WithBodyValidator.
func GetValidatedBody[T any](ctx context.Context) (T, bool) {
	var zero T
	key := validatedBodyKeyType{typ: fmt.Sprintf("%T", zero)}
	v, ok := ctx.Value(key).(T)
	return v, ok
}

// --- Validation Helpers ---

// RequiredField returns a ValidationError if the value is blank.
func RequiredField(field, value string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: "is required"}
	}
	return nil
}

// MinLength returns a ValidationError if value is shorter than min.
func MinLength(field, value string, min int) *ValidationError {
	if len(value) < min {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must be at least %d characters", min)}
	}
	return nil
}

// MaxLength returns a ValidationError if value is longer than max.
func MaxLength(field, value string, max int) *ValidationError {
	if len(value) > max {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must be at most %d characters", max)}
	}
	return nil
}

// IsEmail performs basic email format validation.
func IsEmail(field, value string) *ValidationError {
	if !strings.Contains(value, "@") || strings.LastIndex(value, ".") <= strings.Index(value, "@") {
		return &ValidationError{Field: field, Message: "must be a valid email address"}
	}
	return nil
}

// IsOneOf validates that a value is in an allowed set.
func IsOneOf[T comparable](field string, value T, allowed ...T) *ValidationError {
	for _, a := range allowed {
		if reflect.DeepEqual(value, a) {
			return nil
		}
	}
	return &ValidationError{Field: field, Message: "is not an allowed value"}
}

func defaultValidationError(w http.ResponseWriter, r *http.Request, errs ValidationErrors) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":      "validation failed",
		"details":    errs,
		"request_id": GetRequestID(r.Context()),
	})
}

func defaultDecodeError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":      "invalid request body: " + err.Error(),
		"request_id": GetRequestID(r.Context()),
	})
}
