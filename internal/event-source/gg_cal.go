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

func GetTodayLeavesEvent() ([]string, error) {
	// calendatID := os.Getenv("CALENDAR_ID")
	calendatID := "fdtiu7e9tp0i07753g787egrdo@group.calendar.google.com"
	ctx := context.Background()

	// Create client
	// Get the credentials JSON content from the environment variable
	encodedCredentials := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if encodedCredentials == "" {
		return nil, fmt.Errorf("environment variable GOOGLE_CREDENTIALS_JSON is not set")
	}

	// Decode the Base64-encoded credentials (if applicable)
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

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	t := startOfToday.Format(time.RFC3339)
	endOfToday := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	log.Printf("Get leave event of : %s", now.Format(time.DateOnly))
	todayLeavesEvent, err := srv.Events.List(calendatID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).TimeMax(endOfToday).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	var eventSummaries []string
	for _, item := range todayLeavesEvent.Items {
		eventSummaries = append(eventSummaries, item.Summary)
	}

	return eventSummaries, nil
}
