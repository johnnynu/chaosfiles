package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func PreviewFile(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return utils.ResponseOK("Get file successfully")
}