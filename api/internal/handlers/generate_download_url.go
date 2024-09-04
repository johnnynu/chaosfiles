package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func GenerateDownloadURL(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileID := request.QueryStringParameters["fileID"]
	if fileID == "" {
		return utils.ResponseError(errors.New("fileID is required"))
	}

	file, err := db.GetFile(ctx, fileID)
	if err != nil {
		return utils.ResponseError(err)
	}
	if file == nil {
		return utils.ResponseError(errors.New("file not found"))
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return utils.ResponseError(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Generate pre signed url
	presignClient := s3.NewPresignClient(s3Client)
	presignedUrl, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("chaosfiles-filestorage"),
		Key: aws.String(fileID),
		ResponseContentType: aws.String(file.FileType),
	}, s3.WithPresignExpires(time.Minute * 15))

	if err != nil {
		return utils.ResponseError(err)
	}

	// Return the presigned url
	return utils.ResponseOK(map[string]string{
		"downloadUrl": presignedUrl.URL,
		"fileName": file.FileName,
		"contentType": file.FileType,
	})
}