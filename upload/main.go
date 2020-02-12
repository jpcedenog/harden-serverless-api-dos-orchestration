package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	maxSize = 1024 * 1000
)

type Response events.APIGatewayProxyResponse

type Event struct {
	FileUrL string `json:"file_url"`
}

func UploadFile(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("Uploading file...")

	myClient := httpClient()

	apiEvent := Event{}
	if err := json.Unmarshal([]byte(request.Body), &apiEvent); err != nil {
		return Response{StatusCode: 400}, err
	}

	fileURL, err := url.ParseRequestURI(apiEvent.FileUrL)
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	file, err := myClient.Get(fileURL.String())
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	fileSize := file.ContentLength
	if fileSize > maxSize {
		return Response{StatusCode: 400}, errors.New("file exceeds maximum size")
	}

	fileBuffer, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	contentType := http.DetectContentType(fileBuffer[:513])
	if contentType != "image/jpeg" {
		return Response{StatusCode: 400}, errors.New("file is not an image")
	}

	if err := file.Body.Close(); err != nil {
		return Response{StatusCode: 400}, err
	}

	fileName := strings.TrimSpace(path.Base(file.Request.URL.Path))
	if len(fileName) == 0 {
		return Response{StatusCode: 400}, errors.New("file name is missing")
	}

	cognitoIdentityID := request.RequestContext.Identity.CognitoIdentityID
	os.Setenv("cognitoIdentityID", cognitoIdentityID)

	uploadInput := &s3manager.UploadInput{
		Body:   bytes.NewReader(fileBuffer),
		Bucket: aws.String(os.Getenv("bucketName")),
		Key:    aws.String(strings.Join([]string{cognitoIdentityID, fileName}, "/")),
	}

	fmt.Println("Upload Input:", uploadInput)

	uploader := s3manager.NewUploader(session.Must(session.NewSession()))
	result, err := uploader.Upload(uploadInput)
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"status": result.Location,
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

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func main() {
	lambda.Start(UploadFile)
}
