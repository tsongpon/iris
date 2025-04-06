package eventsource

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetTodayEventFrom(calendarID string) ([]string, error) {
	ctx := context.Background()

	// Create client
	// Get the credentials JSON content from the environment variable
	encodedCredentials := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if encodedCredentials == "" {
		return nil, fmt.Errorf("environment variable GOOGLE_CREDENTIALS_JSON is not set")
	}

	// Decode the Base64-encoded credentials
	credentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credentials, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, fmt.Errorf("failed to load Bangkok timezone: %v", err)
	}
	now := time.Now().In(location)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.Add(24*time.Hour - time.Second)
	log.Printf("Get event of : %s, from calendar : %s", now.Format(time.DateOnly), calendarID)
	todayLeavesEvent, err := srv.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(startOfToday.Format(time.RFC3339)).TimeMax(endOfToday.Format(time.RFC3339)).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	var eventSummaries []string
	for _, item := range todayLeavesEvent.Items {
		eventSummaries = append(eventSummaries, item.Summary)
	}

	return eventSummaries, nil
}
