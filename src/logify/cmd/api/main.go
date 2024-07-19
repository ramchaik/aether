package main

import (
	"fmt"
	"log"
	"logify/internal/server"
	"net"
	"os"
)

func main() {
	// Start gRPC server
	grpcServer := server.NewGRPCServer()
	GRPC_SERVER_ADDRESS := os.Getenv("GRPC_SERVER_ADDRESS")

	listener, err := net.Listen("tcp", GRPC_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("GRPC server running at", GRPC_SERVER_ADDRESS)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// Start HTTP server
	server := server.NewServer()

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
