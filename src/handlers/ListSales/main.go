package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Sale struct {
	ID              string `json:"id"`
	UserId          string `json:"user_id"`
	StoreId         string `json:"store_id"`
	Date            string `json:"date"`
	SalesInCents    int    `json:"sales_in_cents"`
	ExpensesInCents int    `json:"expenses_in_cents"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type Response struct {
	Sales []Sale `json:"sales"`
}

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		log.Println("Failed to connect")
	} else {
		ddb = dynamodb.New(session) // Create DynamoDB client
	}
}

func CreateSale(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("CreateSale")

	storeId := request.RequestContext.Authorizer["SF-Store-Id"].(string)

	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":store_id": {
				S: aws.String(storeId),
			},
		},
		KeyConditionExpression: aws.String("store_id = :store_id"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String(os.Getenv("STORE_INDEX_NAME")),
	}

	if result, err := ddb.Query(queryInput); err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		sales := []Sale{}
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &sales)

		if err != nil {
			return events.APIGatewayProxyResponse{ // Error HTTP response
				Body:       err.Error(),
				StatusCode: 500,
			}, nil
		}

		body, _ := json.Marshal(&Response{
			Sales: sales,
		})

		return events.APIGatewayProxyResponse{ // Success HTTP response
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(CreateSale)
}
