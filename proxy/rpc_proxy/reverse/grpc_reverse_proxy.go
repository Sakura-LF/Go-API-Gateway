package reverse

import (
	"Go-API-Gateway/proxy/rpc_proxy/proto"
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
)

func GrpcProxy() {
	var port = flag.Int("port", 8085, "Proxy Server Port")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Println("Listen Failed:", err)
	}

	// 代理服务器不知道/也不必知道,下游服务器的服务名称
	// UnknownServiceHandler
	// 1.处理未知服务名称的处理程序
	// 2.返回一个 ServerOption,允许添加自定义的位置服务处理程序
	// 3.提供的服务方法是:双向流式处理,因为双向流式处理兼容另外三种服务
	// 4.handler 处理所有的客户端调用(拦截器)
	serverOption := grpc.UnknownServiceHandler(handler)
	grpcserver := grpc.NewServer(serverOption)
	grpcserver.Serve(listen)
}

// 统一的入口(handler),用来处理上游所有的请求
// 同时在这个歌handler里面发送请求给下游
func handler(srv any, proxyServerStream grpc.ServerStream) error {
	// 1.过滤非RPC请求
	// MethodFromServerStream 获取到上游想要调用的函数(服务)
	medthodName, ok := grpc.MethodFromServerStream(proxyServerStream)
	if !ok {
		return errors.New("no RPC_Request in this context")
	}
	// 2.构建一个下游连接器: ClientStream
	ctx := proxyServerStream.Context() // 拿到上游客户端请求到代理服务器生成的context
	// 负载均衡算法获取服务器地址
	//grpc.DialContext(ctx)
	// target:下游真实服务器地址
	proxyClientConn, err := grpc.DialContext(ctx, "localhost:8005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// 从上游请求上下文获取元数据
	md, _ := metadata.FromIncomingContext(ctx)
	// 获取取消函数
	outCtx, clientCancel := context.WithCancel(ctx)
	// 封装小有请求的上下文
	outCtx = metadata.NewOutgoingContext(outCtx, md)

	// 1.2封装下游客户单实例 , grpc.NewClientStream
	// 第二个参数为stream的描述信息
	proxyStreamDesc := &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
	proxyClientStream, err := grpc.NewClientStream(outCtx, proxyStreamDesc, proxyClientConn, medthodName)

	// 2.上游与下游数据拷贝
	// 把上游请求信息,发送给下游真实服务器
	s2cErrChan := severToClient(proxyClientStream, proxyServerStream)
	// 把下游响应信息,发回给上游客户端
	c2sErrChan := clientToServer(proxyServerStream, proxyClientStream)

	// 3.关闭双向流
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan: // 向下游发消息
			if s2cErr == io.EOF {
				// 接收到了发送结束额信号,并且不再发送
				proxyClientStream.CloseSend()
			} else {
				// 其他错误
				// 取消发送,并返回错误
				if clientCancel != nil {
					clientCancel()
				}
				return status.Errorf(codes.Internal, "failed proxying server to client:%v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// 返回error,io.EOF; grpc error
			// Trailer：metadata，当流被关闭（ClientStream），读取消息得到error（gRPC，io.EOF）生成元数据
			proxyServerStream.SetTrailer(proxyClientStream.Trailer())

			if c2sErr != io.EOF {
				return c2sErr
			}
		}
	}
	return nil
}

func severToClient(dst grpc.ClientStream, src grpc.ServerStream) chan error {
	res := make(chan error, 1)
	go func() {
		msg := &proto.EchoRequest{}
		for {
			if err := src.RecvMsg(msg); err != nil {
				res <- err
				break
			}
			if err := dst.SendMsg(msg); err != nil {
				res <- err
				break
			}

		}
	}()
	return res
}

func clientToServer(dst grpc.ServerStream, src grpc.ClientStream) chan error {
	res := make(chan error, 1)
	go func() {
		msg := &proto.EchoResponse{}
		for i := 0; ; i++ {
			// response header进行处理
			// 客户端读取响应时，会先读取响应头，然后作出相应的处理
			// 所以有必要设置响应头
			// 细节:只有第一次响应的时候设置响应头
			if i == 0 {
				md, err := src.Header()
				if err != nil {
					res <- err
					break
				}
				if err = dst.SendHeader(md); err != nil {
					res <- err
					break
				}
			}

			if err := src.RecvMsg(msg); err != nil {
				res <- err
				break
			}
			if err := dst.SendMsg(msg); err != nil {
				res <- err
				break
			}

		}
	}()
	return res
}

//func severToClient(proxyServerStream grpc.ServerStream, proxyClientStream grpc.ServerStream) {
//
//}
