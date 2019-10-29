package renderer

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseError struct {
	Error string `json:"error"`
}

func RenderServerError(errorMessage string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(&ResponseError{
		Error: errorMessage,
	})

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func RenderClientError(errorMessage string, statusCode int) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(&ResponseError{
		Error: errorMessage,
	})

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

func RenderSuccess(bodyMessage []byte) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       string(bodyMessage),
		StatusCode: http.StatusOK,
	}, nil
}
