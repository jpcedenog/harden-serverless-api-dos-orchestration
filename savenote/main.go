package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func SaveNote(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("Request", request)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            request.Body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(SaveNote)
}
