package grpc

import (
	"Go-API-Gateway/proxy/rpc_proxy/helloworld/grpc/Person"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func GRPCSever() {
	// 1.注册服务,启动监听
	grpcserver := grpc.NewServer()
	pb.RegisterHelloServer(grpcserver, &HelloService{})

	// 2.监听
	listen, err := net.Listen("tcp", "127.0.0.1:8005")
	if err != nil {
		log.Println("监听失败,", err)
	}
	defer listen.Close()

	// 3.启动服务
	grpcserver.Serve(listen)
}

type HelloService struct {
	pb.pb
}

func (HelloService) Hello(context.Context, *pb.Person) (*pb.Person, error) {
	person := &pb.Person{
		Name: "Sakura",
		Age:  22,
	}
	return person, nil
}
