package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tsongpon/iris/internal/model"
	"github.com/tsongpon/iris/internal/service"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var _ service.EventRepository = (*GoogleCalendarEventRepository)(nil)

// GoogleCalendarEventRepository reads events from a Google Calendar using a
// service account.
type GoogleCalendarEventRepository struct {
	srv        *calendar.Service
	calendarID string
	loc        *time.Location
}

// NewGoogleCalendarEventRepository creates a repository from raw service
// account credentials JSON. It never touches disk, so credentials can come
// from an environment variable or secret manager. If credentialsJSON is empty,
// Application Default Credentials are used instead, which is the natural fit
// for containers on GCP (Workload Identity, Cloud Run, GKE, GCE): no file and
// no secret need to be provisioned.
func NewGoogleCalendarEventRepository(ctx context.Context, credentialsJSON []byte, calendarID string, loc *time.Location) (*GoogleCalendarEventRepository, error) {
	if calendarID == "" {
		return nil, fmt.Errorf("calendar id is required")
	}
	if loc == nil {
		loc = time.Local
	}
	opts := []option.ClientOption{
		option.WithScopes(calendar.CalendarReadonlyScope),
	}
	if len(credentialsJSON) > 0 {
		opts = append(opts, option.WithAuthCredentialsJSON(option.ServiceAccount, credentialsJSON))
	}
	srv, err := calendar.NewService(ctx, opts...)
	if err != nil {
		slog.ErrorContext(ctx, "create calendar service failed",
			"calendar", calendarID,
			"err", err,
		)
		return nil, fmt.Errorf("create calendar service: %w", err)
	}
	return &GoogleCalendarEventRepository{srv: srv, calendarID: calendarID, loc: loc}, nil
}

func (r *GoogleCalendarEventRepository) GetEvents(ctx context.Context, asOf time.Time) ([]model.Event, error) {
	asOfLoc := asOf.In(r.loc)
	start := time.Date(asOfLoc.Year(), asOfLoc.Month(), asOfLoc.Day(), 0, 0, 0, 0, r.loc)
	end := start.Add(24 * time.Hour)

	var events []model.Event
	pageToken := ""
	for {
		call := r.srv.Events.List(r.calendarID).
			TimeMin(start.Format(time.RFC3339)).
			TimeMax(end.Format(time.RFC3339)).
			SingleEvents(true).
			OrderBy("startTime").
			MaxResults(250)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		page, err := call.Do()
		if err != nil {
			var apiErr *googleapi.Error
			attrs := []any{
				"calendar", r.calendarID,
				"timeMin", start.Format(time.RFC3339),
				"timeMax", end.Format(time.RFC3339),
			}
			if errors.As(err, &apiErr) {
				attrs = append(attrs,
					"status", apiErr.Code,
					"apiMessage", apiErr.Message,
					"apiBody", apiErr.Body,
				)
			}
			slog.ErrorContext(ctx, "list calendar events failed", attrs...)
			return nil, fmt.Errorf("list events from calendar %q: %w", r.calendarID, err)
		}
		for _, item := range page.Items {
			if item.Status == "cancelled" {
				continue
			}
			events = append(events, model.Event{
				Name:        item.Summary,
				Description: item.Description,
				Time:        r.eventTime(item),
			})
		}
		if page.NextPageToken == "" {
			break
		}
		pageToken = page.NextPageToken
	}
	return events, nil
}

// eventTime maps the event start to a time.Time. All-day events only carry a
// date, so they map to midnight in the repository location; timed events are
// parsed from their RFC3339 dateTime.
func (r *GoogleCalendarEventRepository) eventTime(item *calendar.Event) time.Time {
	if item.Start != nil && item.Start.DateTime != "" {
		if t, err := time.Parse(time.RFC3339, item.Start.DateTime); err == nil {
			return t
		}
	}
	if item.Start != nil && item.Start.Date != "" {
		if t, err := time.ParseInLocation("2006-01-02", item.Start.Date, r.loc); err == nil {
			return t
		}
	}
	return time.Time{}
}
