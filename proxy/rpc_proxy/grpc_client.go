package rpc_proxy

import (
	proto2 "Go-API-Gateway/proxy/rpc_proxy/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
)

var msg = "Client Data"

func GrpcClient() {
	conn, err := grpc.Dial("127.0.0.1:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := proto2.NewEchoClient(conn)

	// 1.调用一元RPC方法
	unaryEchoWithMetadata(client, msg)

	// 2.调用服务端流式处理
	serverStreaming(client, msg)

	// 3.调用客户端流式处理
	clientStreaming(client, msg)

	// 4.双向流式处理
	serverclientStreaming(client, msg)
}

// 调用一元RPC方法
func unaryEchoWithMetadata(client proto2.EchoClient, msg string) {
	fmt.Println("---------UnaryEcho Client---------")

	// Pairs封装一个metadata
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	//md.Append("authorization", "token....")
	ctx := metadata.NewOutgoingContext(context.Background(), md) // 即将输出的请求

	resp, err := client.UnaryEcho(ctx, &proto2.EchoRequest{Message: msg})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Grpc Server Response:", resp)
	}
}

// 调用服务端流式处理
func serverStreaming(client proto2.EchoClient, msg string) {
	fmt.Println("---------serverStreaming Client---------")

	// Pairs封装一个metadata
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	//md.Append("authorization", "token....")
	ctx := metadata.NewOutgoingContext(context.Background(), md) // 即将输出的请求

	stream, err := client.ServerStreamingEcho(ctx, &proto2.EchoRequest{Message: msg})
	if err != nil {
		fmt.Println(err)
	}
	// 从流中接收数据,每次读一个消息
	for {
		// err 读取到末尾会返回 EOF
		recv, err := stream.Recv()
		if err != nil {
			log.Println(err)
			break
		} else if err == io.EOF {
			fmt.Println("finish receive")
			return
		}
		fmt.Println("response is :,", recv.Message)
	}
}

// 客户端流式处理
func clientStreaming(client proto2.EchoClient, msg string) {
	fmt.Println("---------clientStreaming---------")

	// Pairs封装一个metadata
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	//md.Append("authorization", "token....")
	ctx := metadata.NewOutgoingContext(context.Background(), md) // 即将输出的请求

	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// 向服务端循环发送消息
	for i := 0; i < 5; i++ {
		err = stream.Send(&proto2.EchoRequest{Message: msg})
		if err != nil {
			log.Println("Fialed to send", err)
		}
	}
	// 关闭并接受
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
	}
	log.Println("Servere response:", recv)
}

// 双向流式处理
func serverclientStreaming(client proto2.EchoClient, msg string) {
	fmt.Println("---------serverclientStreaming---------")

	// Pairs封装一个metadata
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	//md.Append("authorization", "token....")
	ctx := metadata.NewOutgoingContext(context.Background(), md) // 即将输出的请求

	stream, err := client.BidirectionalStreamingEcho(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// 开启一个协程向服务端发送消息
	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&proto2.EchoRequest{Message: msg})
			if err != nil {
				log.Println(err)
			}
		}
		stream.CloseSend()
	}()
	//time.Sleep(time.Second * 10)

	//循环接收服务器的消息
	for {
		recv, err1 := stream.Recv()
		if err1 == io.EOF {
			log.Println("----- Receive Finish ------ ")
			break
		} else if err != nil {
			log.Println("Receive Fail,", err)
		}
		log.Println("Server Response:", recv)
	}
}
