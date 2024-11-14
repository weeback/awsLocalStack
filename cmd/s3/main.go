package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	// import the env package to load the environment variables
	_ "app/env"
)

var (
	// AWS_REGION is the region to use
	AWS_REGION = os.Getenv("LOCALSTACK_DEFAULT_REGION")
	// AWS_ACCESS_KEY_ID is the access key ID
	AWS_ACCESS_KEY_ID = os.Getenv("LOCALSTACK_ACCESS_KEY_ID")
	// AWS_SECRET_ACCESS_KEY is the secret access key
	AWS_SECRET_ACCESS_KEY = os.Getenv("LOCALSTACK_SECRET_ACCESS_KEY")
	// AWS_ENDPOINT is the endpoint for LocalStack
	AWS_ENDPOINT = os.Getenv("LOCALSTACK_ENDPOINT")

	S3_BUCKET_NAME = "my-bucket"
)

func main() {

	// Create a new session with the LocalStack endpoint
	// Load the Shared AWS Configuration (~/.aws/config). Replace with the LocalStack endpoint
	provider, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
		config.WithBaseEndpoint(AWS_ENDPOINT),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
			),
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// Create S3 service client
	svc := s3.NewFromConfig(provider)

	// PutItem
	PutItem(svc)
	// ListBuckets
	ListBuckets(svc)
}

func PutItem(svc *s3.Client) {

	isResourceNotFoundError := func(err error) bool {
		var (
			rnf *types.NotFound
		)
		return errors.As(err, &rnf)
	}
	// Check if bucket exists
	_, err := svc.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(S3_BUCKET_NAME),
	})
	if isResourceNotFoundError(err) {
		// Create bucket
		if _, err := svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(S3_BUCKET_NAME),
		}); err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	} else if err != nil {
		log.Println("Failed to head bucket:", err)
		return
	}
	// Put item
	if _, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3_BUCKET_NAME),
		Key:    aws.String("my-key"),
		Body:   bytes.NewReader([]byte("Hello, World!")),
	}); err != nil {
		log.Fatalf("Failed to put item: %v", err)
	}
}

func ListBuckets(svc *s3.Client) {
	// List buckets
	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Failed to list buckets: %v", err)
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n", aws.ToString(b.Name), b.CreationDate)

		// List objects
		objects, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: b.Name,
		})
		if err != nil {
			log.Printf("Failed to list objects: %v", err)
			continue
		} else {
			for _, o := range objects.Contents {
				fmt.Printf("* %s created on %s\n", aws.ToString(o.Key), o.LastModified)
			}
		}
	}

}
