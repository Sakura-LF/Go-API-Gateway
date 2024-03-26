package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

// RateLimiterMiddleware 限流中间件
func RateLimiterMiddleware(rateLimiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), "limiter", rateLimiter)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
