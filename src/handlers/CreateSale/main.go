package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	uuid "github.com/satori/go.uuid"
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
	Sale Sale `json:"sale"`
}

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session) // Create DynamoDB client
	}
}

func CreateSale(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("CreateSale")

	userId := request.RequestContext.Authorizer["SF-User-Id"].(string)
	storeId := request.RequestContext.Authorizer["SF-Store-Id"].(string)

	var (
		id        = uuid.Must(uuid.NewV4(), nil).String()
		tableName = aws.String(os.Getenv("STORES_TABLE_NAME"))
		timestamp = time.Now().Format(time.RFC3339)
	)

	sale := &Sale{
		ID:        id,
		UserId:    userId,
		StoreId:   storeId,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// Parse request body
	json.Unmarshal([]byte(request.Body), sale)

	// Write to DynamoDB
	item, _ := dynamodbattribute.MarshalMap(sale)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}

	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(&Response{
			Sale: *sale,
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
