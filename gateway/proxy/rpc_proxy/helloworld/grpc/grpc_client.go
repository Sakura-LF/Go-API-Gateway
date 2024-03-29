package grpc

import (
	"Go-API-Gateway/gateway/proxy/rpc_proxy/helloworld/grpc/Person"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCClient() {
	// 1.连接gRPC服务
	//grpcConn, err := grpc.Dial("127.0.0.1:8004")
	// 抑制安全策略
	grpcConn, err := grpc.Dial("127.0.0.1:8005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("grpc Dial error:", err)
		return
	}
	defer grpcConn.Close()

	// 2.初始化客户端
	grpcClient := pb.NewHelloClient(grpcConn)

	// 3.调用远程服务（函数）
	reply, err := grpcClient.Hello(context.Background(), &pb.Person{Name: "李四", Age: 18})
	if err != nil {
		fmt.Println("reply error:", err)
		return
	}
	fmt.Println("reply:", reply)
}
