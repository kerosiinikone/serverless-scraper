package queue

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

func New(s *session.Session) *sqs.SQS {
	return sqs.New(s)
}

func SendMessage(queue *sqs.SQS, msg *models.QueueMessage) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = queue.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonMsg)),
		QueueUrl:   aws.String(queueURL()),
	})
	return err
}

func ReceiveMessage(queue *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	return queue.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(queueURL()),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(60),
	})
}

func DeleteMessage(queue *sqs.SQS, receiptHandle *string) error {
	_, err := queue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:       aws.String(queueURL()),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return err
	}
	return nil 
}

func queueURL() string {
	return os.Getenv("QUEUE_URL")
}