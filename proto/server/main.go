package main

import (
	"Learning/proto/server/pb"
	"Learning/proto/server/service"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

var (
	certFile = "proto/server/server.crt"
	keyFile  = "proto/server/server.key"
)

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc.Creds(creds))           // 创建gRPC服务器
	pb.RegisterGetSchoolServer(s, &service.School{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
