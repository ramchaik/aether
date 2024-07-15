package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"forge/internal/utils"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Message struct {
	RepoURL      string `json:"repoURL"`
	BuildCommand string `json:"buildCommand"`
}

// ProcessMessage takes a message and performs the necessary actions based on the message content.
func ProcessMessage(message *sqs.Message, workerType string) bool {
	if message == nil {
		log.Println("Received nil message")
		return false
	}

	if message.Body == nil {
		log.Println("Message body is nil")
		return false
	}

	if message.MessageAttributes == nil {
		log.Println("Message attributes is nil")
		return false
	}

	messageBody := *message.Body
	messageType := *message.MessageAttributes["MessageType"].StringValue

	fmt.Println("Body:", messageBody)
	fmt.Printf("Attributes: %v", messageType)

	if messageType != workerType {
		return false
	}

	var msg Message
	err := json.Unmarshal([]byte(messageBody), &msg)
	if err != nil {
		log.Printf("Error unmarshaling message body: %v\n", err)
		return false
	}

	repoURL := msg.RepoURL
	buildCommand := msg.BuildCommand

	ctx := context.Background()

	utils.BuildProject(ctx, repoURL, buildCommand)

	// Return true if the message should be deleted, false otherwise
	return true
}

// Run listens to an SQS queue and processes messages.
func Run(queueURL string, workerType string) {
	sqsSvc := utils.GetSQSService()

	fmt.Printf("[Type: %s] Listening to SQS: %v\n\n", workerType, queueURL)

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
			messageStatus := ProcessMessage(message, workerType)
			if !messageStatus {
				log.Printf("Failed to process message [message id: %s]\n", *message.MessageId)

				// ! Testing: Invert this on PROD
				// Allowing delete for the invalid case to avoid clog

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
