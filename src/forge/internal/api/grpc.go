package api

import (
	"context"
	"log"
	"time"

	pb "forge/internal/proto"

	"google.golang.org/grpc"
)

func NewGrpcClient(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	return conn, nil
}

func TestSaveProjectURL(conn *grpc.ClientConn, url string, projectId string) {
	c := pb.NewProjectServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SaveProjectUrl(ctx, &pb.SaveProjectUrlRequest{ProjectUrl: url, ProjectId: projectId})
	if err != nil {
		log.Fatalf("could not save project: %v", err)
	}

	log.Printf("Response: %s", r.GetMessage())
}
