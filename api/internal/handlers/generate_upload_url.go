package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

const (
	maxFileSize = 1 * 1024 * 1024 * 1024 * 1024 // 1TB
	maxParts    = 10000
	multipartThreshold = 100 * 1024 * 1024 // 100MB
)

type UploadURLRequest struct {
    FileName  string `json:"fileName"`
    FileType  string `json:"fileType"`
    FileSize  int64  `json:"fileSize"`
    ChunkSize int64  `json:"chunkSize"`
}

func GenerateUploadURL(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("GenerateUploadURL function started")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return utils.ResponseError(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	var  req UploadURLRequest
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return utils.ResponseError(errors.New("invalid request body"))
	}

	if req.FileName == "" {
		log.Println("fileName is required")
		return utils.ResponseError(errors.New("fileName is required"))
	}

	if req.FileName == "" || req.FileSize == 0 {
		log.Println("fileName and fileSize are required")
		return utils.ResponseError(errors.New("fileName and fileSize are required"))
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
		FileSize: req.FileSize,
        CreatedAt: time.Now().Format(time.RFC3339),
        UpdatedAt: time.Now().Format(time.RFC3339),
    }

    log.Printf("Attempting to create file: %+v", file)

    err = db.CreateFile(ctx, file)
    if err != nil {
        log.Printf("Error creating file: %v", err)
        return utils.ResponseError(err)
    }

	if req.FileSize < multipartThreshold {
	// Handle single part upload
	// Generate pre signed url
	presignClient := s3.NewPresignClient(s3Client)
	presignedUrl, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("chaosfiles-filestorage"),
		Key: aws.String(fileID),
		ContentType: aws.String(req.FileType),
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
	} else {
		// handle multipart upload
		if req.ChunkSize == 0 {
			log.Println("chunkSize is required for multipart uploads")
			return utils.ResponseError(errors.New("chunkSize is required for multipart uploads"))
		}

		numParts := int(math.Ceil(float64(req.FileSize) / float64(req.ChunkSize)))
		if numParts > maxParts {
			log.Printf("Number of parts %d exceeds maximum allowed parts %d", numParts, maxParts)
			return utils.ResponseError(errors.New("file size results in too many parts"))
		}

		// initiate multipart upload
		createResp, err := s3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String("chaosfiles-filestorage"),
			Key: aws.String(fileID),
			ContentType: aws.String(req.FileType),
		})
		if err != nil {
			log.Printf("Error creating multipart upload: %v", err)
			return utils.ResponseError(err)
		}

		uploadID := *createResp.UploadId

		// Generate pre-signed URLs for each part
		partUrls := make([]string, numParts)
		for i := 0; i < numParts; i++ {
			partNumber := int32(i + 1)
			presignedReq, err := presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
				Bucket: aws.String("chaosfiles-filestorage"),
				Key: aws.String(fileID),
				UploadId: aws.String(uploadID),
				PartNumber: &partNumber,
			}, s3.WithPresignExpires(time.Hour*24))

			if err != nil {
				log.Printf("Error generating pre-signed URL for part %d: %v", partNumber, err)
				return utils.ResponseError(err)
			}

			partUrls[i] = presignedReq.URL
		}

		response := struct {
			UploadID string   `json:"uploadId"`
			FileID   string   `json:"fileID"`
			PartUrls []string `json:"partUrls"`
		}{
			UploadID: uploadID,
			FileID:   fileID,
			PartUrls: partUrls,
		}

		return utils.ResponseOK(response)
	}
}