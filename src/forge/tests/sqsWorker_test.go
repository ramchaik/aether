package worker

import (
	"forge/internal/worker"
	"log"
	"testing"

	pb "forge/internal/genprotobuf"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

type MockSQSService struct {
	message types.Message
}

type MockGrpcClient struct {
	conn string
}

func (g *MockGrpcClient) UpdateProjectStatus(projectId string, status pb.ProjectStatus) {
	log.Println("projectId: ", projectId, " status", status)
}

func TestProcessMessage(t *testing.T) {
	// mock grpc client
	mockClient := &MockGrpcClient{
		conn: "test",
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

	isProcessed := worker.ProcessMessage(mockMessage, "Build", mockClient)
	assert.True(t, isProcessed, "Expected message to be processed")

	isProcessed = worker.ProcessMessage(mockMessage, "invalid-type", mockClient)
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

	isProcessed = worker.ProcessMessage(invalidMessage, "Build", mockClient)
	assert.False(t, isProcessed, "Expected message with invalid JSON to be rejected")

	// Test missing message attributes
	missingAttributes := types.Message{
		Body: aws.String(`{"repoURL": "https://github.com/example/repo", "buildCommand": "go build"}`),
	}

	isProcessed = worker.ProcessMessage(missingAttributes, "Build", mockClient)
	assert.False(t, isProcessed, "Expected message with missing attributes to be rejected")
}
