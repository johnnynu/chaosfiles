package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func CreateFile(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Received request: %+v", request)

    var input struct {
        FileName string `json:"file_name"`
        FileSize int64  `json:"file_size"`
        FileType string `json:"file_type"`
    }

    err := json.Unmarshal([]byte(request.Body), &input)
    if err != nil {
        log.Printf("Error unmarshalling request body: %v", err)
        return utils.ResponseError(err)
    }

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

    file := db.File{
        FileID:    uuid.New().String(),
        UserID:    userID,
        FileName:  input.FileName,
        FileSize:  input.FileSize,
        FileType:  input.FileType,
        CreatedAt: time.Now().Format(time.RFC3339),
        UpdatedAt: time.Now().Format(time.RFC3339),
    }

    log.Printf("Attempting to create file: %+v", file)

    err = db.CreateFile(ctx, file)
    if err != nil {
        log.Printf("Error creating file: %v", err)
        return utils.ResponseError(err)
    }

    return utils.ResponseOK(file)
}