package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stockfile/sales/src/db"
	"github.com/stockfile/sales/src/models"
	"github.com/stockfile/sales/src/renderer"
)

type Response struct {
	Sales []models.Sale `json:"sales"`
}

func ListSales(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storeID := request.RequestContext.Authorizer["SF-Store-Id"].(string)
	sales := querySalesByStoreId(storeID)

	body, _ := json.Marshal(&Response{
		Sales: sales,
	})
	return renderer.RenderSuccess(body)
}

func querySalesByStoreId(storeID string) []models.Sale {
	scanIndexForward := false
	sales := []models.Sale{}

	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":store_id": {
				S: aws.String(storeID),
			},
		},
		KeyConditionExpression: aws.String("store_id = :store_id"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String(os.Getenv("STORE_INDEX_NAME")),
		ScanIndexForward:       &scanIndexForward,
	}

	if result, err := db.Query(queryInput); err == nil {
		dynamodbattribute.UnmarshalListOfMaps(result.Items, &sales)
	}

	return sales
}

func main() {
	lambda.Start(ListSales)
}
