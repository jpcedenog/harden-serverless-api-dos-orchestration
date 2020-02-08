package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type Response events.APIGatewayProxyResponse

type note struct {
	UserID    string `json:"userId"`
	ObjectKey string `json:"objectKey"`
}

func SaveNote(ctx context.Context, s3Event events.S3Event) error {
	fmt.Println("Saving Note...")

	for _, record := range s3Event.Records {
		fmt.Printf(
			"[%s - %s] Bucket = %s, Key = %s \n",
			record.EventSource,
			record.EventTime,
			record.S3.Bucket.Name,
			record.S3.Object.Key,
		)
		if err := SaveFileData(record); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func SaveFileData(record events.S3EventRecord) error {
	s3Object := record.S3
	myNote := note{
		UserID:    record.PrincipalID.PrincipalID,
		ObjectKey: s3Object.Object.Key,
	}

	dattr, err := dynamodbattribute.MarshalMap(myNote)
	if err != nil {
		return err
	}

	svc := dynamodb.New(session.Must(session.NewSession()))
	input := &dynamodb.PutItemInput{
		Item:      dattr,
		TableName: aws.String(os.Getenv("tableName")),
	}
	if _, err = svc.PutItem(input); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(SaveNote)
}
