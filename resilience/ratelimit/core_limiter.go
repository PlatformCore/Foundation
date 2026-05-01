package ratelimit

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrLimited          = errors.New("rate limited")
	ErrNoStore          = errors.New("store is nil")
	ErrEmptyIdentity    = errors.New("empty key identity")
	ErrRedisExecutorNil = errors.New("redis executor is nil")
)

type Strategy string

const (
	StrategyFixedWindow   Strategy = "fixed_window"
	StrategySlidingWindow Strategy = "sliding_window"
	StrategyTokenBucket   Strategy = "token_bucket"
	StrategyLeakyBucket   Strategy = "leaky_bucket"
)

const (
	HeaderLimit     = "X-RateLimit-Limit"
	HeaderRemaining = "X-RateLimit-Remaining"
	HeaderReset     = "X-RateLimit-Reset"
)

type Key struct {
	Namespace  string
	Identity   string
	Route      string
	Method     string
	Tenant     string
	Dimensions map[string]string
}

func NewKey(namespace, identity string) Key {
	return Key{Namespace: namespace, Identity: identity, Dimensions: map[string]string{}}
}

func (k Key) Validate() error {
	if strings.TrimSpace(k.Identity) == "" {
		return ErrEmptyIdentity
	}
	return nil
}

func (k Key) String() string {
	identity := normalizePart(k.Identity)
	namespace := normalizePart(k.Namespace)
	if namespace != "" {
		identity = namespace + ":" + identity
	}
	parts := make([]string, 0, 4+len(k.Dimensions))
	if r := normalizePart(k.Route); r != "" {
		parts = append(parts, "route="+r)
	}
	if m := strings.ToUpper(strings.TrimSpace(k.Method)); m != "" {
		parts = append(parts, "method="+m)
	}
	if t := normalizePart(k.Tenant); t != "" {
		parts = append(parts, "tenant="+t)
	}
	if len(k.Dimensions) > 0 {
		keys := make([]string, 0, len(k.Dimensions))
		for n := range k.Dimensions {
			keys = append(keys, n)
		}
		sort.Strings(keys)
		for _, n := range keys {
			parts = append(parts, normalizePart(n)+"="+normalizePart(k.Dimensions[n]))
		}
	}
	if len(parts) == 0 {
		return identity
	}
	return identity + "|" + strings.Join(parts, "|")
}

func normalizePart(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	return strings.ReplaceAll(v, " ", "_")
}

type Policy struct {
	Name                string
	Limit               int64
	Window              time.Duration
	Strategy            Strategy
	Cost                int64
	Burst               int64
	RefillRatePerSecond float64
	ShadowMode          bool
}

func (p Policy) Normalize() Policy {
	if p.Name == "" {
		p.Name = "default"
	}
	if p.Limit <= 0 {
		p.Limit = 100
	}
	if p.Window <= 0 {
		p.Window = time.Minute
	}
	if p.Strategy == "" {
		p.Strategy = StrategyFixedWindow
	}
	if p.Cost <= 0 {
		p.Cost = 1
	}
	if p.Burst <= 0 {
		p.Burst = p.Limit
	}
	if p.RefillRatePerSecond <= 0 && p.Window > 0 {
		p.RefillRatePerSecond = float64(p.Limit) / p.Window.Seconds()
	}
	return p
}

type Result struct {
	Allowed    bool
	Limit      int64
	Remaining  int64
	Used       int64
	ResetAt    time.Time
	RetryAfter time.Duration
	PolicyName string
	Strategy   Strategy
	ShadowMode bool
}

type Decision struct {
	Key       string
	Result    Result
	Err       error
	Timestamp time.Time
}

type DecisionHook func(Decision)

type Store interface {
	Increment(ctx context.Context, key string, window time.Duration, now time.Time) (count int64, resetAt time.Time, err error)
}

type StoreRequest struct {
	Key    string
	Policy Policy
	Now    time.Time
}
type StoreResponse struct {
	Used      int64
	Limit     int64
	Remaining int64
	ResetAt   time.Time
	Allowed   bool
	Metadata  map[string]string
}
type EvalStore interface {
	Eval(ctx context.Context, req StoreRequest) (StoreResponse, error)
}

type Options struct {
	Policy   Policy
	Store    Store
	Now      func() time.Time
	FailOpen bool
	OnResult DecisionHook
}

func (o Options) normalize() Options {
	o.Policy = o.Policy.Normalize()
	if o.Now == nil {
		o.Now = time.Now
	}
	return o
}

