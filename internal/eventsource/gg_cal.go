package eventsource

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

type GoogleCalendar struct {
	calendarID                     string
	base64GoogleCalendarCredential string
}

func NewGoogleCalendar(calendarID string, base64GoogleCalendarCredential string) *GoogleCalendar {
	return &GoogleCalendar{calendarID: calendarID, base64GoogleCalendarCredential: base64GoogleCalendarCredential}
}

func (g *GoogleCalendar) GetEvents(asOf time.Time) ([]string, error) {
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

	startOfToday := time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 0, 0, 0, 0, asOf.Location())
	endOfToday := startOfToday.Add(24*time.Hour - time.Second)
	log.Printf("Get event of : %s, from calendar : %s", asOf.Format(time.DateOnly), g.calendarID)

	todayLeavesEvent, err := srv.Events.List(g.calendarID).ShowDeleted(false).
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
