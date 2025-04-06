package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	eventsource "gitlab.com/tsongpon/iris/internal/event-source"
	"gitlab.com/tsongpon/iris/internal/handler"
	notichannel "gitlab.com/tsongpon/iris/internal/noti-channel"
)

func newEventHandler() *handler.EventHandler {
	holidayEventSource := eventsource.NewGoogleCalendar(os.Getenv("HOLIDAY_CALENDAR_ID"), os.Getenv("GOOGLE_CREDENTIALS_JSON"))
	leaveEventSource := eventsource.NewGoogleCalendar(os.Getenv("LEAVE_CALENDAR_ID"), os.Getenv("GOOGLE_CREDENTIALS_JSON"))
	notiChannel := notichannel.NewLineNoti(os.Getenv("LINE_GROUP_ID"), os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))

	eventHandler := handler.NewEventHandler(leaveEventSource, holidayEventSource, notiChannel)
	return eventHandler
}

func HandleRequest(ctx context.Context) error {
	log.Printf("Running Lambda hendler function")
	handler := newEventHandler()
	err := handler.HandleEvent()
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
		handler := newEventHandler()
		if err := handler.HandleEvent(); err != nil {
			log.Fatalf("Error handling event: %v", err)
		}
	}
}
