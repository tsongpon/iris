package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/tsongpon/iris/internal/handler"
)

func HandleRequest(ctx context.Context) error {
	handler.LeaveEventHandler()
	return nil
}

func main() {
	isLabbda := os.Getenv("IS_LAMBDA")
	if isLabbda == "true" {
		log.Printf("Running in AWS Lambda")
		lambda.Start(HandleRequest)
	} else {
		handler.LeaveEventHandler()
	}
}
