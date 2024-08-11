package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func ListFiles(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return utils.ResponseOK("Listed files successfully")
}