package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if req != nil {
		log.Printf("%+v\n", map[string]any{
			"RequestID":             req.RequestContext.RequestID,
			"HTTPMethod":            req.HTTPMethod,
			"Path":                  req.Path,
			"QueryStringParameters": req.QueryStringParameters,
			"Headers":               req.Headers,
			"Body":                  req.Body,
			"IsBase64Encoded":       req.IsBase64Encoded,
		})
	} else {
		log.Printf("%+v\n", map[string]any{"RequestID": "Unknown"})
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, Lambda!",
	}, nil
}

func main() {
	log.SetPrefix("")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lambda.Start(handler)
}
