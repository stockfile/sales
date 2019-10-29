package db

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		log.Println("Failed to connect")
	} else {
		db = dynamodb.New(session)
	}
}

func GetDB() *dynamodb.DynamoDB {
	return db
}

func PutItem(in interface{}, tableName *string) error {
	item, _ := dynamodbattribute.MarshalMap(in)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}

	_, err := db.PutItem(input)

	return err
}

func Query(queryInput *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return db.Query(queryInput)
}
