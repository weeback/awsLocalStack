package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

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

	LOG_GROUP_NAME  = "my-log-group"
	LOG_STREAM_NAME = "my-log-stream"
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
	// Create CloudWatch Logs client
	svc := cloudwatchlogs.NewFromConfig(provider)

	//
	PutLogEvents(svc)
	GetLogEvents(svc)
}

func PutLogEvents(logsSvc *cloudwatchlogs.Client) {
	//
	isResourceAlreadyExistsError := func(err error) bool {
		var (
			rex *types.ResourceAlreadyExistsException
		)
		return errors.As(err, &rex)
	}

	// Ensure log group exists
	_, err := logsSvc.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(LOG_GROUP_NAME),
	})
	if err != nil {
		if !isResourceAlreadyExistsError(err) {
			log.Println("Failed to create log group:", err)
			return
		}
	}

	// Ensure log stream exists
	_, err = logsSvc.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(LOG_GROUP_NAME),
		LogStreamName: aws.String(LOG_STREAM_NAME),
	})
	if err != nil {
		if !isResourceAlreadyExistsError(err) {
			log.Println("Failed to create log stream:", err)
			return
		}
	}

	// Put log event
	_, err = logsSvc.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String("This is a log message written by the CloudWatch Logs API at " + time.Now().String()),
				Timestamp: aws.Int64(time.Now().Unix() * 1000),
			},
		},
		LogGroupName:  aws.String(LOG_GROUP_NAME),
		LogStreamName: aws.String(LOG_STREAM_NAME),
	})
	if err != nil {
		log.Println("Failed to put log events:", err)
		return
	}

	log.Println("Successfully put log event")
}

func GetLogEvents(logsSvc *cloudwatchlogs.Client) {
	// Get log events
	output, err := logsSvc.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(LOG_GROUP_NAME),
		LogStreamName: aws.String(LOG_STREAM_NAME),
	})
	if err != nil {
		fmt.Println("Failed to get log events:", err)
		return
	}

	fmt.Println("Log events:")
	for _, event := range output.Events {
		fmt.Printf("Timestamp: %s, Message: %s\n", time.UnixMilli(*event.Timestamp).Format(time.DateTime), *event.Message)
	}
}
