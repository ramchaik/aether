package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var inMemoryDB = map[string]string{
	"abc": "abc-prefix",
	"xyz": "xyz-prefix",
}

const (
	bucketName = "aether-bucket"
	buildDir   = "build"
	defaultKey = "index.html"
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
	host := request.Headers["Host"]
	baseDomain := os.Getenv("BASE_DOMAIN")
	subdomain := strings.TrimSuffix(host, "."+baseDomain)

	if prefix, ok := inMemoryDB[subdomain]; ok {
		key := path.Join(prefix, buildDir)
		if request.Path == "/" || request.Path == "" {
			key = path.Join(key, defaultKey)
		} else {
			key = path.Join(key, strings.TrimPrefix(request.Path, "/"))
		}

		// Generate pre-signed URL
		presignClient := s3.NewPresignClient(s3Client)
		presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 3600 * 24 // 24 hours
		})
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Failed to generate pre-signed URL",
			}, nil
		}

		// Parse the pre-signed URL to get the query string
		parsedURL, err := url.Parse(presignResult.URL)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Failed to parse pre-signed URL",
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 307, // Temporary Redirect
			Headers: map[string]string{
				"Location":      fmt.Sprintf("/%s?%s", key, parsedURL.RawQuery),
				"Cache-Control": "public, max-age=31536000", // Cache for 1 year
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       fmt.Sprintf("Subdomain not found: %s", subdomain),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
