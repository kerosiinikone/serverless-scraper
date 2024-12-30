package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	tableName = "Request"
)

func New(s *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(s)
}

func CreateTable(db *dynamodb.DynamoDB) error {
    _, err := db.DescribeTable(&dynamodb.DescribeTableInput{
        TableName: aws.String(tableName),
    })
    if err == nil {
        return nil
    }
    input := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition{
            {
                AttributeName: aws.String("RequestID"),
                AttributeType: aws.String("S"),
            },
            {
                AttributeName: aws.String("ClientID"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement{
            {
                AttributeName: aws.String("RequestID"),
                KeyType:       aws.String("HASH"),
            },
            {
                AttributeName: aws.String("ClientID"),
                KeyType:       aws.String("RANGE"),
            },
        },
        ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
        TableName: aws.String(tableName),
    }
	_, err = db.CreateTable(input)
    if err != nil {
        log.Println(err)
        return err
    }
    err = db.WaitUntilTableExists(&dynamodb.DescribeTableInput{
        TableName: aws.String(tableName),
    })
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func CreateEmptyItem(db *dynamodb.DynamoDB, requestID string, clientID string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"RequestID": {
				S: aws.String(requestID),
			},
			"ClientID": {
				S: aws.String(clientID),
			},
			"Content": {
				S: aws.String(""),
			},
			"CurrentStatus": {
				S: aws.String("Pending"),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := db.PutItem(input)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SetStatusFailed(db *dynamodb.DynamoDB, requestID string, clientID string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String("Failed"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"RequestID": {
				S: aws.String(requestID),
			},
			"ClientID": {
				S: aws.String(clientID),
			},
		},
		TableName: aws.String(tableName),
		UpdateExpression: aws.String("SET CurrentStatus = :s"),
	}
	_, err := db.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}

func SetItemContent(db *dynamodb.DynamoDB, requestID string, clientID string, content string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(content),
			},
			":s": {
				S: aws.String("Completed"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"RequestID": {
				S: aws.String(requestID),
			},
			"ClientID": {
				S: aws.String(clientID),
			},
		},
		TableName: aws.String(tableName),
		UpdateExpression: aws.String("SET Content = :c, CurrentStatus = :s"),
	}
	_, err := db.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}