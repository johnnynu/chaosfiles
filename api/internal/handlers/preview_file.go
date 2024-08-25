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

func PreviewFile(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Received request: %+v", request)
    log.Printf("Context: %+v", ctx)
    log.Printf("Path parameters: %+v", request.PathParameters)
    log.Printf("Query string parameters: %+v", request.QueryStringParameters)
    log.Printf("Headers: %+v", request.Headers)

    fileID, ok := request.PathParameters["fileId"]
    if !ok || fileID == "" {
        log.Printf("FileID not found in path parameters")
        return utils.ResponseError(fmt.Errorf("fileID is required"))
    }
    log.Printf("FileID: %s", fileID)

    // get file details
    file, err := db.GetFile(ctx, fileID)
    if err != nil {
        log.Printf("Error getting file: %v", err)
        return utils.ResponseError(err)
    }

    if file == nil {
        log.Printf("File not found for ID: %s", fileID)
        return utils.ResponseError(utils.ErrNotFound)
    }

    log.Printf("File found: %+v", file)

    res, err := json.Marshal(file)
    if err != nil {
        log.Printf("Error marshalling response: %v", err)
        return utils.ResponseError(err)
    }

    log.Printf("Sending response: %s", string(res))

    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(res),
    }, nil
}