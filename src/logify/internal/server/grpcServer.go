package server

import (
	"context"
	"log"

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
		"log":       req.LogEntry.Log,
		"timestamp": req.LogEntry.Timestamp,
	}

	err := utils.PushDataToKinesisStream(data)
	if err != nil {
		log.Printf("Failed to push log to Kinesis stream: %v", err)
		return &pb.PushLogsResponse{
			Success: false,
			Message: "Failed to push log to Kinesis stream",
		}, err
	}

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
