package utils

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func getSession() *session.Session {
	// Retrieve AWS credentials from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")

	// Create a credentials object using the retrieved credentials
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken)

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_REGION")),
	}))

	return sess
}

func GetSQSService() *sqs.SQS {
	sess := getSession()
	sqsSvc := sqs.New(sess)
	return sqsSvc
}
