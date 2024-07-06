package main

import (
	"forge/internal/worker"
	"os"
)

func main() {
	queueURL := os.Getenv("AWS_SQS_URL")
	workerType := os.Getenv("WORKER_TYPE")

	worker.Run(queueURL, workerType)
}
