// pkg/utils/response.go
package utils

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

var ErrNotFound = errors.New("resource not found")

func ResponseOK(body interface{}) (events.APIGatewayProxyResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return ResponseError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func ResponseError(err error) (events.APIGatewayProxyResponse, error) {
	statusCode := 500
	if errors.Is(err, ErrNotFound) {
		statusCode = 404
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       err.Error(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}