package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	"github.com/kerosiinikone/serverless-scraper/infra"
	"github.com/kerosiinikone/serverless-scraper/infra/blob"
	"github.com/kerosiinikone/serverless-scraper/infra/database"
	"github.com/kerosiinikone/serverless-scraper/infra/queue"
	consumer "github.com/kerosiinikone/serverless-scraper/internal/consumer"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/sashabaranov/go-openai"
)

var (
	downloader  *s3manager.Downloader
	sess    	*session.Session
	q   		*sqs.SQS
	s3svc     	*s3.S3
	client  	*openai.Client
	db      	*dynamodb.DynamoDB
)

func init() {
	if os.Getenv("RUNTIME") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	sess = infra.New(nil)
	downloader = blob.NewDownloader(sess)
	q = queue.New(sess)
	s3svc = blob.New(sess)
	db = database.New(sess)
	client = openai.NewClient(os.Getenv("OPENAI_API"))
}

func handleAnalysis(ctx context.Context, sqsEvent events.SQSEvent) error {
	var (
		msg models.QueueMessage
	)
	message := sqsEvent.Records[0].Body
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return err
	}
	if msg.ClientID == "" || msg.RequestID == "" {
		return nil
	}
	
	c := consumer.New(q, downloader, s3svc, db, client, &msg)

	f, err := c.Get(&msg)
	if err != nil {
		return err
	}
	resp, err := c.Synthesis(f)
	if err != nil {
		return err
	}
	if err = c.Store(resp); err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(handleAnalysis)
}
