// pkg/utils/response.go
package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

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
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       err.Error(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}