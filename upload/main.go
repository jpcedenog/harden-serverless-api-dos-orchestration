package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"time"
)

type Response events.APIGatewayProxyResponse

type Event struct {
	FileUrL string `json:"file_url"`
}

func UploadFile(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	apiEvent := Event{}
	if err := json.Unmarshal([]byte(request.Body), &apiEvent); err != nil {
		return Response{StatusCode: 400}, err
	}

	var myClient = &http.Client{Timeout: 5 * time.Second}
	file, err := myClient.Get(apiEvent.FileUrL)
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	defer file.Body.Close()

	body, err := json.Marshal(map[string]interface{}{
		"is_secret_correct": false,
	})
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(UploadFile)
}
