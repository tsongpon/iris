package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"gitlab.com/tsongpon/iris/internal/handler"
)

func HandleRequest(ctx context.Context) error {
	log.Printf("Running Lambda hendler function")
	handler.LeaveEventHandler()
	return nil
}

func main() {
	isLabbda := os.Getenv("IS_LAMBDA")
	if isLabbda == "true" {
		log.Printf("Running in AWS Lambda")
		lambda.Start(HandleRequest)
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Unable to load .env file")
		}
		handler.LeaveEventHandler()
	}
}
