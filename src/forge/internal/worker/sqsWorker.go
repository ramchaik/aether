package worker

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ProcessMessage takes a message and performs the necessary actions based on the message content.
func ProcessMessage(message *sqs.Message, workerType string) bool {
	messageBody := *message.Body
	messageType := *message.MessageAttributes["MessageType"].StringValue

	fmt.Println("Body:", messageBody)
	fmt.Printf("Attributes: %v", messageType)

	if messageType != workerType {
		return false
	}

	// Process the message here
	fmt.Print("Processing message\n")

	// Return true if the message should be deleted, false otherwise
	return true
}

// Run listens to an SQS queue and processes messages.
func Run(queueURL string, workerType string) {
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

	sqsSvc := sqs.New(sess)

	fmt.Printf("Listening to SQS: %v\n\n", queueURL)

	for {
		result, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:              &queueURL,
			MaxNumberOfMessages:   aws.Int64(10),
			VisibilityTimeout:     aws.Int64(20),
			WaitTimeSeconds:       aws.Int64(0),
			MessageAttributeNames: []*string{aws.String("All")},
		})
		if err != nil {
			log.Fatalf("ReceiveMessage failed %v", err)
		}

		for _, message := range result.Messages {
			if !ProcessMessage(message, workerType) {
				break
			}

			// Delete the message from the queue after processing
			_, err = sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("DeleteMessage Failed %v", err)
			}
		}
	}
}
