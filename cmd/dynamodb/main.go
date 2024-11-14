package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

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

	DBTABLE_NAME  = "my-table"
	DBPRIMARY_KEY = "pkey"
	DBSORT_KEY    = "skey"
)

func main() {

	// // Create a new session with the LocalStack endpoint
	// provider, err := session.NewSession(&aws.Config{
	// 	Endpoint: aws.String(AWS_ENDPOINT),
	// 	Credentials: credentials.NewStaticCredentials(
	// 		AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, "",
	// 	),
	// })
	// if err != nil {
	// 	fmt.Println("Failed to create provider session:", err)
	// 	return
	// }

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

	// Create DynamoDB service client
	svc := dynamodb.NewFromConfig(provider)

	// Create a new table
	CreateTable(svc)

	//
	Insert(svc)

	// Scan the table
	Scan(svc)

	// Query the table
	Query(svc)
}

func CreateTable(svc *dynamodb.Client) {

	// Check if table exists
	_, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(DBTABLE_NAME),
	})
	if err == nil {
		fmt.Println("Table already exists")
		return
	}
	// Create a new table
	_, err = svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(DBTABLE_NAME),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(DBPRIMARY_KEY),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(DBSORT_KEY),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(DBPRIMARY_KEY),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(DBSORT_KEY),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
	fmt.Println("Table created successfully")
}

func Insert(svc *dynamodb.Client) {
	_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("my-table"),
		Item: map[string]types.AttributeValue{
			DBPRIMARY_KEY: &types.AttributeValueMemberS{
				Value: "my-partition-key",
			},
			DBSORT_KEY: &types.AttributeValueMemberS{
				Value: "my-sort-key-" + time.Now().Format(time.RFC3339),
			},
			"attribute": &types.AttributeValueMemberS{
				Value: "my-attribute-value",
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to insert item:", err)
		return
	}
	fmt.Println("Item inserted successfully")
}

func Scan(svc *dynamodb.Client) {
	resp, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(DBTABLE_NAME),
	})
	if err != nil {
		fmt.Println("Failed to scan table:", err)
		return
	}
	fmt.Println("Scan results:")

	// parse the results
	var items []struct {
		PK    string `dynamodbav:"pkey"`
		SK    string `dynamodbav:"skey"`
		Value string `dynamodbav:"attribute"`
	}
	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		fmt.Println("Failed to unmarshal scan items:", err)
		return
	}
	for i, item := range items {
		fmt.Printf("Item(%d): %+v\n", i, item)
	}
}

func Query(svc *dynamodb.Client) {
	expr, err := attributevalue.MarshalMap(map[string]interface{}{
		":pkey": "my-partition-key",
		":val":  "my-sort-key",
	})
	if err != nil {
		fmt.Println("Failed to marshal query expression:", err)
		return
	}
	resp, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(DBTABLE_NAME),
		KeyConditionExpression:    aws.String(fmt.Sprintf("%s = :pkey AND begins_with(%s, :val)", DBPRIMARY_KEY, DBSORT_KEY)),
		ExpressionAttributeValues: expr,
		//  map[string]types.AttributeValue{
		// 	":pkey": &types.AttributeValueMemberS{
		// 		Value: "my-partition-key",
		// 	},
		// 	":val": &types.AttributeValueMemberS{
		// 		Value: "my-sort-key",
		// 	},
		// },
		// 	// or use attributevalue.MarshalMap(object)
	})
	if err != nil {
		fmt.Println("Failed to query table:", err)
		return
	}
	fmt.Println("Query results:")
	// parse the results
	var items []struct {
		PK    string `dynamodbav:"pkey"`
		SK    string `dynamodbav:"skey"`
		Value string `dynamodbav:"attribute"`
	}
	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		fmt.Println("Failed to unmarshal scan items:", err)
		return
	}
	for i, item := range items {
		fmt.Printf("Item(%d): %+v\n", i, item)
	}
}
