package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// RedisExecutor is the small subset needed by RedisStore, so the ratelimit
// core can work with go-redis or a compatible adapter without adding a compat layer.
type RedisEvalResult interface {
	Result() (interface{}, error)
}

type RedisExecutor interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) RedisEvalResult
}

type RedisStore struct {
	client RedisExecutor
	prefix string
}

func NewRedisStore(client RedisExecutor, prefix string) (*RedisStore, error) {
	if client == nil {
		return nil, ErrRedisExecutorNil
	}
	if prefix == "" {
		prefix = "ratelimit"
	}
	return &RedisStore{client: client, prefix: prefix}, nil
}

func (s *RedisStore) key(key string) string { return s.prefix + ":" + key }

func (s *RedisStore) Increment(ctx context.Context, key string, window time.Duration, now time.Time) (int64, time.Time, error) {
	p := Policy{Limit: 1<<62 - 1, Window: window, Cost: 1, Strategy: StrategyFixedWindow}.Normalize()
	resp, err := s.Eval(ctx, StoreRequest{Key: key, Policy: p, Now: now})
	return resp.Used, resp.ResetAt, err
}

func (s *RedisStore) Eval(ctx context.Context, req StoreRequest) (StoreResponse, error) {
	p := req.Policy.Normalize()
	now := req.Now
	if now.IsZero() {
		now = time.Now()
	}
	script := redisFixedWindowScript
	if p.Strategy == StrategySlidingWindow {
		script = redisSlidingWindowScript
	}
	// Token bucket is usually better handled by RedisLimiter in enterprise_limiter.go;
	// this store still supports fixed/sliding for the direct core Limiter.
	raw, err := s.client.Eval(ctx, script, []string{s.key(req.Key)}, now.UnixMilli(), int64(p.Window/time.Millisecond), p.Limit, p.Cost).Result()
	if err != nil {
		return StoreResponse{}, err
	}
	vals, err := toInt64Slice(raw)
	if err != nil {
		return StoreResponse{}, err
	}
	if len(vals) < 4 {
		return StoreResponse{}, fmt.Errorf("ratelimit redis response: expected 4 fields, got %d", len(vals))
	}
	resetAt := time.UnixMilli(vals[2])
	return StoreResponse{Allowed: vals[0] == 1, Used: vals[1], ResetAt: resetAt, Limit: p.Limit, Remaining: vals[3], Metadata: map[string]string{"mode": string(p.Strategy)}}, nil
}

func toInt64Slice(v interface{}) ([]int64, error) {
	list, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected redis eval result %T", v)
	}
	out := make([]int64, 0, len(list))
	for _, it := range list {
		switch x := it.(type) {
		case int64:
			out = append(out, x)
		case int:
			out = append(out, int64(x))
		case string:
			n, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return nil, err
			}
			out = append(out, n)
		case []byte:
			n, err := strconv.ParseInt(string(x), 10, 64)
			if err != nil {
				return nil, err
			}
			out = append(out, n)
		default:
			return nil, fmt.Errorf("unexpected redis eval value %T", it)
		}
	}
	return out, nil
}

const redisFixedWindowScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window_ms = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])
local reset = redis.call('PTTL', key)
if reset < 0 then
  redis.call('SET', key, cost, 'PX', window_ms)
  local reset_at = now + window_ms
  return {1, cost, reset_at, limit - cost}
end
local used = redis.call('INCRBY', key, cost)
local reset_at = now + reset
local remaining = limit - used
if remaining < 0 then remaining = 0 end
if used <= limit then return {1, used, reset_at, remaining} end
return {0, used, reset_at, remaining}
`

const redisSlidingWindowScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window_ms = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])
redis.call('ZREMRANGEBYSCORE', key, 0, now - window_ms)
local used = redis.call('ZCARD', key)
for i=1,cost do redis.call('ZADD', key, now, now .. ':' .. i .. ':' .. redis.call('INCR', key .. ':seq')) end
redis.call('PEXPIRE', key, window_ms)
redis.call('PEXPIRE', key .. ':seq', window_ms)
used = used + cost
local remaining = limit - used
if remaining < 0 then remaining = 0 end
local reset_at = now + window_ms
if used <= limit then return {1, used, reset_at, remaining} end
return {0, used, reset_at, remaining}
`
