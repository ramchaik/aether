package server

import (
	"context"

	pb "logify/internal/genprotobuf/project_log"
	"logify/internal/utils"

	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedProjectLogServiceServer
}

func (s *grpcServer) PushLogs(ctx context.Context, req *pb.PushLogsRequest) (*pb.PushLogsResponse, error) {

	data := map[string]any{
		"projectId": req.ProjectId,
		"log":       req.Log,
		"timestamp": req.Timestamp,
	}

	utils.PushDataToKinesisStream(data)

	return &pb.PushLogsResponse{
		Success: true,
		Message: "Log pushed successfully",
	}, nil
}

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterProjectLogServiceServer(s, &grpcServer{})
	return s
}
