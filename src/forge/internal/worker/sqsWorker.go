package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"forge/internal/utils"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type Message struct {
	// ProjectId    string `json:"projectId"`
	RepoURL      string `json:"repoURL"`
	BuildCommand string `json:"buildCommand"`
}

// ProcessMessage takes a message and performs the necessary actions based on the message content.
func ProcessMessage(message types.Message, workerType string) bool {
	if message.Body == nil || message.MessageAttributes == nil {
		log.Println("Invalid message received")
		return false
	}

	messageBody := *message.Body
	messageType := *message.MessageAttributes["MessageType"].StringValue

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

	cli, buildDir, imageName, err := utils.BuildProject(ctx, repoURL, buildCommand)
	if err != nil {
		log.Fatalf("Failed to build project: %v", err)
	}

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	// TODO: use the projectID for instead of uuid
	prefix := fmt.Sprintf("projects/%s/build/", uuid.New().String())

	s3Client, err := utils.GetS3Service()
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	if err := utils.UploadToS3(ctx, buildDir, bucketName, prefix, s3Client); err != nil {
		log.Fatalf("Failed to upload files to S3: %v", err)
	}

	utils.Cleanup(ctx, cli, buildDir, imageName)

	return true
}

// Run listens to an SQS queue and processes messages.
func Run(queueURL string, workerType string) {
	sqsSvc, err := utils.GetSQSService()
	if err != nil {
		log.Fatalf("Failed to get SQS service %v", err)
	}

	fmt.Printf("[Type: %s] Listening to SQS: %v\n", workerType, queueURL)

	for {
		result, err := sqsSvc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:              &queueURL,
			MaxNumberOfMessages:   10,
			VisibilityTimeout:     20,
			WaitTimeSeconds:       0,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			log.Fatalf("ReceiveMessage failed %v", err)
		}

		for _, message := range result.Messages {
			messageStatus := ProcessMessage(message, workerType)
			if !messageStatus {
				log.Printf("Failed to process message [message id: %s]\n", *message.MessageId)
				break
			}

			// Delete the message from the queue after processing
			_, err = sqsSvc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("DeleteMessage Failed %v", err)
			}
		}
	}
}
