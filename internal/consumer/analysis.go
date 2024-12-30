package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kerosiinikone/serverless-scraper/infra/blob"
	"github.com/kerosiinikone/serverless-scraper/infra/database"
	"github.com/kerosiinikone/serverless-scraper/infra/queue"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/sashabaranov/go-openai"
)

type Consumer struct {
	q *sqs.SQS
	d *s3manager.Downloader
	s *s3.S3
	db *dynamodb.DynamoDB
	ai *openai.Client

	msg *models.QueueMessage
}

type Message struct {
	Body     string `json:"body"`
	Comments []models.ReplyTree `json:"comments"`
}

// Envelope for the messages
type Messages struct {
	Messages []Message `json:"messages"`
}

func New(q *sqs.SQS, d *s3manager.Downloader, s *s3.S3, db *dynamodb.DynamoDB, ai *openai.Client, msg *models.QueueMessage) *Consumer {
	return &Consumer{
		q: q,
		d: d,
		s: s,
		db: db,
		ai: ai,
		msg: msg,
	}
}

func (c *Consumer) Get(sqsMsg *models.QueueMessage) (string, error) {
	msgs, err := blob.RetrieveFiles(c.s, c.d, sqsMsg.ClientID, sqsMsg.RequestID)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	final, err := formatMessages(msgs)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	return string(final), nil
}

func (c *Consumer) Synthesis(final string) (string, error) {
	resp, err := c.ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Analyze the messages in JSON and produce a synthesis on the problems and pain points users are experiencing. Emphasize the problems that could be solved with a startup without giving solutions to those problems. Filter out problems regarding politics or sexuality. The response should be a JSON object with the synthesis. Messages: %s", final),
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: "json_schema",
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "Synthesis",
					Schema: Schema,
					Strict: true,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Consumer) Store(content string) error {
	if content == "" {
		return fmt.Errorf("No content in the response")
	}
	err := database.SetItemContent(c.db, c.msg.RequestID, c.msg.ClientID, content)
	if err != nil {
		return err
	}
	return nil
}

func (c *Consumer) _retrieveQueueMessage() (*models.QueueMessage ,error) {
	var (
		message *models.QueueMessage
	)
	msg, err := queue.ReceiveMessage(c.q)
	if err != nil {
		return nil, err
	}
	if len(msg.Messages) == 0 {
		return nil, fmt.Errorf("No messages in the queue")
	}
	m := msg.Messages[0]
	defer func() {
		if err := queue.DeleteMessage(c.q, m.ReceiptHandle); err != nil {
			log.Println(err)
		}	
	}() 
	if err := json.Unmarshal([]byte(*m.Body), &message); err != nil {
		return nil, err
	}
	
	c.msg = message
	return message, nil
}

func formatMessages(msgs []models.DataEntry) ([]byte, error) {
	var (
		messages Messages
	)
	for _, msg := range msgs {
		messages.Messages = append(messages.Messages, Message{Body: msg.Post.Selftext, Comments: msg.Post.Comments})
	}
	return json.Marshal(messages)
}
