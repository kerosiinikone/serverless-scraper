package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kerosiinikone/serverless-scraper/infra"
	"github.com/kerosiinikone/serverless-scraper/infra/database"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
)

type Response events.APIGatewayProxyResponse

type Request events.APIGatewayProxyRequest

var (
	sess *session.Session
	db *dynamodb.DynamoDB
)

func init() {
	if os.Getenv("RUNTIME") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	sess = infra.New(nil)
	db = database.New(sess)
}

func handleAPIRequest(req Request) (Response, error) {
    var (
		payload scraper.APIRequest
		lambdaClient = lambdaService.New(sess)
	)

	rId := uuid.New().String()
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return Response{StatusCode: http.StatusBadRequest, Body: err.Error()}, err
	}
	payload.ID = rId

	// Separate AWS logic from handlers
	if err := database.CreateTable(db); err != nil {
		return Response{StatusCode: http.StatusInternalServerError, Body: err.Error()}, err
	}
	if err := database.CreateEmptyItem(db, payload.ID, payload.ClientID); err != nil {
		return Response{StatusCode: http.StatusInternalServerError, Body: err.Error()}, err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError, Body: err.Error()}, err
	}

    _, err = lambdaClient.Invoke(&lambdaService.InvokeInput{
        FunctionName: aws.String("scraper-lambda"),
		InvocationType: aws.String("Event"),
		Payload: payloadBytes,
    })
    if err != nil {
        return Response{StatusCode: http.StatusInternalServerError, Body: err.Error()}, err
    }

	return Response{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            rId,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handleAPIRequest)
}
