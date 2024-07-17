package service

import (
	"context"
	pb "forge/internal/proto"
	"log"
	"time"

	"google.golang.org/grpc"
)

type ProjectService interface {
	SaveProjectURL(url string, projectId string)
}

type project struct {
	grpc *grpc.ClientConn
}

func NewProjectServiceClient(grpcConn *grpc.ClientConn) *project {
	return &project{
		grpc: grpcConn,
	}
}

func (p *project) SaveProjectURL(url string, projectId string) {
	c := pb.NewProjectServiceClient(p.grpc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SaveProjectUrl(ctx, &pb.SaveProjectUrlRequest{ProjectUrl: url, ProjectId: projectId})
	if err != nil {
		log.Fatalf("could not save project: %v", err)
	}

	log.Printf("Response: %s", r.GetMessage())
}
