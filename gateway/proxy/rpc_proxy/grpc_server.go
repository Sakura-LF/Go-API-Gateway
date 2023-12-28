package rpc_proxy

import (
	"Go-API-Gateway/gateway/proxy/rpc_proxy/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 8005, "Grpc server port")

func GrpcServer() {
	flag.Parse()
	// 监听端口
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Println("监听失败,", err)
	}
	server := grpc.NewServer()

	proto.RegisterEchoServer(server, &Servers{})

	server.Serve(listen)
}

type Servers struct {
	proto.UnimplementedEchoServer
}

// UnaryEcho 一元RPC服务的实现
// 元数据,map[string][string]
// token,timestamp,授权
func (s *Servers) UnaryEcho(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	fmt.Println("----------UnaryEcho-----------")
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Println("metadata:", incomingContext)
	return &proto.EchoResponse{Message: req.Message}, nil
}

// ServerStreamingEcho 服务端流式处理RPC
func (s *Servers) ServerStreamingEcho(req *proto.EchoRequest, stream proto.Echo_ServerStreamingEchoServer) error {
	fmt.Println("----------ServerStreamingEcho-----------")

	// 服务端向客户单响应使用send函数
	for i := 0; i < 5; i++ { // 以流的形式发送多次
		err := stream.Send(&proto.EchoResponse{Message: req.Message})
		if err != nil {
			return err
		}
	}

	return nil
}

// ClientStreamingEcho 客户端流式处理
func (s *Servers) ClientStreamingEcho(stream proto.Echo_ClientStreamingEchoServer) error {
	fmt.Println("--------------ClientStreamingEcho------------------------")
	// 接收客户端消息
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("receive finish")
			// 发送并关闭
			return stream.SendAndClose(&proto.EchoResponse{Message: "Receive Finish"})
		} else if err != nil {
			log.Println(err)
			return err
		}
		// 打印接收到的消息
		fmt.Println("request reveived: ", recv.Message)
	}
}

// 双向流式处理
func (s *Servers) BidirectionalStreamingEcho(stream proto.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Println("--------------BidirectionalStreamingEcho------------------------")
	// 每接收到客户端的消息,就向客户端发送响应
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("消息发送完毕")
			stream.Send(&proto.EchoResponse{Message: "Receive over"})
			return nil
		} else if err != nil {
			log.Println("Receive Failed", err)
			return err
		}
		// 打印接收到消息
		log.Println("Client Message:", recv)
		// 回复响应
		err = stream.Send(&proto.EchoResponse{Message: "Receive success!"})
		if err != nil {
			log.Println("Send Fail,", err)
		}
	}
}
