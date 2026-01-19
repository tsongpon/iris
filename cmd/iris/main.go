package main

import (
	"context"
	"log"
	"os"
	"time"

	"gitbub.com/tsongpon/iris/internal/repository"
	"gitbub.com/tsongpon/iris/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func newEventNotifyServive() (service.EventNotifyService, error) {
	leaveCalendarID := os.Getenv("LEAVE_CALENDAR_ID")
	holidayCalendarID := os.Getenv("HOLIDAY_CALENDAR_ID")
	onCallCalendarID := os.Getenv("ON_CALL_CALENDAR_ID")
	googleCalendarCredential := os.Getenv("GOOGLE_CREDENTIALS_JSON")

	lineGroupID := os.Getenv("LINE_GROUP_ID")
	lineChannelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")

	leaveEventRepository := repository.NewGoogleCalendar(googleCalendarCredential, leaveCalendarID)
	holidayEventRepository := repository.NewGoogleCalendar(googleCalendarCredential, holidayCalendarID)
	onCallEventRepository := repository.NewGoogleCalendar(googleCalendarCredential, onCallCalendarID)
	notificationRepo := repository.NewLineNotificationRepository(lineGroupID, lineChannelSecret, lineChannelToken)
	eventNotify := service.NewEventNotifyService(leaveEventRepository, holidayEventRepository, onCallEventRepository, notificationRepo)

	return eventNotify, nil
}

// Handle call from AWS Lambda
func HandleRequest(ctx context.Context) error {
	log.Printf("Running Lambda hendler function")
	var err error
	service, err := newEventNotifyServive()
	if err != nil {
		log.Printf("Error creating event handler: %v", err)
		return err
	}

	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal("Error loading location ", err)
	}
	asOf := time.Now().In(bangkok)
	err = service.Notify(asOf)
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
		} else {
			log.Println("Loaded .env file")
		}
		service, err := newEventNotifyServive()
		if err != nil {
			log.Printf("Unable to create event handler: %v", err)
			log.Fatal("Error while create service ", err)
		}
		bangkok, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			log.Fatal("Error loading location ", err)
		}
		asOf := time.Now().In(bangkok)
		service.Notify(asOf)
	}
}
