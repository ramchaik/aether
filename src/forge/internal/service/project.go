package service

import (
	"context"
	pb "forge/internal/genprotobuf"
	"log"
	"time"

	"google.golang.org/grpc"
)

type ProjectService interface {
	UpdateProjectStatus(projectId string, status pb.ProjectStatus)
}

type project struct {
	grpc *grpc.ClientConn
}

func NewProjectServiceClient(grpcConn *grpc.ClientConn) *project {
	return &project{
		grpc: grpcConn,
	}
}

func (p *project) UpdateProjectStatus(projectId string, status pb.ProjectStatus) {
	c := pb.NewProjectServiceClient(p.grpc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UpdateProjectStatus(ctx, &pb.UpdateProjectStatusRequest{
		ProjectId: projectId,
		Status:    status,
	})
	if err != nil {
		log.Fatalf("could not update project status: %v", err)
	}

	log.Printf("Response: %s", r.GetMessage())
}
