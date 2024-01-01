package middleware

import (
	"context"
	"net/http"
	"strings"
)

// HandlerFunc 路由处理函数,存储在slice中
type HandlerFunc func(*SliceRouteContext)

// SliceRouter 每个路由器对应多个路由
type SliceRouter struct {
	group []*sliceRoute
}

// 路由:每个路由对应多个处理器
// 维护一个路由数组,同一个请求可以有多个函数处理
type sliceRoute struct {
	// 反向指针,每个路由可以值到它所属的路由组
	*SliceRouter

	// 请求路径
	path string

	// 请求处理列表
	handlers []HandlerFunc
}

// SliceRouteContext 路由上下文
// 每个路由对应一个上下文实例,同时维护请求和响应
type SliceRouteContext struct {
	// 可以通过上下文知道是哪一个路由
	*sliceRoute
	// 当前请求处理执行到哪个位置
	index int8

	Ctx context.Context
	Req *http.Request
	RW  http.ResponseWriter
}

// NewRouter 构造路由器实例
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

// Group 按照指定路径构造路由
func (group *SliceRouter) Group(path string) *sliceRoute {
	return &sliceRoute{
		SliceRouter: group,
		path:        path,
	}
}

// Use 为路由绑定中加件
func (route *sliceRoute) Use(middleware ...HandlerFunc) *sliceRoute {
	// 将中间件追加到handler切片中
	route.handlers = append(route.handlers, middleware...)

	// 当前路由在路由器中是否存在
	var flag bool
	for _, r := range route.SliceRouter.group {
		if route == r {
			flag = true
			break
		}
	}
	if !flag {
		// 不存在,则添加
		route.SliceRouter.group = append(route.SliceRouter.group, route)
	}
	return route
}

// 定义处理器类型函数
// 接收 *SliceRouteContext 类型作为参数
// 返回 http.Handler 结果
type handler func(*SliceRouteContext) http.Handler

// SliceRouterHandler 方法数组路由器的核心处理器
//
//	维护一个方法数组路由器的指针：*SliceRouter
//	支持用户自定义处理器
type SliceRouterHandler struct {
	h handler
	// 维护一个方法海事局路由器的指针,目的将handler和路由绑定
	router *SliceRouter
}

// ServeHTTP 实现了 http.Handler 接口的方法
//
//	作为当前路由器的 http 服务的处理器入口
//
// 步骤:
// 1.初始化路由上下文实例
// 2.检查有无用户自己定义的处理函数
// 3.依次执行路由的处理函数(中间件)
func (rh *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 获得当前请求上下文
	c := newSliceRouterContext(rw, req, rh.router)
	if rh.h != nil {
		c.handlers = append(c.handlers, func(*SliceRouteContext) {
			rh.h(c).ServeHTTP(c.RW, c.Req)
		})
	}
	// 依次执行路由的处理函数(中间件)
	c.Reset()
	c.Next()
}

// 初始化路由上下文实例
func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouteContext {
	// 初始化最长url匹配路由
	sr := &sliceRoute{}
	// 最长url前缀匹配
	matchUrlLen := 0
	for _, route := range r.group {
		// uri匹配成功：前缀匹配
		if strings.HasPrefix(req.RequestURI, route.path) {
			// 记录最长匹配 uri
			pathLen := len(route.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				// 浅拷贝：拷贝数组指针
				*sr = *route
			}
		}
	}

	c := &SliceRouteContext{
		RW:         rw,
		Req:        req,
		Ctx:        req.Context(),
		sliceRoute: sr}
	c.Reset()
	return c
}

// NewSliceRouterHandler 创建 http 服务的处理器
// 将实现了 http.Handler 接口的实例返回
func NewSliceRouterHandler(h handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		h:      h,
		router: router,
	}
}

// Next 从最先加入中间件开始回调
func (c *SliceRouteContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		// 循环调用每一个handler
		c.handlers[c.index](c)
		c.index++
	}
}

//// Abort 跳出中间件方法
//func (c *SliceRouteContext) Abort() {
//	c.index = abortIndex
//}
//
//// IsAborted 是否跳过了回调
//func (c *SliceRouteContext) IsAborted() bool {
//	return c.index >= abortIndex
//}

// Reset 重置回调
func (c *SliceRouteContext) Reset() {
	c.index = -1
}
