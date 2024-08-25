package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/johnnynu/agreatchaos/api/internal/db"
)

func HandleStream(ctx context.Context, e events.DynamoDBEvent) error {
	for _, record := range e.Records {
		if record.EventName == "INSERT" || record.EventName == "MODIFY" {
			convertedImage, err := convertDDBStreamImage(record.Change.NewImage)
			if err != nil {
				log.Printf("Error converting DynamoDB stream image: %v", err)
				continue
			}

			var file db.File
			err = attributevalue.UnmarshalMap(convertedImage, &file)
			if err != nil {
				log.Printf("Error unmarshalling DynamoDB event: %v", err)
				continue
			}

			fmt.Printf("File updated: %s\n", file.FileName)
			// Here you can implement real-time notifications or trigger other processes
		} else if record.EventName == "REMOVE" {
			convertedImage, err := convertDDBStreamImage(record.Change.OldImage)
			if err != nil {
				log.Printf("Error converting DynamoDB stream image: %v", err)
				continue
			}

			var file db.File
			err = attributevalue.UnmarshalMap(convertedImage, &file)
			if err != nil {
				log.Printf("Error unmarshalling DynamoDB event: %v", err)
				continue
			}

			fmt.Printf("File deleted: %s\n", file.FileName)
			// Here you can implement cleanup processes or notifications
		}
	}

	return nil
}

func convertDDBStreamImage(image map[string]events.DynamoDBAttributeValue) (map[string]types.AttributeValue, error) {
	converted := make(map[string]types.AttributeValue)

	for k, v := range image {
		switch v.DataType() {
		case events.DataTypeString:
			converted[k] = &types.AttributeValueMemberS{Value: v.String()}
		case events.DataTypeNumber:
			converted[k] = &types.AttributeValueMemberN{Value: v.Number()}
		case events.DataTypeBinary:
			converted[k] = &types.AttributeValueMemberB{Value: v.Binary()}
		case events.DataTypeBoolean:
			converted[k] = &types.AttributeValueMemberBOOL{Value: v.Boolean()}
		case events.DataTypeNull:
			converted[k] = &types.AttributeValueMemberNULL{Value: v.IsNull()}
		case events.DataTypeList:
			listValues, err := convertDDBStreamList(v.List())
			if err != nil {
				return nil, err
			}
			converted[k] = &types.AttributeValueMemberL{Value: listValues}
		case events.DataTypeMap:
			mapValues, err := convertDDBStreamImage(v.Map())
			if err != nil {
				return nil, err
			}
			converted[k] = &types.AttributeValueMemberM{Value: mapValues}
		case events.DataTypeStringSet:
			converted[k] = &types.AttributeValueMemberSS{Value: v.StringSet()}
		case events.DataTypeNumberSet:
			converted[k] = &types.AttributeValueMemberNS{Value: v.NumberSet()}
		case events.DataTypeBinarySet:
			converted[k] = &types.AttributeValueMemberBS{Value: v.BinarySet()}
		default:
			return nil, fmt.Errorf("unsupported data type: %v", v.DataType())
		}
	}

	return converted, nil
}

func convertDDBStreamList(list []events.DynamoDBAttributeValue) ([]types.AttributeValue, error) {
	converted := make([]types.AttributeValue, len(list))

	for i, v := range list {
		switch v.DataType() {
		case events.DataTypeString:
			converted[i] = &types.AttributeValueMemberS{Value: v.String()}
		case events.DataTypeNumber:
			converted[i] = &types.AttributeValueMemberN{Value: v.Number()}
		case events.DataTypeBinary:
			converted[i] = &types.AttributeValueMemberB{Value: v.Binary()}
		case events.DataTypeBoolean:
			converted[i] = &types.AttributeValueMemberBOOL{Value: v.Boolean()}
		case events.DataTypeNull:
			converted[i] = &types.AttributeValueMemberNULL{Value: v.IsNull()}
		case events.DataTypeList:
			listValues, err := convertDDBStreamList(v.List())
			if err != nil {
				return nil, err
			}
			converted[i] = &types.AttributeValueMemberL{Value: listValues}
		case events.DataTypeMap:
			mapValues, err := convertDDBStreamImage(v.Map())
			if err != nil {
				return nil, err
			}
			converted[i] = &types.AttributeValueMemberM{Value: mapValues}
		case events.DataTypeStringSet:
			converted[i] = &types.AttributeValueMemberSS{Value: v.StringSet()}
		case events.DataTypeNumberSet:
			converted[i] = &types.AttributeValueMemberNS{Value: v.NumberSet()}
		case events.DataTypeBinarySet:
			converted[i] = &types.AttributeValueMemberBS{Value: v.BinarySet()}
		default:
			return nil, fmt.Errorf("unsupported data type: %v", v.DataType())
		}
	}

	return converted, nil
}