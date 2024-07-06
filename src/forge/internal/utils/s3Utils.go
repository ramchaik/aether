package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3Client *s3.S3

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	s3Client = s3.New(sess)
}

func UploadFilesToS3(projectID string, files [][]byte) error {
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	// Create a new session (assumes AWS credentials are set up)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return err
	}

	// Create a new Uploader instance
	uploader := s3manager.NewUploader(sess)

	for _, fileContent := range files {
		key := fmt.Sprintf("build/%s/%s", projectID, "filename.ext")
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   bytes.NewReader(fileContent),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
