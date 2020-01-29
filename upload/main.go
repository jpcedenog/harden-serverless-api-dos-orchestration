package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"os"
)

type Response events.APIGatewayProxyResponse

type Event struct {
	ControlValue string `json:"control_value"`
}

func GuessSecret(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	apiEvent := Event{}
	if err := json.Unmarshal([]byte(request.Body), &apiEvent); err != nil {
		return Response{StatusCode: 400}, err
	}

	kmsSvc := kms.New(session.Must(session.NewSession()))
	//Get encrypted password from environment
	encryptedPassword, err := b64.URLEncoding.DecodeString(os.Getenv("password"))
	if err != nil {
		return Response{StatusCode: 400}, err
	}
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedPassword,
	}

	//Decode encrypted password
	result, err := kmsSvc.Decrypt(input)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	//Do something with plain text password value
	isSecretValueCorrect := apiEvent.ControlValue == string(result.Plaintext)

	body, err := json.Marshal(map[string]interface{}{
		"is_secret_correct": isSecretValueCorrect,
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
			"Content-Type":           "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(GuessSecret)
}