type Limiter struct {
	policy   Policy
	store    Store
	now      func() time.Time
	failOpen bool
	hook     DecisionHook
}

func New(opts Options) *Limiter {
	opts = opts.normalize()
	return &Limiter{
		policy:   opts.Policy,
		store:    opts.Store,
		now:      opts.Now,
		failOpen: opts.FailOpen,
		hook:     opts.OnResult,
	}
}

func (l *Limiter) Allow(ctx context.Context, key Key) (Result, error) {
	now := l.now()
	res, err := l.allow(ctx, key, now)
	if l.hook != nil {
		l.hook(Decision{
			Key:       key.String(),
			Result:    res,
			Err:       err,
			Timestamp: now,
		})
	}
	return res, err
}

func (l *Limiter) allow(ctx context.Context, key Key, now time.Time) (Result, error) {
	if err := key.Validate(); err != nil {
		return Result{}, err
	}
	if l.store == nil {
		if l.failOpen {
			return Result{Allowed: true, PolicyName: l.policy.Name}, nil
		}
		return Result{}, ErrNoStore
	}
	p := l.policy.Normalize()
	if eval, ok := l.store.(EvalStore); ok {
		resp, err := eval.Eval(ctx, StoreRequest{Key: key.String(), Policy: p, Now: now})
		if err != nil {
			if l.failOpen {
				return Result{Allowed: true, Limit: p.Limit, Remaining: p.Limit, ResetAt: now.Add(p.Window), PolicyName: p.Name, Strategy: p.Strategy}, nil
			}
			return Result{}, err
		}
		res := resultFromStoreResponse(p, now, resp, resp.Allowed || p.ShadowMode)
		if !res.Allowed && !p.ShadowMode {
			return res, ErrLimited
		}
		return res, nil
	}
	count := int64(0)
	var resetAt time.Time
	for i := int64(0); i < p.Cost; i++ {
		c, r, err := l.store.Increment(ctx, key.String(), p.Window, now)
		if err != nil {
			if l.failOpen {
				return Result{Allowed: true, Limit: p.Limit, Remaining: p.Limit, ResetAt: now.Add(p.Window), PolicyName: p.Name, Strategy: p.Strategy}, nil
			}
			return Result{}, err
		}
		count = c
		resetAt = r
	}
	remaining := max64(0, p.Limit-count)
	res := Result{Allowed: count <= p.Limit || p.ShadowMode, Limit: p.Limit, Remaining: remaining, Used: count, ResetAt: resetAt, PolicyName: p.Name, Strategy: p.Strategy, ShadowMode: p.ShadowMode}
	if !res.Allowed && !p.ShadowMode {
		res.RetryAfter = maxDuration(0, resetAt.Sub(now))
		return res, ErrLimited
	}
	return res, nil
}

func resultFromStoreResponse(p Policy, now time.Time, resp StoreResponse, allowed bool) Result {
	res := Result{Allowed: allowed, Limit: resp.Limit, Remaining: resp.Remaining, Used: resp.Used, ResetAt: resp.ResetAt, PolicyName: p.Name, Strategy: p.Strategy, ShadowMode: p.ShadowMode}
	if res.Limit == 0 {
		res.Limit = p.Limit
	}
	if res.ResetAt.IsZero() {
		res.ResetAt = now.Add(p.Window)
	}
	if !res.Allowed {
		res.RetryAfter = maxDuration(0, res.ResetAt.Sub(now))
	}
	return res
}

func (r Result) Headers() http.Header {
	h := make(http.Header)
	h.Set(HeaderLimit, strconv.FormatInt(r.Limit, 10))
	h.Set(HeaderRemaining, strconv.FormatInt(r.Remaining, 10))
	if !r.ResetAt.IsZero() {
		h.Set(HeaderReset, strconv.FormatInt(r.ResetAt.Unix(), 10))
	}
	if r.RetryAfter > 0 {
		h.Set("Retry-After", strconv.Itoa(int(r.RetryAfter.Seconds())))
	}
	return h
}

func (l *Limiter) Peek(ctx context.Context, key Key) (Result, error) {
	return l.Allow(ctx, key)
}

func (l *Limiter) MustAllow(ctx context.Context, key Key) Result {
	res, err := l.Allow(ctx, key)
	if err != nil && !errors.Is(err, ErrLimited) {
		panic(err)
	}
	return res
}

func (l *Limiter) UpdatePolicy(p Policy) {
	l.policy = p.Normalize()
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
