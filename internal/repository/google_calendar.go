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

type GoogleCalendar struct {
	base64GoogleCalendarCredential string
	calendarID                     string
}

func NewGoogleCalendar(base64GoogleCalendarCredential, calendarID string) GoogleCalendar {
	return GoogleCalendar{
		base64GoogleCalendarCredential: base64GoogleCalendarCredential,
		calendarID:                     calendarID,
	}
}

func (g GoogleCalendar) GetEvents(asOf time.Time) ([]string, error) {
	ctx := context.Background()
	credential, err := base64.StdEncoding.DecodeString(g.base64GoogleCalendarCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 credential: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credential, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	beginningOfDay := time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 9, 0, 0, 0, asOf.Location()).Format(time.RFC3339)
	endOfDay := time.Date(asOf.Year(), asOf.Month(), asOf.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), asOf.Location()).Format(time.RFC3339)
	log.Printf("Get event of : %s, from calendar : %s", asOf.Format(time.DateOnly), g.calendarID)

	todayLeavesEvent, err := srv.Events.List(g.calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(beginningOfDay).TimeMax(endOfDay).MaxResults(50).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	var eventSummaries []string
	for _, item := range todayLeavesEvent.Items {
		eventSummaries = append(eventSummaries, item.Summary)
	}

	return eventSummaries, nil
}

func (g GoogleCalendar) GetEventsBetween(start, end time.Time) ([]string, error) {
	ctx := context.Background()
	credential, err := base64.StdEncoding.DecodeString(g.base64GoogleCalendarCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 credential: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credential, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	beginningOfPeriod := start.Format(time.RFC3339)
	endOfPeriod := end.Format(time.RFC3339)
	log.Printf("Get event between : %s and %s, from calendar : %s", beginningOfPeriod, endOfPeriod, g.calendarID)

	events, err := srv.Events.List(g.calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(beginningOfPeriod).TimeMax(endOfPeriod).MaxResults(50).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	var eventSummaries []string
	for _, item := range events.Items {
		date := item.Start.Date
		eventSummaries = append(eventSummaries, date+": "+item.Summary)
	}

	return eventSummaries, nil
}
