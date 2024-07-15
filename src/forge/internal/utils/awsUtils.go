package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// getConfig retrieves AWS credentials and region from environment variables and returns a session.
func getConfig() (*aws.Config, error) {
	// Retrieve AWS credentials from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")

	// Create a credentials object using the retrieved credentials
	creds := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		secretAccessKey,
		sessionToken,
	)

	region := os.Getenv("AWS_REGION")

	// Load the Shared AWS Config
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	return &cfg, nil
}

// GetSQSService create and returns SQS client
func GetSQSService() (*sqs.Client, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	// Create SQS service client
	sqsClient := sqs.NewFromConfig(*cfg)

	return sqsClient, nil
}

// GetS3Service creates and returns an S3 client.
func GetS3Service() (*s3.Client, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	// Create S3 service client
	s3Client := s3.NewFromConfig(*cfg)

	return s3Client, nil
}
