package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/johnnynu/agreatchaos/api/internal/db"
	"github.com/johnnynu/agreatchaos/api/pkg/utils"
)

func DeleteFile(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %+v", request)

	fileID, ok := request.PathParameters["fileId"]
	if !ok || fileID == "" {
		log.Println("FileID not found in path parameters")
		return utils.ResponseError(fmt.Errorf("fileID is required"))
	}
	log.Printf("FileID to delete: %s", fileID)

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

	if userID == "" {
        log.Println("Unable to extract user ID from JWT claims")
        return utils.ResponseError(fmt.Errorf("unable to extract user ID from JWT claims"))
    }

	err := db.DeleteFile(ctx, fileID, userID)
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return utils.ResponseError(err)
	}

	err = deleteFileFromS3(ctx, fileID)
	if err != nil {
		log.Printf("Error deleting file from S3: %v", err)
		return utils.ResponseError(err)
	}

	log.Printf("File %s deleted successfully", fileID)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: fmt.Sprintf("File %s deleted successfully", fileID),
	}, nil
}

func deleteFileFromS3(ctx context.Context, fileID string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String("chaosfiles-filestorage"),
		Key: aws.String(fileID),
	})

	if err != nil {
		return fmt.Errorf("unable to delete file from S3, %v", err)
	}

	return nil
}