package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"

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
	svc := scheduler.NewFromConfig(provider)

	// Create a new scheduler
	CreateScheduler(svc)

	//
	ListScheduler(svc)
}

func CreateScheduler(svc *scheduler.Client) {

	// Create Group if not exists
	if _, err := svc.CreateScheduleGroup(context.TODO(), &scheduler.CreateScheduleGroupInput{
		Name: aws.String("default"),
	}); err != nil {
		log.Fatalf("unable to create scheduler group, %v", err)
	}
	// Create a new scheduler
	if _, err := svc.CreateSchedule(context.TODO(), &scheduler.CreateScheduleInput{
		GroupName:          aws.String("default"),
		Name:               aws.String("my-scheduler"),
		Target: &types.Target{
			Arn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:my-function"),
			RoleArn: aws.String("arn:aws:iam::000000000000:role/service-role/MySchedulerRole"),
		},
		FlexibleTimeWindow: &types.FlexibleTimeWindow{
			Mode: types.FlexibleTimeWindowModeOff,
		},
		Description:        aws.String("my scheduler"),
		ScheduleExpression: aws.String("rate(1 minute)"),
	}); err != nil {
		log.Fatalf("unable to create scheduler, %v", err)
	}
}

func ListScheduler(svc *scheduler.Client) {
	sgl, err := svc.ListScheduleGroups(context.TODO(), &scheduler.ListScheduleGroupsInput{})
	if err != nil {
		log.Fatalf("unable to list scheduler groups, %v", err)
	}
	for _, sg := range sgl.ScheduleGroups {
		log.Printf("Group: %v\n", aws.ToString(sg.Name))
		// ListSchedulers
		result, err := svc.ListSchedules(context.TODO(), &scheduler.ListSchedulesInput{
			GroupName: sg.Name,
		})
		if err != nil {
			log.Fatalf("unable to list schedulers, %v", err)
		}

		for _, s := range result.Schedules {
			log.Printf("Scheduler: %v created on %s\n", aws.ToString(s.Name), s.CreationDate)
		}
	}
}
