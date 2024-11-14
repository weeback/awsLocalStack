package main

import (
	"context"
	"log"
	"os"

	_ "app/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	DBTABLE_NAME  = "my-table"
	DBPRIMARY_KEY = "pkey"
	DBSORT_KEY    = "skey"
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

	// Create SQS service client
	svc := sqs.NewFromConfig(provider)

	// Create a new queue
	CreateQueue(svc)
	// List queues
	ListQueues(svc)
}

func CreateQueue(svc *sqs.Client) {
	// Create a new queue
	result, err := svc.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String("my-queue"),
	})
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	log.Printf("Queue URL: %s\n", *result.QueueUrl)

	// put message
	_, err = svc.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String("Hello World!"),
		QueueUrl:    result.QueueUrl,
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

}

func ListQueues(svc *sqs.Client) {
	// List queues
	result, err := svc.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		log.Fatalf("Failed to list queues: %v", err)
	}

	log.Println("Queues:")
	for _, qUrl := range result.QueueUrls {
		log.Printf("* %s\n", qUrl)

		// receive message
		result, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(qUrl),
		})
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		for _, msg := range result.Messages {
			log.Printf("  Message ID: %s\n", *msg.MessageId)
			log.Printf("  Message Body: %s\n", *msg.Body)
		}
	}

}
