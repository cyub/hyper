package redis

import (
	"errors"
	"time"

	"github.com/cyub/hyper/pkg/limiter"
	"github.com/go-redis/redis/v7"
)

type redisLimiter struct {
	client *redis.Client
}

var (
	_ limiter.Limiter = (*redisLimiter)(nil)
	// ErrInvalidLimitResult invalid limit result
	ErrInvalidLimitResult = errors.New("invalid limit result")
)

// NewRedisLimiter return a limiter base redis
func NewRedisLimiter(client *redis.Client) limiter.Limiter {
	return &redisLimiter{client}
}

func (r *redisLimiter) Acquire(key string, limit limiter.Limit) (result limiter.Result, err error) {
	value, err := luaLimiter.Run(r.client, []string{limiter.KeyPrefix + key}, float64(time.Now().UnixNano())/1e9, time.Now().Unix(), limit.Period.Seconds(), limit.Burst).Result()
	if err != nil {
		return
	}

	v, ok := value.([]interface{})
	if !ok || len(v) != 3 {
		err = ErrInvalidLimitResult
		return
	}

	if acquired, _ := v[0].(int64); acquired == 1 {
		result.Acquired = true
	}

	if decay, _ := v[1].(int64); decay > 0 {
		decayTime := time.Unix(decay, 0)
		result.ResetAfer = int(decayTime.Unix())
		result.RetryAfter = int(decayTime.Sub(time.Now()).Seconds())
	}

	if remaining, _ := v[2].(int64); remaining > 0 {
		result.Remaining = int(remaining)
	}
	result.Limit = limit
	return
}

// lockLuaScript
// KEYS[1] - The limiter name
// ARGV[1] - Current time in microseconds
// ARGV[2] - Current time in seconds
// ARGV[3] - Duration of the bucket
// ARGV[4] - Allowed number of tasks
var luaLimiter = redis.NewScript(`
local function reset()
    redis.call('HMSET', KEYS[1], 'start', ARGV[2], 'end', ARGV[2] + ARGV[3], 'count', 1)
    return redis.call('EXPIRE', KEYS[1], ARGV[3] * 2)
end

if redis.call('EXISTS', KEYS[1]) == 0 then
    return {reset(), ARGV[2] + ARGV[3], ARGV[4] - 1}
end

if ARGV[1] >= redis.call('HGET', KEYS[1], 'start') and ARGV[1] <= redis.call('HGET', KEYS[1], 'end') then
	local acquire = 0
	if tonumber(redis.call('HINCRBY', KEYS[1], 'count', 1)) <= tonumber(ARGV[4]) then
		acquire = 1
	end
	return {
		acquire,
		tonumber(redis.call('HGET', KEYS[1], 'end')),
		ARGV[4] - redis.call('HGET', KEYS[1], 'count')
	}
end

return {reset(), ARGV[2] + ARGV[3], ARGV[4] - 1}
`)
