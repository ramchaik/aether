package main

import (
	"forge/internal/worker"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	queueURL := os.Getenv("AWS_SQS_URL")
	workerType := os.Getenv("WORKER_TYPE")

	worker.Run(queueURL, workerType)
}
