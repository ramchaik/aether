package tests

import (
	"forge/internal/worker"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

// Define the MockSQSService struct
type MockSQSService struct {
	messages []sqs.Message
}

func TestProcessMessage(t *testing.T) {
	mockService := &MockSQSService{
		messages: []sqs.Message{
			{
				Body: aws.String("Test Message"),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"MessageType": {
						DataType:    aws.String("String"),
						StringValue: aws.String("Build"),
					},
				},
			},
		},
	}

	isProcessed := worker.ProcessMessage(&mockService.messages[0], "Build")
	assert.True(t, isProcessed, "Expected message to be processed")

	isProcessed = worker.ProcessMessage(&mockService.messages[0], "invalid-type")
	assert.False(t, isProcessed, "Expected message to be processed")
}
