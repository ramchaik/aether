package api

import (
	"context"
	"log"
	"time"

	pb "forge/internal/proto"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn *grpc.ClientConn
}

func NewGrpcClient(addr string) (*GrpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	return &GrpcClient{
		conn: conn,
	}, nil
}

func (gc *GrpcClient) TestSaveProjectURL(url string, projectId string) {
	c := pb.NewProjectServiceClient(gc.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SaveProjectUrl(ctx, &pb.SaveProjectUrlRequest{ProjectUrl: url, ProjectId: projectId})
	if err != nil {
		log.Fatalf("could not save project: %v", err)
	}

	log.Printf("Response: %s", r.GetMessage())
}
