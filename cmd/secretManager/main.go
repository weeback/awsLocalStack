package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

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
	svc := secretsmanager.NewFromConfig(provider)

	// Create a new secret
	CreateSecret(svc)
	// Update a secret
	UpdateSecret(svc)
	// List all secrets
	ListSecrets(svc)
}

func CreateSecret(svc *secretsmanager.Client) {
	// Create a secret
	_, err := svc.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("my-secret"),
		SecretString: aws.String("my-secret-value"),
	})
	if err != nil {
		log.Printf("unable to create secret, %v\n", err)
		return
	}

	log.Println("secret created")
}

func UpdateSecret(svc *secretsmanager.Client) {
	// Update a secret
	_, err := svc.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String("my-secret"),
		SecretString: aws.String("my-new-secret-value"),
	})
	if err != nil {
		log.Fatalf("unable to update secret, %v\n", err)
	}
	log.Println("secret updated")
}

func ListSecrets(svc *secretsmanager.Client) {
	// List all secrets
	resp, err := svc.ListSecrets(context.TODO(), &secretsmanager.ListSecretsInput{})
	if err != nil {
		log.Fatalf("unable to list secrets, %v", err)
	}
	for _, secret := range resp.SecretList {
		log.Printf("secret: %s", aws.ToString(secret.Name))
	}
}
