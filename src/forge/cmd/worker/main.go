package main

import (
	"forge/internal/api"
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
	conn, err := api.NewGrpcClient(GRPC_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	api.TestSaveProjectURL(conn, "https://abc.com", "123-proj-id")

	queueURL := os.Getenv("AWS_SQS_URL")
	workerType := os.Getenv("WORKER_TYPE")

	worker.Run(queueURL, workerType)
}
