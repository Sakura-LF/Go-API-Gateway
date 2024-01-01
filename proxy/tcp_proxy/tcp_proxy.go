package tcp_proxy

import (
	"context"
	"io"
	"log"
	"net"
	"time"
)

type TCPReverseProxy struct {
	// 下游真实服务器地址host:port
	Addr string

	DialTimeout     time.Duration // 拨号超时时间,持续时间
	Deadline        time.Duration // 拨号截止时间,截止日期
	KeepalivePeriod time.Duration // 长连接超时时间

	// 拨号器,支持自定义:拨号成功,拨号失败,返回error
	DialContext func(ctx context.Context, network, address string) (net.Conn, error)

	// 修改响应
	// 前提:已经向下游服务器发送了请求,并且收到了对应的响应
	// 如果返回错误,则由 ErrorHandler处理
	ModifyResponse func(net.Conn) error

	// 错误处理
	ErrorHandler func(net.Conn, error)
}

func NewSingleHostReverseProxy(addr string) *TCPReverseProxy {
	if addr == "" {
		panic("TCP Server Addr was Empty")
	}

	return &TCPReverseProxy{
		Addr:            addr,            // 下游服务器地址
		DialTimeout:     5 * time.Second, // 拨号超时地址
		Deadline:        time.Minute,     // 拨号截止时间 1ming
		KeepalivePeriod: time.Hour,       // 保活时间: 1h
	}
}

// ServeTCP TCP服务器,实现TCPHandler接口
// 完成上下游连接,及数据的交换
// 1.接收上游连接
// 2.向下游发送请求
// 3.接收下游响应
// 4.拷贝/修改,响应到上游连接
func (tcpProxy *TCPReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	var cancel context.CancelFunc  // 煎炒是否有取消操作
	if tcpProxy.DialTimeout >= 0 { // 连接超时时间时间
		ctx, cancel = context.WithTimeout(ctx, tcpProxy.DialTimeout)
	}
	if tcpProxy.Deadline >= 0 { // 连接截止时间
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(tcpProxy.DialTimeout))
	}
	if cancel != nil {
		defer cancel()
	}
	// 拨号器,使用系统默认的还是自己定义的拨号器
	if tcpProxy.DialContext == nil {
		tcpProxy.DialContext = (&net.Dialer{
			Timeout:   tcpProxy.DialTimeout,              // 连接超时
			Deadline:  time.Now().Add(tcpProxy.Deadline), // 连接截止时间
			KeepAlive: tcpProxy.KeepalivePeriod,          // 长连接超时
		}).DialContext
	}

	// 向下游发送请求
	dst, err := tcpProxy.DialContext(ctx, "tcp", tcpProxy.Addr)
	if err != nil {

		tcpProxy.getErrorHandler()(src, err)
		return
	}
	defer dst.Close() // 关闭下游连接

	// 修改下游服务器响应
	// 如果返回false,说明修改失败,进行错误处理
	if !tcpProxy.modifyResponse(dst) {
		return
	}

	// 修改完信息,进行拷贝
	// 从下游拷贝到上游
	_, err = byteCopy(src, dst)
	if err != nil {
		// 错误处理
		tcpProxy.getErrorHandler()(dst, err)
		return
	}

}

// 通过此函数修改响应,如果没有问题,则返回true,否则返回false
func (tcpProxy *TCPReverseProxy) modifyResponse(res net.Conn) bool {
	// 如果没有定义ModifyResponse,直接返回true,不修改响应
	if tcpProxy.ModifyResponse == nil {
		return true
	}
	// 如果修改想要失败了,进行错误处理
	if err := tcpProxy.ModifyResponse(res); err != nil {
		res.Close() // 关闭连接
		// 错误处理
		tcpProxy.getErrorHandler()(res, err)

	}
	return true
}

func (tcpProxy *TCPReverseProxy) getErrorHandler() func(net.Conn, error) {
	if tcpProxy.ErrorHandler == nil {
		// 执行默认的处理函数
		return tcpProxy.defaultErrorHandler
	}
	return tcpProxy.ErrorHandler
}

func (tcpProxy *TCPReverseProxy) defaultErrorHandler(conn net.Conn, err error) {
	log.Printf("TCP proxy : for conn %v, error: %v\n", conn.RemoteAddr().String(), err)
}

// byteCopy 拷贝两个连接中的数据
// 第一个参数:目标位置
// 第二个参数: 源位置
func byteCopy(dst net.Conn, src net.Conn) (len int64, err error) {
	len, err = io.Copy(dst, src)
	return
}
