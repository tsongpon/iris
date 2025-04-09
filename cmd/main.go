package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"gitlab.com/tsongpon/iris/internal/eventsource"
	"gitlab.com/tsongpon/iris/internal/handler"
	"gitlab.com/tsongpon/iris/internal/notichannel"
)

func newEventHandler() (*handler.EventHandler, error) {
	leaveEventSource := eventsource.NewGoogleCalendar(os.Getenv("LEAVE_CALENDAR_ID"), os.Getenv("GOOGLE_CREDENTIALS_JSON"))
	notiChannel := notichannel.NewLineNoti(os.Getenv("LINE_GROUP_ID"), os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))

	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return nil, err
	}

	nowInBangkok := time.Now().In(bangkok)

	eventHandler := handler.NewEventHandler(leaveEventSource, notiChannel, nowInBangkok)
	return eventHandler, nil
}

func HandleRequest(ctx context.Context) error {
	log.Printf("Running Lambda hendler function")
	var err error
	handler, err := newEventHandler()
	if err != nil {
		log.Printf("Error creating event handler: %v", err)
		return err
	}
	err = handler.HandleEvent()
	if err != nil {
		log.Printf("Error handling event: %v", err)
		return err
	}
	log.Printf("Lambda handler function finished")
	return nil
}

func main() {
	isLabbda := os.Getenv("IS_LAMBDA")
	if isLabbda == "true" {
		log.Printf("Running in AWS Lambda")
		lambda.Start(HandleRequest)
	} else {
		if err := godotenv.Load(); err != nil {
			log.Printf("Unable to load .env file")
		}
		handler, err := newEventHandler()
		if err != nil {
			log.Fatalf("Error creating event handler: %v", err)
		}
		if err := handler.HandleEvent(); err != nil {
			log.Fatalf("Error handling event: %v", err)
		}
	}
}
