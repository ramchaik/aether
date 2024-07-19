package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		dataString := string(dataBytes)

		log.Printf("Partition Key: %s", kinesisRecord.PartitionKey)
		log.Printf("Sequence Number: %s", kinesisRecord.SequenceNumber)
		log.Printf("Data: %s", dataString)

		// Upload to S3
		// Need to think about the structure

	}

	return nil
}

func main() {
	lambda.Start(handler)
}
