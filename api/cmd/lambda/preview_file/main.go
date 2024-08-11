package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/johnnynu/agreatchaos/api/internal/handlers"
)

func main() {
	log.Println("Lambda function starting")
	lambda.Start(handlers.PreviewFile)
}