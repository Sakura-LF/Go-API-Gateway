package tcp_proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// TCP服务器，实现服务与代理分离
// 创建一个TCP服务器：
// 1.监听端口
// 2.获取连接
// 3.封装新连接对象，设置服务参数上下文、超时、连接关闭
// 4.回调handLer（需要定义一个接口）

// TCPServer TCP 服务和兴结构体,监听指定主机,并提供服务
// Addr 必须,主机地址
// Handler 必选,回调函数,处理TCP请求,提供默认实现
type TCPServer struct {
	Addr    string     // 主机地址
	Handler TCPHandler // 回调函数,处理TCP请求

	BaseContext context.Context // 上下文,收集取消,终止,错误等信息
	err         error           // TCP Error

	ReadTimeout      time.Duration // 读超时
	WriteTimeout     time.Duration // 写超时
	KeepAliveTimeout time.Duration // 长连接超时

	mu         sync.Mutex         // 连接关闭等关键工作需要实现
	doneChan   chan struct{}      // 当前的服务完成会向channel写入信号
	inShutdown int32              // 服务终止: 0-未关闭,1-已关闭
	l          *onceCloseListener //服务监听器,使用完成后进行关
	//onShutdown []func()      // 关掉服务的时候,会有一个回调函数,会去执行对应的逻辑

}

type TCPHandler interface {
	// ServeTCP 提供TCP服务
	// ctx:连接上下文
	// conn: TCP连接实例,用于读写操作
	ServeTCP(ctx context.Context, conn net.Conn)
}

// 定义一个默认实现TCPHandler接口的结构体
type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	conn.Write([]byte("kong! TCP Server Handler"))
}

// ListenAndServe 模仿HTTP创建ListenAndServer 方法
// 一个无参数,绑定结构体的  func ListenAndServe(addr string, handler Handler) error
// 一个有参数直接调用的 func (tcpserver *TCPServer) ListenAndServe() error

func ListenAndServe(addr string, handler TCPHandler) error {
	server := &TCPServer{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

// 定义一个错误服务器关闭错误
var (
	ErrServerClosed     = errors.New("http: Server closed")
	ErrAbortHandler     = errors.New("net/http: abort Handler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

func (tcpserver *TCPServer) ListenAndServe() error {
	// 如果返回true说明当前服务器已经关闭
	if tcpserver.shuttingDown() {
		return ErrServerClosed
	}
	addr := tcpserver.Addr
	if addr == "" {
		return errors.New("addr was empty")
	}
	if tcpserver.Handler == nil {
		tcpserver.Handler = &tcpHandler{} // 如果为空,执行默认handler
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return tcpserver.Serve(ln)
}

func (tcpserver *TCPServer) Serve(l net.Listener) error {
	tcpserver.l = &onceCloseListener{Listener: l}
	defer l.Close() // 执行监听器的关闭

	if tcpserver.BaseContext == nil {
		tcpserver.BaseContext = context.Background()
	}

	baseCtx := tcpserver.BaseContext
	ctx := context.WithValue(baseCtx, ServerContextKey, tcpserver)
	for {
		rw, err := l.Accept()
		if err != nil {
			if tcpserver.shuttingDown() {
				return ErrServerClosed
			}
			return err
		}
		c := tcpserver.newConn(rw) // 对TCPConn的二次封装
		go c.serve(ctx)            // 有连接了,启动一个协程提供服务
		// 提供的服服务就是ServeTCP,如果重写了就执行重写的ServeTCP,如果没有,就执行默认的
	}
}

// Create new connection from rwc
// 封装一个连接
func (tcpserver *TCPServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server:     tcpserver,
		rwc:        rwc,
		remoteAddr: rwc.RemoteAddr().String(),
	}
	// 设置参数,从TCPServer 中取字段,赋值给TCPConn
	if t := tcpserver.ReadTimeout; t != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(t))
	}
	if t := tcpserver.WriteTimeout; t != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(t))
	}
	if t := tcpserver.KeepAliveTimeout; t != 0 {
		// 将Conn断言为TCPConn
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(t)
		}
	}
	return c
}

// 封装一个连接去提供服务
// 核心方法serve
func (c *conn) serve(ctx context.Context) {
	// 复制context中的远程地址
	if ra := c.rwc.RemoteAddr(); ra != nil {
		c.remoteAddr = ra.String()
	}
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("tcp: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		c.rwc.Close()
	}()

	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	if c.server.Handler == nil { // 如果handler为nil
		panic("TCP Handler Empty")
	}
	// 对于http来说是request和response,对于tcp来说,只需要一个conn
	c.server.Handler.ServeTCP(ctx, c.rwc) // 上下文和连接,想要读取的数据从哪个里面获取方便

}

type conn struct {
	// 当前服务器的实例
	server *TCPServer
	// 连接
	rwc net.Conn
	// 远程地址(当前请求的地址)
	remoteAddr string
}

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }

func (tcpserver *TCPServer) Close() {
	// 将inShutdown置为1
	atomic.StoreInt32(&tcpserver.inShutdown, 1) // 使用原子操作
	close(tcpserver.doneChan)                   // 关闭channel
	tcpserver.l.Close()                         // 关闭监听
	//return nil
}

func (tcpserver *TCPServer) shuttingDown() bool {
	// 0:正在关闭 1-已关闭
	return atomic.LoadInt32(&tcpserver.inShutdown) != 0
}

type contextKey struct {
	name string
}
