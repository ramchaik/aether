package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func GetAWSConfig() aws.Config {
	AWS_REGION := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return cfg
}

func GetKinesisClient() *kinesis.Client {
	cfg := GetAWSConfig()
	client := kinesis.NewFromConfig(cfg)
	return client
}

func PushDataToKinesisStream(data map[string]any) {
	client := GetKinesisClient()

	streamName := os.Getenv("AWS_KINESIS_STREAM")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
	}

	partitionKey := os.Getenv("AWS_KINESIS_STREAM_PARTITION_KEY")
	input := &kinesis.PutRecordInput{
		Data:         jsonData,
		PartitionKey: aws.String(partitionKey),
		StreamName:   aws.String(streamName),
	}

	result, err := client.PutRecord(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to put record to Kinesis: %v", err)
	}

	fmt.Printf("Successfully put logs record to Kinesis. Shard ID: %s, Sequence number: %s\n",
		*result.ShardId, *result.SequenceNumber)
}
