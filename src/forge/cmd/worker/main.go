package main

import (
	"forge/internal"
	"forge/internal/service"
	"forge/internal/worker"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	GRPC_SERVER_ADDRESS := os.Getenv("GRPC_SERVER_ADDRESS")
	grpcClient := internal.NewGrpcClient(GRPC_SERVER_ADDRESS)
	defer grpcClient.Close()

	projectService := service.NewProjectServiceClient(grpcClient)

	queueURL := os.Getenv("AWS_SQS_URL")
	workerType := os.Getenv("WORKER_TYPE")

	worker.Run(queueURL, workerType, projectService)
}
