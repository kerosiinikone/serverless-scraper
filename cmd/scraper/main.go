package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	"github.com/kerosiinikone/serverless-scraper/infra"
	"github.com/kerosiinikone/serverless-scraper/infra/blob"
	"github.com/kerosiinikone/serverless-scraper/infra/database"
	"github.com/kerosiinikone/serverless-scraper/infra/queue"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper/reddit"
)


const (
	delayMax = 2
    delayMin = 1
    maxDepth = 100
)

var (
    storage *s3manager.Uploader
    sess 	*session.Session
    q 		*sqs.SQS
	db 		*dynamodb.DynamoDB
)

func init() {
	if os.Getenv("RUNTIME") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	sess := infra.New(nil)
	storage = blob.NewUploader(sess)
    q = queue.New(sess)
	db = database.New(sess)
}

func handleScrape(ctx context.Context, event json.RawMessage) error {
	var (
		c = make(chan struct{})
		failed = make(chan error)
		input scraper.APIRequest
	)
	
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Minute))
	defer cancel()

	cfg := &scraper.Config{
		DelayMax: delayMax,
		DelayMin: delayMin,
		MaxDepth: maxDepth,
	}

	if err := json.Unmarshal(event, &input); err != nil {
		err = database.SetStatusFailed(db, input.ID, input.ClientID)
		return err
	}

	rs := reddit.New(cfg, storage, q)

	go func() {
		if err := rs.Scrape(ctx, &input, c); err != nil {
			log.Println("Error scraping: ", err)
			failed <- err
		}
	}()

	select {
	case <- c:
		log.Println("Finished scraping")
	case fail := <- failed:
		if err := database.SetStatusFailed(db, input.ID, input.ClientID); err != nil {
			return err
		}
		return fail
	}
	return nil
}

func main() {
	lambda.Start(handleScrape)
}