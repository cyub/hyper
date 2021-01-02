package middleware

import (
	"net/http"
	"strconv"

	"github.com/cyub/hyper"
	"github.com/cyub/hyper/helper"
	"github.com/cyub/hyper/pkg/limiter"
	redisLimiter "github.com/cyub/hyper/pkg/limiter/redis"
	"github.com/cyub/hyper/redis"
	"github.com/gin-gonic/gin"
)

// RequestSignatureFunc define an function to generate an signature of request
type RequestSignatureFunc func(c *gin.Context) string

// Throttle middleware
func Throttle(limit limiter.Limit, signatureFuncs ...RequestSignatureFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		rdb := redis.Instance()
		if rdb == nil {
			hyper.Logger().Error("please boot with redis as limiter store backend")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "please boot with redis as limiter store backend",
			})
			return
		}
		limiter := redisLimiter.NewRedisLimiter(rdb)
		var key string
		if len(signatureFuncs) > 0 {
			key = signatureFuncs[0](c)
		} else {
			key = requestURLSignature(c)
		}

		rest, err := limiter.Acquire(key, limit)
		if err != nil {
			hyper.Logger().Errorf("limiter acquire fail: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Oops, something went wrong, please try again later",
			})
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(rest.Limit.Burst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(rest.Remaining))
		if !rest.Acquired {
			c.Header("Retry-After", strconv.Itoa(rest.RetryAfter))
			c.Header("X-RateLimit-Reset", strconv.Itoa(rest.ResetAfer))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}

func requestURLSignature(c *gin.Context) string {
	return helper.Md5(c.Request.Method + ":" + c.Request.URL.Path + ":" + c.ClientIP())
}
