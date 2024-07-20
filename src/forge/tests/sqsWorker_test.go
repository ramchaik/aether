package worker

import (
	"forge/internal/service"
	"forge/internal/worker"
	"log"
	"testing"

	pbProject "forge/internal/genprotobuf/project"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

type MockSQSService struct {
	message types.Message
}

type MockGrpcClient1 struct {
	conn string
}

func (g *MockGrpcClient1) UpdateProjectStatus(projectId string, status pbProject.ProjectStatus) {
	log.Println("projectId: ", projectId, " status", status)
}

type MockGrpcClient2 struct {
	conn string
}

func (g *MockGrpcClient2) PushLogs(projectId string, buildLogs []service.LogEntry) (bool, string) {
	log.Println("projectId: ", projectId)
	return true, "success"
}
func TestProcessMessage(t *testing.T) {
	// mock grpc client
	mockClient1 := &MockGrpcClient1{
		conn: "project-test",
	}

	mockClient2 := &MockGrpcClient2{
		conn: "logs-test",
	}
	// Mock SQS message
	mockMessage := types.Message{
		Body: aws.String(`{"repoURL": "https://github.com/example/repo", "buildCommand": "go build"}`),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Build"),
			},
		},
	}

	isProcessed := worker.ProcessMessage(mockMessage, "Build", mockClient1, mockClient2)
	assert.True(t, isProcessed, "Expected message to be processed")

	isProcessed = worker.ProcessMessage(mockMessage, "invalid-type", mockClient1, mockClient2)
	assert.False(t, isProcessed, "Expected message to be rejected due to invalid type")

	// Test invalid JSON message body
	invalidMessage := types.Message{
		Body: aws.String("Invalid JSON"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Build"),
			},
		},
	}

	isProcessed = worker.ProcessMessage(invalidMessage, "Build", mockClient1, mockClient2)
	assert.False(t, isProcessed, "Expected message with invalid JSON to be rejected")

	// Test missing message attributes
	missingAttributes := types.Message{
		Body: aws.String(`{"repoURL": "https://github.com/example/repo", "buildCommand": "go build"}`),
	}

	isProcessed = worker.ProcessMessage(missingAttributes, "Build", mockClient1, mockClient2)

	assert.False(t, isProcessed, "Expected message with missing attributes to be rejected")
}
