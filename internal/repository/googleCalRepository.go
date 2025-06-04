package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarRepository struct {
	calendarID                     string
	base64GoogleCalendarCredential string
}

func NewGoogleCalendarRepository(calendarID string, base64GoogleCalendarCredential string) *GoogleCalendarRepository {
	return &GoogleCalendarRepository{calendarID: calendarID, base64GoogleCalendarCredential: base64GoogleCalendarCredential}
}

func (g *GoogleCalendarRepository) GetEvents(asOf time.Time) ([]string, error) {
	ctx := context.Background()

	// Decode the Base64-encoded credentials
	credentials, err := base64.StdEncoding.DecodeString(g.base64GoogleCalendarCredential)
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

	startOfToday := time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 0, 0, 0, 0, asOf.Location()).Format(time.RFC3339)
	endOfToday := time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), asOf.Location()).Format(time.RFC3339)
	log.Printf("Get event of : %s, from calendar : %s", asOf.Format(time.DateOnly), g.calendarID)

	todayLeavesEvent, err := srv.Events.List(g.calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(startOfToday).TimeMax(endOfToday).MaxResults(50).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	var eventSummaries []string
	for _, item := range todayLeavesEvent.Items {
		eventSummaries = append(eventSummaries, item.Summary)
	}

	return eventSummaries, nil
}
