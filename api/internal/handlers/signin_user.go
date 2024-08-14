package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/johnnynu/agreatchaos/api/internal/db"
)

func SigninUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("TESTReceived request headers: %v+", request.Headers)
	log.Println()

	// Extract cognito token from request
	authHeader := request.Headers["authorization"]
	log.Printf("TEST AUTH HEADER: %s", authHeader)
	log.Println()
	if authHeader == "" {
		log.Println("No auth token provided")
		return events.APIGatewayProxyResponse{StatusCode: 401, Body: "No token provided"}, nil
	}

	log.Printf("Token received: %s...", authHeader[:10])

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
        log.Println("Empty token after removing Bearer prefix")
        return events.APIGatewayProxyResponse{StatusCode: 401, Body: "Invalid token format"}, nil
    }

	log.Printf("Token received (first 20 chars): %s...", token[:20])

	// verify cognito token
	claims, err := verifyToken(ctx, token)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 401}, err
	}

	// extract user info from claims
	uid, ok := claims["sub"].(string)
	if !ok {
		log.Println("Unable to extract user ID from token claims")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid token content"}, nil
	}
	email, _ := claims["email"].(string)

	log.Printf("Extracted user info - UID: %s, Email: %s", uid, email)

	// process any additional metadata from request body

	// Check if user already exists in db
	existingUser, err := db.GetUser(ctx, uid)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal server error"}, err
	}

	isNewUser := existingUser == nil

	if isNewUser {
		user := db.User{
			UID:    uid,
			Email:     email,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		// Call CreateUser function
		err = db.CreateUser(ctx, user)
		if err != nil {
			log.Printf("Error creating new user: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to create new user"}, err
		}
		log.Println("New user created successfully")
	} else {
		log.Println("Existing user signed in")
	}

	res := struct {
		Message   string `json:"message"`
		IsNewUser bool   `json:"isNewUser"`
	}{
		Message:   "Sign-in successful",
		IsNewUser: isNewUser,
	}

	resBody, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error creating response: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating response"}, nil
	}

	log.Println("Signin process completed successfully")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(resBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func verifyToken(ctx context.Context, token string) (map[string]interface{}, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	result, err := cognitoClient.GetUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error verifying token, %v", err)
	}

	claims := make(map[string]interface{})
	for _, attr := range result.UserAttributes {
		claims[*attr.Name] = *attr.Value
	}

	log.Println("Token verified successfully")
	return claims, nil
}
