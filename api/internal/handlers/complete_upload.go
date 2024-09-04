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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

type CompleteUploadRequest struct {
	FileID   string `json:"fileID"`
	UploadID string `json:"uploadId"`
	Parts    []Part `json:"parts"`
}

type Part struct {
	ETag       string `json:"ETag"`
	PartNumber int32  `json:"PartNumber"`
}

func CompleteUpload(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("CompleteUpload function started")

	var req CompleteUploadRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("error unmarshalling request body: %v", err)
		return utils.ResponseError(errors.New("invalid request body"))
	}

	if req.FileID == "" || req.UploadID == "" || len(req.Parts) == 0 {
		log.Println("fileID, uploadID, and/or parts are required")
		return utils.ResponseError(errors.New("fileID, uploadID, and/or parts are required"))
	}

	// Extract user ID from JWT claims
	var userID string
	if jwt, ok := request.RequestContext.Authorizer["jwt"].(map[string]interface{}); ok {
		if claims, ok := jwt["claims"].(map[string]interface{}); ok {
			if sub, ok := claims["sub"].(string); ok {
				userID = sub
			}
		}
	}

	if userID == "" {
		log.Println("Unable to extract user ID from JWT claims")
		return utils.ResponseError(fmt.Errorf("unable to extract user ID from JWT claims"))
	}

	// verify file ownership and get file metadata
	file, err := db.GetFile(ctx, req.FileID)
	if err != nil {
		log.Printf("error fetching file metadata: %v", err)
		return utils.ResponseError(err)
	}

	if file == nil {
		log.Printf("file not found: %s", req.FileID)
		return utils.ResponseError(errors.New("file not found"))
	}

	if file.UserID != userID {
		log.Printf("user %s does not own file %s", userID, req.FileID)
	}

	// set up s3 client
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return utils.ResponseError(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Prepare completed parts for s3 api
	completedParts := make([]types.CompletedPart, len(req.Parts))
	for i, part := range req.Parts {
		completedParts[i] = types.CompletedPart{
			ETag: aws.String(part.ETag),
			PartNumber: aws.Int32(part.PartNumber),
		}
	}

	// call CompleteMultipartUpload
	_, err = s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket: aws.String("chaosfiles-filestorage"),
		Key: aws.String(req.FileID),
		UploadId: aws.String(req.UploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	if err != nil {
		log.Printf("error completing multipart upload: %v", err)
		return utils.ResponseError(err)
	}

	// update file status in the database
	file.UpdatedAt = time.Now().Format(time.RFC3339)
	err = db.UpdateFile(ctx, *file)
	if err != nil {
		log.Printf("error updating file metadata: %v", err)
		return utils.ResponseError(err)
	}

	log.Printf("multipart upload completed successfully for file: %s", req.FileID)

	return utils.ResponseOK(map[string]string{
		"message": "Upload completed successfully",
		"fileID":  req.FileID,
	})
}