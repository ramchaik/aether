package service

import (
	"context"
	pb "forge/internal/genprotobuf/project_log"
	"log"
	"time"

	"google.golang.org/grpc"
)

type ProjectLogService interface {
	PushLogs(projectId string, logs []LogEntry) (bool, string)
}

type LogEntry struct {
	Log       string
	Timestamp int64
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

func (p *projectLog) PushLogs(projectId string, logs []LogEntry) (bool, string) {
	c := pb.NewProjectLogServiceClient(p.grpc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pbLogs := make([]*pb.LogEntry, len(logs))
	for i, entry := range logs {
		pbLogs[i] = &pb.LogEntry{
			Log:       entry.Log,
			Timestamp: entry.Timestamp,
		}
	}

	req := &pb.PushLogsRequest{
		ProjectId: projectId,
		Logs:      pbLogs,
	}

	r, err := c.PushLogs(ctx, req)
	if err != nil {
		log.Printf("could not push logs: %v", err)
		return false, "Failed to push logs"
	}

	return r.GetSuccess(), r.GetMessage()
}
