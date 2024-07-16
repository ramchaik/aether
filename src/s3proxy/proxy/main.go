package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	bucketName = "aether-bucket"
	buildDir   = "build"
	defaultKey = "index.html"
	bucketURL  = "https://aether-bucket.s3.amazonaws.com"
)

var s3Client *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	s3Client = s3.NewFromConfig(cfg)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract projectId from the path parameters
	projectId := request.PathParameters["projectId"]
	if projectId == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Project ID is required",
		}, nil
	}

	// Construct the key using the projectId
	key := path.Join("projects", projectId, buildDir)
	if request.Path == "/"+projectId || request.Path == "/"+projectId+"/" {
		key = path.Join(key, defaultKey)
	} else {
		key = path.Join(key, strings.TrimPrefix(request.Path, "/"+projectId+"/"))
	}

	// Check if the object exists in S3
	_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var nske *types.NoSuchKey
		if strings.Contains(err.Error(), "NotFound") || errors.As(err, &nske) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not Found",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	// Construct the full S3 bucket URL
	redirectURL := fmt.Sprintf("%s/%s", bucketURL, key)

	return events.APIGatewayProxyResponse{
		StatusCode: 307, // Temporary Redirect
		Headers: map[string]string{
			"Location":      redirectURL,
			"Cache-Control": "public, max-age=31536000", // Cache for 1 year
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
