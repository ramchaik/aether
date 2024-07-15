package utils

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(ctx context.Context, buildDir, bucketName, prefix string, s3Client *s3.Client) error {
	return filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %q: %v", path, err)
		}

		// Skip the root directory itself
		if path == buildDir {
			return nil
		}

		// Create the S3 key (path within the bucket)
		relPath, err := filepath.Rel(buildDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}
		s3Key := filepath.Join(prefix, relPath)

		// Replace backslashes with forward slashes for S3 paths
		s3Key = strings.ReplaceAll(s3Key, "\\", "/")

		if info.IsDir() {
			// For directories, we create an empty object with a trailing slash
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(s3Key + "/"),
				ContentType: aws.String("application/x-directory"),
			})
			if err != nil {
				return fmt.Errorf("failed to create directory %s in S3: %v", s3Key, err)
			}
			fmt.Printf("Created directory: s3://%s/%s\n", bucketName, s3Key)
		} else {
			// For files, upload the content with the correct content type
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %v", path, err)
			}
			defer file.Close()

			// Detect content type
			contentType := detectContentType(path)

			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(s3Key),
				Body:        file,
				ContentType: aws.String(contentType),
			})
			if err != nil {
				return fmt.Errorf("failed to upload file %s to S3: %v", path, err)
			}
			fmt.Printf("Uploaded file: s3://%s/%s (Content-Type: %s)\n", bucketName, s3Key, contentType)
		}

		return nil
	})
}

func detectContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	if ext != "" {
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			return mimeType
		}
	}

	// If mime type is not found, try to detect it by reading the file
	file, err := os.Open(filePath)
	if err != nil {
		return "application/octet-stream" // default to binary data
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "application/octet-stream" // default to binary data
	}

	// DetectContentType always returns a valid MIME type
	return http.DetectContentType(buffer)
}
