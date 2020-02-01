package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func SaveNote(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

}

func main() {
	lambda.Start(SaveNote)
}
