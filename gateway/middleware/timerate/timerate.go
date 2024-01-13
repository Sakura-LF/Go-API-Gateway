package timerate

import (
	"Go-API-Gateway/gateway/middleware/router"
	"fmt"
	"golang.org/x/time/rate"
)

func RateLimiter(params ...int) func(c *router.SliceRouteContext) {
	// 默认参数
	var r rate.Limit = 1 // 没秒生成的令牌数量
	var b = 2            // 令牌总数
	// 获取传入的参数 r,b
	if len(params) == 2 {
		r = rate.Limit(params[0])
		b = params[1]
	}
	// 创建限流器
	limiter := rate.NewLimiter(r, b)

	return func(c *router.SliceRouteContext) {
		// 1.如果无法获取到token,则跳出中间件,直接返回
		if !limiter.Allow() {
			c.RW.Write([]byte(fmt.Sprintf("rate limit:")))
			c.Abort()
			return
		}
		// 2.可以获取到token,执行中间件
		c.Next()
	}
}
