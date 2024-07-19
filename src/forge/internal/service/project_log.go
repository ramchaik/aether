package service

import (
	"context"
	pb "forge/internal/genprotobuf/project_log"
	"log"
	"time"

	"google.golang.org/grpc"
)

type ProjectLogService interface {
	PushLogs(projectId string, logs string) (bool, string)
}

type projectLog struct {
	grpc *grpc.ClientConn
	pb.UnimplementedProjectLogServiceServer
}

func NewProjectLogServiceClient(grpcConn *grpc.ClientConn) *projectLog {
	return &projectLog{
		grpc: grpcConn,
	}
}

func (p *projectLog) PushLogs(projectId string, buildLog string) (bool, string) {
	c := pb.NewProjectLogServiceClient(p.grpc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.PushLogsRequest{
		ProjectId: projectId,
		Log:       buildLog,
		Timestamp: time.Now().Unix(), // Current time as Unix timestamp (seconds since epoch)
	}

	r, err := c.PushLogs(ctx, req)
	if err != nil {
		log.Fatalf("could not push logs: %v", err)
		return false, "Failed to push logs"
	}

	log.Printf("Response: %s", r.GetMessage())
	return r.Success, r.Message
}
