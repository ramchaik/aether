package server

import (
	"context"
	"log"

	pb "logify/internal/genprotobuf/project_log"

	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedProjectLogServiceServer
}

func (s *grpcServer) PushLogs(ctx context.Context, req *pb.PushLogsRequest) (*pb.PushLogsResponse, error) {
	pId := req.ProjectId
	logs := req.Logs
	log.Println("Adding logs for pId: ", pId, " logs: ", logs)
	return &pb.PushLogsResponse{
		Success: true,
		Message: "Logs pushed successfully",
	}, nil
}

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterProjectLogServiceServer(s, &grpcServer{})
	return s
}
