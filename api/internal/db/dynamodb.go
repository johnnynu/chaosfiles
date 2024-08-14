package db

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	FileID    string `dynamodbav:"file_id"`
	UserID    string `dynamodbav:"user_id"`
	FileName  string `dynamodbav:"file_name"`
	FileSize  int64  `dynamodbav:"file_size"`
	FileType  string `dynamodbav:"file_type"`
	CreatedAt string `dynamodbav:"created_at"`
	UpdatedAt string `dynamodbav:"updated_at"`
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