package service

import (
	"Learning/proto/server/pb"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

type School struct {
	pb.UnimplementedGetSchoolServer
}

func (s *School) GetSchool(ctx context.Context, req *pb.GetSchoolRequest) (*pb.GetSchoolResponse, error) {
	return &pb.GetSchoolResponse{
		Data: &pb.School{
			Name: "School of Computer Science and Technology",
			Classes: []*pb.Class{
				{
					Grade: pb.Grade_FirstGrade,
					Class: 5,
					Students: []*pb.Student{
						{Name: "Alice", Male: true, Scores: []int32{80, 85, 90}},
						{Name: "Bob", Male: false, Scores: []int32{60, 70, 65}},
					},
				},
			},
		},
		Code:    1,
		Message: "success",
	}, nil
}

func (s *School) GetStudents(in *emptypb.Empty, stream pb.GetSchool_GetStudentsServer) error {
	students := []*pb.Student{
		{Name: "Alice", Male: true, Scores: []int32{80, 85, 90}},
		{Name: "Bob", Male: false, Scores: []int32{60, 70, 65}},
	}
	for _, student := range students {
		if err := stream.Send(student); err != nil {
			return err
		}
	}
	return nil
}

func (s *School) SendStudents(stream pb.GetSchool_SendStudentsServer) error {
	for {
		// 接收客户端发来的流式数据
		res, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			return stream.SendAndClose(&pb.WebResponse{
				Code:    1,
				Message: "finish",
			})
		}
		if err != nil {
			return err
		}
		fmt.Println(res)
	}
}
