package httpmw

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// CacheBackend is the interface for pluggable cache storage.
type CacheBackend interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// CacheConfig configures response caching.
type CacheConfig struct {
	// Backend is the storage backend. If nil, an in-process LRU-like map is used.
	Backend CacheBackend
	// TTL is the default cache duration. Default: 1 minute.
	TTL time.Duration
	// MaxBodySize is the maximum response body size to cache in bytes. Default: 1MB.
	MaxBodySize int
	// Methods lists the HTTP methods to cache. Default: [GET, HEAD].
	Methods []string
	// KeyFunc builds the cache key from the request.
	// If nil, SHA-256 of method+path+sorted query string is used.
	KeyFunc func(r *http.Request) string
	// ShouldCache determines whether the response should be cached.
	// If nil, only 200-299 responses are cached.
	ShouldCache func(status int, header http.Header) bool
	// SkipPaths lists paths exempt from caching.
	SkipPaths []string
}

// cachedResponse stores a serialized HTTP response.
type cachedResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

// inMemoryCache is a simple thread-safe in-process cache.
type inMemoryCache struct {
	mu    sync.RWMutex
	items map[string]*inMemoryItem
}

type inMemoryItem struct {
	value   []byte
	expires time.Time
}

func newInMemoryCache() *inMemoryCache {
	c := &inMemoryCache{items: make(map[string]*inMemoryItem)}
	go c.evict()
	return c
}

func (c *inMemoryCache) Get(_ context.Context, key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok || time.Now().After(item.expires) {
		return nil, false
	}
	return item.value, true
}

func (c *inMemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &inMemoryItem{value: value, expires: time.Now().Add(ttl)}
	return nil
}

func (c *inMemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *inMemoryCache) evict() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		c.mu.Lock()
		for k, v := range c.items {
			if now.After(v.expires) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

// WithCache returns a response caching middleware.
// Cache hits return the stored response without invoking downstream handlers.
func WithCache(cfg CacheConfig) Middleware {
	if cfg.Backend == nil {
		cfg.Backend = newInMemoryCache()
	}
	if cfg.TTL == 0 {
		cfg.TTL = time.Minute
	}
	if cfg.MaxBodySize == 0 {
		cfg.MaxBodySize = 1 << 20 // 1 MB
	}
	if len(cfg.Methods) == 0 {
		cfg.Methods = []string{http.MethodGet, http.MethodHead}
	}
	if cfg.ShouldCache == nil {
		cfg.ShouldCache = func(status int, _ http.Header) bool {
			return status >= 200 && status < 300
		}
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = defaultCacheKey
	}

	methodSet := make(map[string]struct{}, len(cfg.Methods))
	for _, m := range cfg.Methods {
		methodSet[m] = struct{}{}
	}
	skipSet := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipSet[p] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, skip := skipSet[r.URL.Path]; skip {
				next.ServeHTTP(w, r)
				return
			}
			if _, ok := methodSet[r.Method]; !ok {
				next.ServeHTTP(w, r)
				return
			}
			// Skip if client sends Cache-Control: no-cache.
			if r.Header.Get("Cache-Control") == "no-cache" {
				next.ServeHTTP(w, r)
				return
			}

			key := cfg.KeyFunc(r)

			// Cache hit?
			if raw, ok := cfg.Backend.Get(r.Context(), key); ok {
				var cached cachedResponse
				if err := json.Unmarshal(raw, &cached); err == nil {
					for k, vals := range cached.Headers {
						for _, v := range vals {
							w.Header().Set(k, v)
						}
					}
					w.Header().Set("X-Cache", "HIT")
					w.WriteHeader(cached.Status)
					_, _ = w.Write(cached.Body)
					return
				}
			}

			// Cache miss — record the response.
			buf := &bytes.Buffer{}
			crw := &cachingWriter{ResponseWriter: w, buf: buf, maxSize: cfg.MaxBodySize}
			next.ServeHTTP(crw, r)

			if cfg.ShouldCache(crw.status, w.Header()) && !crw.overflow {
				entry := cachedResponse{
					Status:  crw.status,
					Headers: w.Header().Clone(),
					Body:    buf.Bytes(),
				}
				raw, _ := json.Marshal(entry)
				_ = cfg.Backend.Set(r.Context(), key, raw, cfg.TTL)
			}
			w.Header().Set("X-Cache", "MISS")
		})
	}
}

// cachingWriter captures the response body for caching.
type cachingWriter struct {
	http.ResponseWriter
	buf      *bytes.Buffer
	maxSize  int
	status   int
	overflow bool
}

func (cw *cachingWriter) WriteHeader(code int) {
	cw.status = code
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *cachingWriter) Write(b []byte) (int, error) {
	if cw.status == 0 {
		cw.status = http.StatusOK
	}
	if !cw.overflow {
		if cw.buf.Len()+len(b) > cw.maxSize {
			cw.overflow = true
		} else {
			cw.buf.Write(b)
		}
	}
	return cw.ResponseWriter.Write(b)
}

func defaultCacheKey(r *http.Request) string {
	q := r.URL.Query()
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+strings.Join(q[k], ","))
	}
	raw := r.Method + ":" + r.URL.Path + "?" + strings.Join(parts, "&")
	sum := sha256.Sum256([]byte(raw))
	return "cache:" + hex.EncodeToString(sum[:])
}
