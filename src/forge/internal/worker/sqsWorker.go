package worker

import (
	"context"
	"encoding/json"
	"fmt"
	pb "forge/internal/genprotobuf/project"
	"forge/internal/service"
	"forge/internal/utils"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Message struct {
	ProjectId    string `json:"projectId"`
	RepoURL      string `json:"repoURL"`
	BuildCommand string `json:"buildCommand"`
}

// ProcessMessage takes a message and performs the necessary actions based on the message content.
func ProcessMessage(
	message types.Message,
	workerType string,
	projectService service.ProjectService,
	logService service.ProjectLogService,
) bool {
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

	projectId := msg.ProjectId
	repoURL := msg.RepoURL
	buildCommand := msg.BuildCommand

	ctx := context.Background()

	// Push log entry
	pushLogs := func(logMessage string) {
		logEntry := service.LogEntry{
			Log:       logMessage,
			Timestamp: time.Now().Unix(),
		}
		logService.PushLogs(projectId, logEntry)
	}

	cli, buildDir, imageName, err := utils.BuildProject(ctx, repoURL, buildCommand, pushLogs)
	if err != nil {
		log.Fatalf("Failed to build project: %v", err)
	}

	// Deploying to S3
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	prefix := fmt.Sprintf("projects/%s/build/", projectId)

	s3Client, err := utils.GetS3Service()
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	if err := utils.UploadToS3(ctx, buildDir, bucketName, prefix, s3Client); err != nil {
		log.Fatalf("Failed to upload files to S3: %v", err)
	}

	utils.Cleanup(ctx, cli, buildDir, imageName)

	// Update launchpad as the project is deployed
	projectService.UpdateProjectStatus(projectId, pb.ProjectStatus_LIVE)
	return true
}

// Run listens to an SQS queue and processes messages.
func Run(
	queueURL string,
	workerType string,
	projectService service.ProjectService,
	logService service.ProjectLogService,
) {
	sqsSvc, err := utils.GetSQSService()
	if err != nil {
		log.Fatalf("Failed to get SQS service %v", err)
	}

	fmt.Printf("[Type: %s] Listening to SQS: %v\n", workerType, queueURL)

	for {
		result, err := sqsSvc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:              &queueURL,
			MaxNumberOfMessages:   10,
			VisibilityTimeout:     300,
			WaitTimeSeconds:       20,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			log.Fatalf("ReceiveMessage failed %v", err)
		}

		for _, message := range result.Messages {
			messageStatus := ProcessMessage(message, workerType, projectService, logService)
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
