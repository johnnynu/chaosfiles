package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

type UploadURLRequest struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
}

func GenerateUploadURL(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("GenerateUploadURL function started")

	var  req UploadURLRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return utils.ResponseError(errors.New("invalid request body"))
	}

	if req.FileName == "" {
		log.Println("fileName is required")
		return utils.ResponseError(errors.New("fileName is required"))
	}

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

	fileID := uuid.New().String()

    file := db.File{
        FileID:    fileID,
        UserID:    userID,
        FileName:  req.FileName,
		FileType: req.FileType,
        CreatedAt: time.Now().Format(time.RFC3339),
        UpdatedAt: time.Now().Format(time.RFC3339),
    }

    log.Printf("Attempting to create file: %+v", file)

    err = db.CreateFile(ctx, file)
    if err != nil {
        log.Printf("Error creating file: %v", err)
        return utils.ResponseError(err)
    }

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return utils.ResponseError(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Generate pre signed url
	presignClient := s3.NewPresignClient(s3Client)
	presignedUrl, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("chaosfiles-filestorage"),
		Key: aws.String(fileID),
	}, s3.WithPresignExpires(time.Minute * 15))

	if err != nil {
		return utils.ResponseError(err)
	}

    response := struct {
        UploadURL string `json:"uploadUrl"`
        FileID    string `json:"fileID"`
    }{
        UploadURL: presignedUrl.URL,
        FileID:    fileID,
    }

	// Return the presigned url
	return utils.ResponseOK(response)
}