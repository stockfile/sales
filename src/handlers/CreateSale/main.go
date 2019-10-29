package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stockfile/sales/src/db"
	"github.com/stockfile/sales/src/models"
	"github.com/stockfile/sales/src/renderer"

	uuid "github.com/satori/go.uuid"
)

type Response struct {
	Sale models.Sale `json:"sale"`
}

func CreateSale(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sale := buildSale(request)
	tableName := os.Getenv("TABLE_NAME")

	if err := db.PutItem(sale, &tableName); err != nil {
		return renderer.RenderServerError(err.Error())
	} else {
		body, _ := json.Marshal(&Response{
			Sale: *sale,
		})
		return renderer.RenderSuccess(body)
	}
}

func buildSale(request events.APIGatewayProxyRequest) *models.Sale {
	userID := request.RequestContext.Authorizer["SF-User-Id"].(string)
	storeID := request.RequestContext.Authorizer["SF-Store-Id"].(string)

	id := uuid.Must(uuid.NewV4(), nil).String()
	timestamp := time.Now().Format(time.RFC3339)

	sale := &models.Sale{
		ID:        id,
		UserId:    userID,
		StoreId:   storeID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	json.Unmarshal([]byte(request.Body), sale)

	return sale
}

func main() {
	lambda.Start(CreateSale)
}
