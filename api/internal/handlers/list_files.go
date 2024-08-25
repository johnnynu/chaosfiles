package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func ListFiles(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %+v: ", request)

    // Extract userID from JWT claims
    var userID string
	if jwt, ok := request.RequestContext.Authorizer["jwt"].(map[string]interface{}); ok {
		if claims, ok := jwt["claims"].(map[string]interface{}); ok {
			log.Printf("claims: %v", claims)
			if sub, ok := claims["sub"].(string); ok {
				log.Printf("sub: %v", sub)
				userID = sub
			}
		}
	}

	if userID == "" {
        log.Println("Unable to extract user ID from JWT claims")
        return utils.ResponseError(fmt.Errorf("unable to extract user ID from JWT claims"))
    }

	files, err := db.ListUserFiles(ctx, userID)
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return utils.ResponseError(err)
	}

	// Convert files to JSON
	resBody, err := json.Marshal(files)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return utils.ResponseError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(resBody),
	}, nil
}