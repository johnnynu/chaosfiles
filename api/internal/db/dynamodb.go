package db

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	dbClient = dynamodb.NewFromConfig(cfg)
}

type User struct {
	UID    string `dynamodbav:"uid"`
	Username  string `dynamodbav:"username"`
	Email     string `dynamodbav:"email"`
	CreatedAt string `dynamodbav:"created_at"`
}

type File struct {
	FileID    string `dynamodbav:"FileID"`
	UserID    string `dynamodbav:"UserID"`
	FileName  string `dynamodbav:"FileName"`
	FileSize  int64  `dynamodbav:"FileSize"`
	FileType  string `dynamodbav:"FileType"`
	CreatedAt string `dynamodbav:"CreatedAt"`
	UpdatedAt string `dynamodbav:"UpdatedAt"`
}

func CreateUser(ctx context.Context, user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item: item,
	})

	if err != nil {
		return fmt.Errorf("failed to create user in DynamoDB: %v", err)
	}

	return nil
}

func GetUser(ctx context.Context, uid string) (*User, error) {
	res, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(("users")),
		Key: map[string]types.AttributeValue {
			"uid": &types.AttributeValueMemberS{Value: uid},
		},
	})

	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, nil
	}

	var user User
	err = attributevalue.UnmarshalMap(res.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateFile(ctx context.Context, file File) error {
	item, err := attributevalue.MarshalMap(file)
	if err != nil {
		return fmt.Errorf("failed to marshal file: %v", err)
	}

	// Log the marshalled item
	log.Printf("Marshalled item: %+v", item)

	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("FileMetadata"), // Make sure this matches your actual table name
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to create file in DynamoDB: %v", err)
	}

	return nil
}

func GetFile(ctx context.Context, fileID string) (*File, error) {
	res, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("FileMetadata"),
		Key: map[string]types.AttributeValue{
			"FileID": &types.AttributeValueMemberS{Value: fileID},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	if res.Item == nil {
		return nil, nil
	}

	var file File
	err = attributevalue.UnmarshalMap(res.Item, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func ListUserFiles(ctx context.Context, userID string) ([]File, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String("FileMetadata"),
		IndexName: aws.String("UserID-index"),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	}

	res, err := dbClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query user files: %v", err)
	}

	var files []File
	err = attributevalue.UnmarshalListOfMaps(res.Items, &files)
	if err != nil {
		return nil, fmt.Errorf("failed the unmarshal files: %v", err)
	}

	return files, nil
}

func UpdateFile(ctx context.Context, file File) error {
	update := expression.Set(expression.Name("FileSize"), expression.Value(file.FileSize)).
		Set(expression.Name("FileType"), expression.Value(file.FileType)).
		Set(expression.Name("UpdatedAt"), expression.Value(file.UpdatedAt))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("couldnt build expression for update: %v\n", err)
		return err
	}

	_, err = dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("FileMetadata"),
		Key: map[string]types.AttributeValue{
			"FileID": &types.AttributeValueMemberS{Value: file.FileID},
		},
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression: expr.Update(),
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		log.Printf("couldnt update file %v: %v", file.FileID, err)
		return err
	}

	log.Printf("successfully updated file metadata for fileID: %s", file.FileID)
	return nil
}

func DeleteFile(ctx context.Context, fileID string, userID string) error {
	file, err := GetFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}
	if file == nil {
		return fmt.Errorf("file not found")
	}

	if file.UserID != userID {
		return fmt.Errorf("unauthorized: file does not belong to the user")
	}

	_, err = dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("FileMetadata"),
		Key: map[string]types.AttributeValue{
			"FileID": &types.AttributeValueMemberS{Value: fileID},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}