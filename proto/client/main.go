package main

import (
	"Learning/proto/client/pb"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"time"
)

var (
	addr     = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	certFile = "./proto/client/server.crt"
)

func main() {
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGetSchoolClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetSchool(ctx, &pb.GetSchoolRequest{Name: ""})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("res: %v", r)
	fmt.Println("================================================================")
	GetStudents(client)
	fmt.Println("================================================================")
	SendStudent(client)
}

func GetStudents(c pb.GetSchoolClient) {
	// server端流式RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.GetStudents(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("c.GetStudents failed, err: %v", err)
	}
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.GetStudents failed, err: %v", err)
		}
		log.Printf("got res: %v\n", res)
	}
}
func SendStudent(c pb.GetSchoolClient) {
	students := []*pb.Student{
		{Name: "Alice", Male: true, Scores: []int32{80, 85, 90}},
		{Name: "Bob", Male: false, Scores: []int32{60, 70, 65}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端流式RPC
	stream, err := c.SendStudents(ctx)
	if err != nil {
		log.Fatalf("c.SendStudents failed, err: %v", err)
	}
	for _, student := range students {
		// 发送流式数据
		err = stream.Send(student)
		if err != nil {
			log.Fatalf("c.SendStudents stream.Send(%v) failed, err: %v", student, err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.SendStudents failed: %v", err)
	}
	log.Printf("got res: %v", res)
}
