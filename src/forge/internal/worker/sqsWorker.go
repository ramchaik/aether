package worker

import (
	"fmt"
	"forge/internal/utils"
	"log"

	"github.com/aws/aws-sdk-go/aws"
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
	sqsSvc := utils.GetSQSService()

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
