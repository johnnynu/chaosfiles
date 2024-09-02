package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/johnnynu/agreatchaos/api/internal/handlers"
)

func main() {
	lambda.Start(handlers.ProcessUpload)
}